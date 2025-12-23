package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EquipmentCatalogCSVRow represents a single CSV row for equipment catalog
type EquipmentCatalogCSVRow struct {
	ProductCode                  string
	ProductName                  string
	ManufacturerName             string
	ModelNumber                  string
	Category                     string
	Subcategory                  string
	Description                  string
	BasePrice                    string
	Currency                     string
	WeightKg                     string
	RecommendedServiceIntervalDays string
	EstimatedLifespanYears       string
	MaintenanceComplexity        string
}

// CatalogBulkImportHandler handles equipment catalog CSV uploads
type CatalogBulkImportHandler struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewCatalogBulkImportHandler creates a new catalog bulk import handler
func NewCatalogBulkImportHandler(db *pgxpool.Pool, logger *slog.Logger) *CatalogBulkImportHandler {
	return &CatalogBulkImportHandler{
		db:     db,
		logger: logger,
	}
}

// BulkImportResponse represents the import result
type BulkImportResponse struct {
	TotalRows    int           `json:"total_rows"`
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
	Errors       []ImportError `json:"errors,omitempty"`
	ImportedIDs  []string      `json:"imported_ids,omitempty"`
	DryRun       bool          `json:"dry_run"`
}

// ImportError represents an error for a specific row
type ImportError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// HandleCatalogBulkImport processes the equipment catalog CSV upload
func (h *CatalogBulkImportHandler) HandleCatalogBulkImport(w http.ResponseWriter, r *http.Request) {
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
	result, err := h.ProcessImport(r.Context(), file, createdBy, dryRun, updateMode)
	if err != nil {
		h.logger.Error("Catalog bulk import failed", slog.String("error", err.Error()))
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
func (h *CatalogBulkImportHandler) ProcessImport(ctx context.Context, file io.Reader, createdBy string, dryRun, updateMode bool) (*BulkImportResponse, error) {
	response := &BulkImportResponse{
		Errors:      []ImportError{},
		ImportedIDs: []string{},
		DryRun:      dryRun,
	}

	// Parse CSV
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Detect column mapping
	columnMap := h.detectColumnMapping(headers)
	if columnMap["product_code"] == -1 || columnMap["product_name"] == -1 {
		return nil, fmt.Errorf("required columns missing: product_code and product_name are required")
	}

	// Start transaction if not dry run
	var tx pgx.Tx
	if !dryRun {
		tx, err = h.db.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback(ctx)
	}

	// Process rows
	rowNum := 1
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
		catalogRow, err := h.parseRow(record, columnMap)
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
		if err := h.validateRow(catalogRow); err != nil {
			response.Errors = append(response.Errors, ImportError{
				Row:     rowNum,
				Message: err.Error(),
				Data:    fmt.Sprintf("Product: %s", catalogRow.ProductName),
			})
			response.FailureCount++
			continue
		}

		// Insert or update
		if dryRun {
			response.SuccessCount++
		} else {
			catalogID, err := h.insertOrUpdateCatalog(ctx, tx, catalogRow, updateMode, createdBy)
			if err != nil {
				response.Errors = append(response.Errors, ImportError{
					Row:     rowNum,
					Message: fmt.Sprintf("Insert failed: %v", err),
					Data:    catalogRow.ProductName,
				})
				response.FailureCount++
				continue
			}

			response.SuccessCount++
			response.ImportedIDs = append(response.ImportedIDs, catalogID)
		}
	}

	// Commit transaction
	if !dryRun && tx != nil {
		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	h.logger.Info("Equipment catalog bulk import completed",
		slog.Int("total", response.TotalRows),
		slog.Int("success", response.SuccessCount),
		slog.Int("failures", response.FailureCount),
		slog.Bool("dry_run", dryRun),
	)

	return response, nil
}

// detectColumnMapping maps CSV columns to struct fields
func (h *CatalogBulkImportHandler) detectColumnMapping(headers []string) map[string]int {
	mapping := map[string]int{
		"product_code":                     -1,
		"product_name":                     -1,
		"manufacturer_name":                -1,
		"model_number":                     -1,
		"category":                         -1,
		"subcategory":                      -1,
		"description":                      -1,
		"base_price":                       -1,
		"currency":                         -1,
		"weight_kg":                        -1,
		"recommended_service_interval_days": -1,
		"estimated_lifespan_years":         -1,
		"maintenance_complexity":           -1,
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
		case strings.Contains(normalized, "product_code") || strings.Contains(normalized, "sku"):
			mapping["product_code"] = i
		case strings.Contains(normalized, "product_name") || strings.Contains(normalized, "name"):
			mapping["product_name"] = i
		case strings.Contains(normalized, "manufacturer"):
			mapping["manufacturer_name"] = i
		case strings.Contains(normalized, "model"):
			mapping["model_number"] = i
		case strings.Contains(normalized, "category"):
			mapping["category"] = i
		case strings.Contains(normalized, "subcategory"):
			mapping["subcategory"] = i
		case strings.Contains(normalized, "description"):
			mapping["description"] = i
		case strings.Contains(normalized, "price"):
			mapping["base_price"] = i
		case strings.Contains(normalized, "currency"):
			mapping["currency"] = i
		case strings.Contains(normalized, "weight"):
			mapping["weight_kg"] = i
		case strings.Contains(normalized, "service_interval"):
			mapping["recommended_service_interval_days"] = i
		case strings.Contains(normalized, "lifespan"):
			mapping["estimated_lifespan_years"] = i
		case strings.Contains(normalized, "complexity"):
			mapping["maintenance_complexity"] = i
		}
	}

	return mapping
}

// parseRow parses a CSV row
func (h *CatalogBulkImportHandler) parseRow(record []string, mapping map[string]int) (*EquipmentCatalogCSVRow, error) {
	getField := func(key string) string {
		idx := mapping[key]
		if idx == -1 || idx >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[idx])
	}

	row := &EquipmentCatalogCSVRow{
		ProductCode:                  getField("product_code"),
		ProductName:                  getField("product_name"),
		ManufacturerName:             getField("manufacturer_name"),
		ModelNumber:                  getField("model_number"),
		Category:                     getField("category"),
		Subcategory:                  getField("subcategory"),
		Description:                  getField("description"),
		BasePrice:                    getField("base_price"),
		Currency:                     getField("currency"),
		WeightKg:                     getField("weight_kg"),
		RecommendedServiceIntervalDays: getField("recommended_service_interval_days"),
		EstimatedLifespanYears:       getField("estimated_lifespan_years"),
		MaintenanceComplexity:        getField("maintenance_complexity"),
	}

	// Set defaults
	if row.Currency == "" {
		row.Currency = "USD"
	}

	return row, nil
}

