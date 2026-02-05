package api

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// EngineerCSVRow represents a row from engineer CSV import
type EngineerCSVRow struct {
	Name            string
	Phone           string
	Email           string
	Location        string
	EngineerLevel   int
	EquipmentTypes  []string
	ExperienceYears int
}

// EngineerImportResult represents the result of CSV import
type EngineerImportResult struct {
	TotalRows    int      `json:"total_rows"`
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors,omitempty"`
	ImportedIDs  []string `json:"imported_ids,omitempty"`
}

// ImportEngineersCSV handles POST /engineers/import
func (h *AssignmentHandler) ImportEngineersCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the CSV file
	file, _, err := r.FormFile("csv_file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Process the CSV
	result, err := h.processEngineerCSV(ctx, file)
	if err != nil {
		h.logger.Error("Failed to process engineer CSV", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to process CSV: %v", err), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *AssignmentHandler) processEngineerCSV(ctx context.Context, file io.Reader) (*EngineerImportResult, error) {
	result := &EngineerImportResult{
		Errors:      []string{},
		ImportedIDs: []string{},
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Validate header
	expectedHeaders := []string{"name", "phone", "email", "location", "engineer_level", "equipment_types", "experience_years"}
	if len(header) < len(expectedHeaders) {
		return nil, fmt.Errorf("invalid header: expected at least %d columns, got %d", len(expectedHeaders), len(header))
	}

	rowNum := 1 // Start from 1 (header is row 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to read: %v", rowNum, err))
			result.FailureCount++
			rowNum++
			continue
		}

		result.TotalRows++

		// Parse row
		if len(record) < len(expectedHeaders) {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Insufficient columns", rowNum))
			result.FailureCount++
			rowNum++
			continue
		}

		engineerLevel, err := strconv.Atoi(strings.TrimSpace(record[4]))
		if err != nil || engineerLevel < 1 || engineerLevel > 3 {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Invalid engineer_level (must be 1, 2, or 3)", rowNum))
			result.FailureCount++
			rowNum++
			continue
		}

		equipmentTypes := strings.Split(strings.TrimSpace(record[5]), "|")
		for i := range equipmentTypes {
			equipmentTypes[i] = strings.TrimSpace(equipmentTypes[i])
		}

		experienceYears := 0
		if len(record) > 6 && record[6] != "" {
			experienceYears, _ = strconv.Atoi(strings.TrimSpace(record[6]))
		}

		csvRow := EngineerCSVRow{
			Name:            strings.TrimSpace(record[0]),
			Phone:           strings.TrimSpace(record[1]),
			Email:           strings.TrimSpace(record[2]),
			Location:        strings.TrimSpace(record[3]),
			EngineerLevel:   engineerLevel,
			EquipmentTypes:  equipmentTypes,
			ExperienceYears: experienceYears,
		}

		// Validate required fields
		if csvRow.Name == "" || csvRow.Phone == "" || csvRow.Email == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Missing required fields (name, phone, or email)", rowNum))
			result.FailureCount++
			rowNum++
			continue
		}

		// Create engineer in database
		engineerID, err := h.createEngineerFromCSV(ctx, csvRow)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to create engineer: %v", rowNum, err))
			result.FailureCount++
			rowNum++
			continue
		}

		result.ImportedIDs = append(result.ImportedIDs, engineerID)
		result.SuccessCount++
		rowNum++
	}

	h.logger.Info("Engineer CSV import completed",
		slog.Int("total", result.TotalRows),
		slog.Int("success", result.SuccessCount),
		slog.Int("failure", result.FailureCount))

	return result, nil
}

func (h *AssignmentHandler) createEngineerFromCSV(ctx context.Context, row EngineerCSVRow) (string, error) {
	// Get organization ID from context
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok || orgID == "" {
		return "", fmt.Errorf("organization_id not found in context")
	}
	
	h.logger.Info("Creating engineer from CSV",
		slog.String("name", row.Name),
		slog.String("email", row.Email),
		slog.Int("engineer_level", row.EngineerLevel))
	
	// Create engineer through service layer
	engineerID, err := h.service.CreateEngineer(
		ctx,
		row.Name,
		row.Phone,
		row.Email,
		row.Location,
		row.EngineerLevel,
		row.EquipmentTypes,
		row.ExperienceYears,
		orgID,
	)
	
	if err != nil {
		return "", fmt.Errorf("failed to create engineer: %w", err)
	}
	
	h.logger.Info("Engineer created successfully",
		slog.String("engineer_id", engineerID),
		slog.String("name", row.Name))
	
	return engineerID, nil
}
