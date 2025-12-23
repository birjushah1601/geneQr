package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BulkImportRequest represents the CSV import request
type BulkImportRequest struct {
	File         io.Reader
	CreatedBy    string
	DryRun       bool   // If true, validate only, don't insert
	UpdateMode   bool   // If true, update existing records
}

// BulkImportResponse represents the import result
type BulkImportResponse struct {
	TotalRows     int                  `json:"total_rows"`
	SuccessCount  int                  `json:"success_count"`
	FailureCount  int                  `json:"failure_count"`
	Errors        []ImportError        `json:"errors,omitempty"`
	ImportedIDs   []string             `json:"imported_ids,omitempty"`
	DryRun        bool                 `json:"dry_run"`
}

// ImportError represents an error for a specific row
type ImportError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// OrganizationCSVRow represents a single CSV row
type OrganizationCSVRow struct {
	Name         string
	OrgType      string // manufacturer|supplier|distributor|dealer|hospital|service_provider|other
	Status       string // active|inactive|suspended
	ExternalRef  string
	GSTIN        string
	PAN          string
	Website      string
	Email        string
	Phone        string
	Address      string
	City         string
	State        string
	Country      string
	Pincode      string
}

// BulkImportHandler handles CSV upload and bulk import
type BulkImportHandler struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewBulkImportHandler creates a new bulk import handler
func NewBulkImportHandler(db *pgxpool.Pool, logger *slog.Logger) *BulkImportHandler {
	return &BulkImportHandler{
		db:     db,
		logger: logger,
	}
}

