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
	"github.com/jackc/pgx/v5/pgxpool"
)

// PartsImportHandler handles parts CSV import
type PartsImportHandler struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewPartsImportHandler creates a new parts import handler
func NewPartsImportHandler(pool *pgxpool.Pool, logger *slog.Logger) *PartsImportHandler {
	return &PartsImportHandler{
		pool:   pool,
		logger: logger.With(slog.String("component", "parts_import_handler")),
	}
}

// PartCSVRow represents a row from parts CSV import
type PartCSVRow struct {
	PartNumber       string
	PartName         string
	Description      string
	Category         string
	Subcategory      string
	PartType         string
	IsOEMPart        bool
	ManufacturerName string
	UnitPrice        float64
	Currency         string
	MinimumStock     int
	LeadTimeDays     int
	WeightKg         float64
	Dimensions       string
	WarrantyMonths   int
	Specifications   string
}

// PartsImportResult represents the result of parts import
type PartsImportResult struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors,omitempty"`
	ImportedIDs  []string `json:"imported_ids,omitempty"`
}

// ImportParts handles POST /api/v1/parts/import
func (h *PartsImportHandler) ImportParts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("csv_file")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "CSV file is required")
		return
	}
	defer file.Close()

	h.logger.Info("Starting parts import", slog.String("filename", header.Filename))

	// Parse CSV
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	headerRow, err := reader.Read()
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to read CSV header: "+err.Error())
		return
	}

	h.logger.Info("CSV header", slog.Any("columns", headerRow))

	result := &PartsImportResult{
		Errors:      []string{},
		ImportedIDs: []string{},
	}

	rowNum := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to read: %v", rowNum, err))
			result.FailureCount++
			rowNum++
			continue
		}

		// Parse row
		if len(row) < 16 {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Insufficient columns (expected 16, got %d)", rowNum, len(row)))
			result.FailureCount++
			rowNum++
			continue
		}

		// Create part from CSV row
		partID, err := h.createPartFromCSV(ctx, row)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %v", rowNum, err))
			result.FailureCount++
			rowNum++
			continue
		}

		result.ImportedIDs = append(result.ImportedIDs, partID)
		result.SuccessCount++
		rowNum++
	}

	h.logger.Info("Parts import complete",
		slog.Int("success", result.SuccessCount),
		slog.Int("failures", result.FailureCount))

	h.respondJSON(w, http.StatusOK, result)
}

func (h *PartsImportHandler) createPartFromCSV(ctx context.Context, row []string) (string, error) {
	// Parse CSV row
	csvRow, err := h.parseCSVRow(row)
	if err != nil {
		return "", err
	}

	partID := uuid.New().String()

	// Insert into spare_parts_catalog table
	query := `
		INSERT INTO spare_parts_catalog (
			id, part_number, part_name, description, category, subcategory,
			part_type, manufacturer_part_number, manufacturer_name,
			unit_price, currency, is_available, stock_status,
			lead_time_days, minimum_order_quantity, is_oem_part,
			weight_kg, dimensions_cm, warranty_months, specifications,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $2, $8,
			$9, $10, true, 'in_stock',
			$11, $12, $13,
			$14, $15, $16, $17,
			NOW(), NOW()
		)
		ON CONFLICT (part_number) DO UPDATE SET
			part_name = EXCLUDED.part_name,
			description = EXCLUDED.description,
			unit_price = EXCLUDED.unit_price,
			updated_at = NOW()
		RETURNING id
	`

	err = h.pool.QueryRow(ctx, query,
		partID,
		csvRow.PartNumber,
		csvRow.PartName,
		csvRow.Description,
		csvRow.Category,
		csvRow.Subcategory,
		csvRow.PartType,
		csvRow.ManufacturerName,
		csvRow.UnitPrice,
		csvRow.Currency,
		csvRow.LeadTimeDays,
		csvRow.MinimumStock,
		csvRow.IsOEMPart,
		csvRow.WeightKg,
		csvRow.Dimensions,
		csvRow.WarrantyMonths,
		csvRow.Specifications,
	).Scan(&partID)

	if err != nil {
		return "", fmt.Errorf("failed to insert part: %w", err)
	}

	h.logger.Info("Part created",
		slog.String("part_id", partID),
		slog.String("part_number", csvRow.PartNumber),
		slog.String("part_name", csvRow.PartName))

	return partID, nil
}

func (h *PartsImportHandler) parseCSVRow(row []string) (*PartCSVRow, error) {
	// Validate required fields
	if row[0] == "" {
		return nil, fmt.Errorf("part_number is required")
	}
	if row[1] == "" {
		return nil, fmt.Errorf("part_name is required")
	}

	// Parse numeric fields
	unitPrice, err := strconv.ParseFloat(row[8], 64)
	if err != nil {
		unitPrice = 0
	}

	minimumStock, err := strconv.Atoi(row[10])
	if err != nil {
		minimumStock = 0
	}

	leadTimeDays, err := strconv.Atoi(row[11])
	if err != nil {
		leadTimeDays = 0
	}

	weightKg, err := strconv.ParseFloat(row[12], 64)
	if err != nil {
		weightKg = 0
	}

	warrantyMonths, err := strconv.Atoi(row[14])
	if err != nil {
		warrantyMonths = 0
	}

	// Parse boolean
	isOEMPart := false
	if strings.ToLower(row[6]) == "true" || row[6] == "1" {
		isOEMPart = true
	}

	return &PartCSVRow{
		PartNumber:       row[0],
		PartName:         row[1],
		Description:      row[2],
		Category:         row[3],
		Subcategory:      row[4],
		PartType:         row[5],
		IsOEMPart:        isOEMPart,
		ManufacturerName: row[7],
		UnitPrice:        unitPrice,
		Currency:         row[9],
		MinimumStock:     minimumStock,
		LeadTimeDays:     leadTimeDays,
		WeightKg:         weightKg,
		Dimensions:       row[13],
		WarrantyMonths:   warrantyMonths,
		Specifications:   row[15],
	}, nil
}

func (h *PartsImportHandler) respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *PartsImportHandler) respondError(w http.ResponseWriter, code int, message string) {
	h.logger.Error("Parts import error", slog.String("message", message))
	h.respondJSON(w, code, ErrorResponse{
		Error:   "Import failed",
		Message: message,
	})
}