// validateRow validates an equipment catalog row
func (h *CatalogBulkImportHandler) validateRow(row *EquipmentCatalogCSVRow) error {
	// Required fields
	if row.ProductCode == "" {
		return fmt.Errorf("product_code is required")
	}
	if row.ProductName == "" {
		return fmt.Errorf("product_name is required")
	}
	if row.ManufacturerName == "" {
		return fmt.Errorf("manufacturer_name is required")
	}
	if row.ModelNumber == "" {
		return fmt.Errorf("model_number is required")
	}
	if row.Category == "" {
		return fmt.Errorf("category is required")
	}

	// Validate category
	validCategories := map[string]bool{
		"MRI": true, "CT": true, "X-Ray": true, "Ultrasound": true,
		"Patient Monitor": true, "Ventilator": true, "Anesthesia": true,
		"Dialysis": true, "Infusion Pump": true, "Surgical": true,
		"Laboratory": true, "Other": true,
	}
	if !validCategories[row.Category] {
		return fmt.Errorf("invalid category: %s", row.Category)
	}

	return nil
}

// insertOrUpdateCatalog inserts or updates equipment catalog
func (h *CatalogBulkImportHandler) insertOrUpdateCatalog(ctx context.Context, tx pgx.Tx, row *EquipmentCatalogCSVRow, updateMode bool, createdBy string) (string, error) {
	// Parse numeric fields
	var basePrice, weightKg *float64
	var serviceInterval, lifespan *int

	if row.BasePrice != "" {
		if val, err := strconv.ParseFloat(row.BasePrice, 64); err == nil {
			basePrice = &val
		}
	}
	if row.WeightKg != "" {
		if val, err := strconv.ParseFloat(row.WeightKg, 64); err == nil {
			weightKg = &val
		}
	}
	if row.RecommendedServiceIntervalDays != "" {
		if val, err := strconv.Atoi(row.RecommendedServiceIntervalDays); err == nil {
			serviceInterval = &val
		}
	}
	if row.EstimatedLifespanYears != "" {
		if val, err := strconv.Atoi(row.EstimatedLifespanYears); err == nil {
			lifespan = &val
		}
	}

	// Check if exists
	var existingID string
	checkQuery := `SELECT id FROM equipment_catalog WHERE product_code = $1 LIMIT 1`
	err := tx.QueryRow(ctx, checkQuery, row.ProductCode).Scan(&existingID)

	if err == nil && existingID != "" {
		// Exists
		if !updateMode {
			return "", fmt.Errorf("equipment already exists: %s (use update_mode=true to update)", row.ProductCode)
		}

		// Update
		updateQuery := `
			UPDATE equipment_catalog 
			SET product_name = $1, manufacturer_name = $2, model_number = $3, 
			    category = $4, subcategory = $5, description = $6,
			    base_price = $7, currency = $8, weight_kg = $9,
			    recommended_service_interval_days = $10, estimated_lifespan_years = $11,
			    maintenance_complexity = $12, updated_at = CURRENT_TIMESTAMP
			WHERE id = $13
			RETURNING id
		`
		err = tx.QueryRow(ctx, updateQuery,
			row.ProductName, row.ManufacturerName, row.ModelNumber,
			row.Category, row.Subcategory, row.Description,
			basePrice, row.Currency, weightKg,
			serviceInterval, lifespan, row.MaintenanceComplexity,
			existingID,
		).Scan(&existingID)

		if err != nil {
			return "", fmt.Errorf("failed to update equipment: %w", err)
		}
		return existingID, nil
	}

	// Insert new
	catalogID := uuid.New().String()
	insertQuery := `
		INSERT INTO equipment_catalog (
			id, product_code, product_name, manufacturer_name, model_number,
			category, subcategory, description, base_price, currency, weight_kg,
			recommended_service_interval_days, estimated_lifespan_years,
			maintenance_complexity, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`

	err = tx.QueryRow(ctx, insertQuery,
		catalogID, row.ProductCode, row.ProductName, row.ManufacturerName, row.ModelNumber,
		row.Category, row.Subcategory, row.Description, basePrice, row.Currency, weightKg,
		serviceInterval, lifespan, row.MaintenanceComplexity, createdBy,
	).Scan(&catalogID)

	if err != nil {
		return "", fmt.Errorf("failed to insert equipment: %w", err)
	}

	return catalogID, nil
}