// HandleBulkImport processes the CSV upload
func (h *BulkImportHandler) HandleBulkImport(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the CSV file
	file, _, err := r.FormFile("csv_file")
	if err != nil {
		http.Error(w, "CSV file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get optional parameters
	createdBy := r.FormValue("created_by")
	if createdBy == "" {
		createdBy = "bulk-import"
	}

	dryRun := r.FormValue("dry_run") == "true"
	updateMode := r.FormValue("update_mode") == "true"

	// Process the import
	request := BulkImportRequest{
		File:       file,
		CreatedBy:  createdBy,
		DryRun:     dryRun,
		UpdateMode: updateMode,
	}

	result, err := h.ProcessImport(r.Context(), request)
	if err != nil {
		h.logger.Error("Bulk import failed", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Import failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	if result.FailureCount > 0 {
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

// ProcessImport processes the CSV import
func (h *BulkImportHandler) ProcessImport(ctx context.Context, req BulkImportRequest) (*BulkImportResponse, error) {
	response := &BulkImportResponse{
		Errors:      []ImportError{},
		ImportedIDs: []string{},
		DryRun:      req.DryRun,
	}

	// Parse CSV
	reader := csv.NewReader(req.File)
	reader.TrimLeadingSpace = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Detect column mapping
	columnMap := h.detectColumnMapping(headers)
	if columnMap["name"] == -1 || columnMap["org_type"] == -1 {
		return nil, fmt.Errorf("required columns missing: name and org_type are required")
	}

	// Start transaction if not dry run
	var tx pgx.Tx
	if !req.DryRun {
		tx, err = h.db.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback(ctx)
	}

	// Process rows
	rowNum := 1 // Start at 1 (header is row 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			response.Errors = append(response.Errors, ImportError{
				Row:     rowNum,
				Message: fmt.Sprintf("CSV parse error: %v", err),
			})
			response.FailureCount++
			rowNum++
			continue
		}

		rowNum++
		response.TotalRows++

		// Parse row
		orgRow, err := h.parseRow(record, columnMap)
		if err != nil {
			response.Errors = append(response.Errors, ImportError{
				Row:     rowNum,
				Message: err.Error(),
				Data:    strings.Join(record, ","),
			})
			response.FailureCount++
			continue
		}

		// Validate row
		if err := h.validateRow(orgRow); err != nil {
			response.Errors = append(response.Errors, ImportError{
				Row:     rowNum,
				Message: err.Error(),
				Data:    fmt.Sprintf("Name: %s, Type: %s", orgRow.Name, orgRow.OrgType),
			})
			response.FailureCount++
			continue
		}

		// Insert or update
		if req.DryRun {
			// Dry run - just count success
			response.SuccessCount++
		} else {
			orgID, err := h.insertOrUpdateOrganization(ctx, tx, orgRow, req.UpdateMode, req.CreatedBy)
			if err != nil {
				response.Errors = append(response.Errors, ImportError{
					Row:     rowNum,
					Message: fmt.Sprintf("Insert failed: %v", err),
					Data:    orgRow.Name,
				})
				response.FailureCount++
				continue
			}

			response.SuccessCount++
			response.ImportedIDs = append(response.ImportedIDs, orgID)
		}
	}

	// Commit transaction
	if !req.DryRun && tx != nil {
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	h.logger.Info("Bulk import completed",
		slog.Int("total", response.TotalRows),
		slog.Int("success", response.SuccessCount),
		slog.Int("failures", response.FailureCount),
		slog.Bool("dry_run", req.DryRun),
	)

	return response, nil
}

// detectColumnMapping maps CSV columns to struct fields
func (h *BulkImportHandler) detectColumnMapping(headers []string) map[string]int {
	mapping := map[string]int{
		"name":         -1,
		"org_type":     -1,
		"status":       -1,
		"external_ref": -1,
		"gstin":        -1,
		"pan":          -1,
		"website":      -1,
		"email":        -1,
		"phone":        -1,
		"address":      -1,
		"city":         -1,
		"state":        -1,
		"country":      -1,
		"pincode":      -1,
	}

	for i, header := range headers {
		normalized := strings.ToLower(strings.TrimSpace(header))
		normalized = strings.ReplaceAll(normalized, " ", "_")

		// Try exact match
		if _, ok := mapping[normalized]; ok {
			mapping[normalized] = i
			continue
		}

		// Try fuzzy matches
		switch {
		case strings.Contains(normalized, "name") && !strings.Contains(normalized, "company"):
			mapping["name"] = i
		case strings.Contains(normalized, "type") || strings.Contains(normalized, "org_type"):
			mapping["org_type"] = i
		case strings.Contains(normalized, "status"):
			mapping["status"] = i
		case strings.Contains(normalized, "gstin") || strings.Contains(normalized, "gst"):
			mapping["gstin"] = i
		case strings.Contains(normalized, "pan"):
			mapping["pan"] = i
		case strings.Contains(normalized, "website") || strings.Contains(normalized, "url"):
			mapping["website"] = i
		case strings.Contains(normalized, "email"):
			mapping["email"] = i
		case strings.Contains(normalized, "phone") || strings.Contains(normalized, "mobile"):
			mapping["phone"] = i
		case strings.Contains(normalized, "address") && !strings.Contains(normalized, "email"):
			mapping["address"] = i
		case strings.Contains(normalized, "city"):
			mapping["city"] = i
		case strings.Contains(normalized, "state"):
			mapping["state"] = i
		case strings.Contains(normalized, "country"):
			mapping["country"] = i
		case strings.Contains(normalized, "pin") || strings.Contains(normalized, "zip") || strings.Contains(normalized, "postal"):
			mapping["pincode"] = i
		}
	}

	return mapping
}

// parseRow parses a CSV row into OrganizationCSVRow
func (h *BulkImportHandler) parseRow(record []string, mapping map[string]int) (*OrganizationCSVRow, error) {
	getField := func(key string) string {
		idx := mapping[key]
		if idx == -1 || idx >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[idx])
	}

	row := &OrganizationCSVRow{
		Name:         getField("name"),
		OrgType:      getField("org_type"),
		Status:       getField("status"),
		ExternalRef:  getField("external_ref"),
		GSTIN:        getField("gstin"),
		PAN:          getField("pan"),
		Website:      getField("website"),
		Email:        getField("email"),
		Phone:        getField("phone"),
		Address:      getField("address"),
		City:         getField("city"),
		State:        getField("state"),
		Country:      getField("country"),
		Pincode:      getField("pincode"),
	}

	// Set defaults
	if row.Status == "" {
		row.Status = "active"
	}
	if row.Country == "" {
		row.Country = "India"
	}

	return row, nil
}

// validateRow validates an organization row
func (h *BulkImportHandler) validateRow(row *OrganizationCSVRow) error {
	// Required fields
	if row.Name == "" {
		return fmt.Errorf("name is required")
	}
	if row.OrgType == "" {
		return fmt.Errorf("org_type is required")
	}

	// Validate org_type
	validTypes := map[string]bool{
		"manufacturer":     true,
		"supplier":         true,
		"distributor":      true,
		"dealer":           true,
		"hospital":         true,
		"clinic":           true,
		"service_provider": true,
		"other":            true,
	}
	if !validTypes[row.OrgType] {
		return fmt.Errorf("invalid org_type: %s (must be one of: manufacturer, supplier, distributor, dealer, hospital, clinic, service_provider, other)", row.OrgType)
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":    true,
		"inactive":  true,
		"suspended": true,
	}
	if row.Status != "" && !validStatuses[row.Status] {
		return fmt.Errorf("invalid status: %s (must be: active, inactive, or suspended)", row.Status)
	}

	return nil
}

// insertOrUpdateOrganization inserts or updates an organization
func (h *BulkImportHandler) insertOrUpdateOrganization(ctx context.Context, tx pgx.Tx, row *OrganizationCSVRow, updateMode bool, createdBy string) (string, error) {
	// Build metadata JSON
	metadata := map[string]interface{}{}
	if row.GSTIN != "" {
		metadata["gstin"] = row.GSTIN
	}
	if row.PAN != "" {
		metadata["pan"] = row.PAN
	}
	if row.Website != "" {
		metadata["website"] = row.Website
	}
	if row.Email != "" {
		metadata["email"] = row.Email
	}
	if row.Phone != "" {
		metadata["phone"] = row.Phone
	}
	if row.Address != "" || row.City != "" || row.State != "" {
		address := map[string]string{
			"address": row.Address,
			"city":    row.City,
			"state":   row.State,
			"country": row.Country,
			"pincode": row.Pincode,
		}
		metadata["address"] = address
	}

	metadataJSON, _ := json.Marshal(metadata)

	// Check if organization exists (by name or external_ref)
	var existingID string
	checkQuery := `SELECT id FROM organizations WHERE name = $1 OR (external_ref IS NOT NULL AND external_ref = $2) LIMIT 1`
	err := tx.QueryRow(ctx, checkQuery, row.Name, row.ExternalRef).Scan(&existingID)

	if err == nil && existingID != "" {
		// Organization exists
		if !updateMode {
			return "", fmt.Errorf("organization already exists: %s (use update_mode=true to update)", row.Name)
		}

		// Update existing
		updateQuery := `
			UPDATE organizations 
			SET org_type = $1, status = $2, external_ref = $3, metadata = $4, updated_at = NOW()
			WHERE id = $5
			RETURNING id
		`
		err = tx.QueryRow(ctx, updateQuery, row.OrgType, row.Status, row.ExternalRef, metadataJSON, existingID).Scan(&existingID)
		if err != nil {
			return "", fmt.Errorf("failed to update organization: %w", err)
		}

		return existingID, nil
	}

	// Insert new organization
	orgID := uuid.New().String()
	insertQuery := `
		INSERT INTO organizations (id, name, org_type, status, external_ref, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	err = tx.QueryRow(ctx, insertQuery, orgID, row.Name, row.OrgType, row.Status, row.ExternalRef, metadataJSON).Scan(&orgID)
	if err != nil {
		return "", fmt.Errorf("failed to insert organization: %w", err)
	}

	return orgID, nil
}
