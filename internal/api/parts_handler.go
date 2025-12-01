package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/aby-med/medical-platform/internal/parts"
	"github.com/gorilla/mux"
)

// PartsHandler handles parts recommendation HTTP requests
type PartsHandler struct {
	engine *parts.Engine
	db     *sql.DB
}

// NewPartsHandler creates a new parts handler
func NewPartsHandler(aiManager *ai.Manager, db *sql.DB) *PartsHandler {
	return &PartsHandler{
		engine: parts.NewEngine(aiManager, db),
		db:     db,
	}
}

// RegisterRoutes registers parts routes
func (h *PartsHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/parts/recommend", h.RecommendParts).Methods("POST")
	r.HandleFunc("/api/parts/recommendations/{requestId}", h.GetRecommendation).Methods("GET")
	r.HandleFunc("/api/tickets/{ticketId}/parts-recommendations", h.GetTicketRecommendations).Methods("GET")
	r.HandleFunc("/api/tickets/{ticketId}/parts", h.GetTicketParts).Methods("GET")
	r.HandleFunc("/api/parts/recommendations/{requestId}/feedback", h.ProvideFeedback).Methods("POST")
	r.HandleFunc("/api/parts/recommendations/{requestId}/usage", h.RecordPartsUsage).Methods("POST")
	r.HandleFunc("/api/parts/analytics", h.GetAnalytics).Methods("GET")
	r.HandleFunc("/api/parts/catalog", h.SearchCatalog).Methods("GET")
}

// RecommendParts handles POST /api/parts/recommend
func (h *PartsHandler) RecommendParts(w http.ResponseWriter, r *http.Request) {
	var req parts.RecommendationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Set default options if not provided
	if req.Options == (parts.RecommendationOptions{}) {
		req.Options = parts.DefaultRecommendationOptions()
	}

	// Get recommendations
	result, err := h.engine.RecommendParts(r.Context(), &req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Recommendation failed",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// GetRecommendation handles GET /api/parts/recommendations/{requestId}
func (h *PartsHandler) GetRecommendation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestId"]

	var result parts.RecommendationResponse

	query := `
		SELECT 
			request_id,
			ticket_id,
			replacement_parts,
			accessories,
			preventive_parts,
			metadata,
			created_at
		FROM parts_recommendations
		WHERE request_id = $1
	`

	var replacementJSON, accessoriesJSON, preventiveJSON, metadataJSON []byte
	err := h.db.QueryRowContext(r.Context(), query, requestID).Scan(
		&result.RequestID,
		&result.TicketID,
		&replacementJSON,
		&accessoriesJSON,
		&preventiveJSON,
		&metadataJSON,
		&result.CreatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   "Recommendation not found",
			Message: "No recommendation with this request ID",
		})
		return
	}

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve recommendation",
			Message: err.Error(),
		})
		return
	}

	json.Unmarshal(replacementJSON, &result.ReplacementParts)
	json.Unmarshal(accessoriesJSON, &result.Accessories)
	json.Unmarshal(preventiveJSON, &result.PreventiveParts)
	json.Unmarshal(metadataJSON, &result.Metadata)

	respondJSON(w, http.StatusOK, result)
}

// GetTicketRecommendations handles GET /api/tickets/{ticketId}/parts-recommendations
func (h *PartsHandler) GetTicketRecommendations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketIDStr := vars["ticketId"]

	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ticket ID",
			Message: err.Error(),
		})
		return
	}

	query := `
		SELECT 
			request_id,
			ticket_id,
			replacement_parts,
			accessories,
			preventive_parts,
			metadata,
			was_accurate,
			created_at
		FROM parts_recommendations
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`

    rows, err := h.db.QueryContext(r.Context(), query, ticketID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve recommendations",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var recommendations []map[string]interface{}
	for rows.Next() {
		var result parts.RecommendationResponse
		var replacementJSON, accessoriesJSON, preventiveJSON, metadataJSON []byte
        var wasAccurate sql.NullBool

        err := rows.Scan(
			&result.RequestID,
			&result.TicketID,
			&replacementJSON,
			&accessoriesJSON,
			&preventiveJSON,
			&metadataJSON,
			&wasAccurate,
			&result.CreatedAt,
		)
        if err != nil {
            continue
        }

		json.Unmarshal(replacementJSON, &result.ReplacementParts)
		json.Unmarshal(accessoriesJSON, &result.Accessories)
		json.Unmarshal(preventiveJSON, &result.PreventiveParts)
		json.Unmarshal(metadataJSON, &result.Metadata)

        // Handle nullable boolean
        var wasAccuratePtr *bool
        if wasAccurate.Valid {
            v := wasAccurate.Bool
            wasAccuratePtr = &v
        }

        recommendations = append(recommendations, map[string]interface{}{
			"request_id":        result.RequestID,
			"replacement_parts": result.ReplacementParts,
			"accessories":       result.Accessories,
			"preventive_parts":  result.PreventiveParts,
			"metadata":          result.Metadata,
            "was_accurate":      wasAccuratePtr,
			"created_at":        result.CreatedAt,
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_id":       ticketID,
		"count":           len(recommendations),
		"recommendations": recommendations,
	})
}

// PartsUsageRequest represents parts that were actually used
type PartsUsageRequest struct {
	PartsUsed       []PartUsage `json:"parts_used"`
	AccessoriesSold []int64     `json:"accessories_sold"`
}

// PartUsage represents a part that was used
type PartUsage struct {
	PartID   int64   `json:"part_id"`
	Quantity int     `json:"quantity"`
	Cost     float64 `json:"cost"`
	Notes    string  `json:"notes"`
}

// RecordPartsUsage handles POST /api/parts/recommendations/{requestId}/usage
func (h *PartsHandler) RecordPartsUsage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestId"]

	var req PartsUsageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Get ticket ID from recommendation
	var ticketID int64
	err := h.db.QueryRowContext(r.Context(), `
		SELECT ticket_id FROM parts_recommendations WHERE request_id = $1
	`, requestID).Scan(&ticketID)

	if err != nil {
		respondJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   "Recommendation not found",
			Message: err.Error(),
		})
		return
	}

	// Record parts usage in ticket_parts
	for _, partUsage := range req.PartsUsed {
		_, err := h.db.ExecContext(r.Context(), `
			INSERT INTO ticket_parts (ticket_id, part_id, quantity_used, was_recommended, cost, notes)
			VALUES ($1, $2, $3, true, $4, $5)
			ON CONFLICT (ticket_id, part_id) DO UPDATE
			SET quantity_used = ticket_parts.quantity_used + $3,
				cost = $4,
				notes = $5
		`, ticketID, partUsage.PartID, partUsage.Quantity, partUsage.Cost, partUsage.Notes)

		if err != nil {
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to record part usage",
				Message: err.Error(),
			})
			return
		}
	}

	// Update recommendation with usage
	partsUsedIDs := make([]int64, len(req.PartsUsed))
	for i, p := range req.PartsUsed {
		partsUsedIDs[i] = p.PartID
	}

	_, err = h.db.ExecContext(r.Context(), `
		UPDATE parts_recommendations
		SET parts_used = $1,
			accessories_sold = $2,
			updated_at = NOW()
		WHERE request_id = $3
	`, partsUsedIDs, req.AccessoriesSold, requestID)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update recommendation",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parts usage recorded successfully",
	})
}

// PartsFeedbackRequest represents feedback on recommendation
type PartsFeedbackRequest struct {
	WasAccurate      bool   `json:"was_accurate"`
	AccuracyFeedback string `json:"accuracy_feedback"`
}

// ProvideFeedback handles POST /api/parts/recommendations/{requestId}/feedback
func (h *PartsHandler) ProvideFeedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestId"]

	var req PartsFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	query := `
		UPDATE parts_recommendations
		SET was_accurate = $1,
			accuracy_feedback = $2,
			updated_at = NOW()
		WHERE request_id = $3
	`

	result, err := h.db.ExecContext(r.Context(), query, req.WasAccurate, req.AccuracyFeedback, requestID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to save feedback",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   "Recommendation not found",
			Message: "No recommendation with this request ID",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Feedback recorded successfully",
	})
}

// GetAnalytics handles GET /api/parts/analytics
func (h *PartsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// Get recommendation accuracy
    accuracyQuery := `
        SELECT 
            recommendation_date,
            total_recommendations,
            accurate_count,
            ROUND(accuracy_rate, 2) as accuracy_rate,
            pending_feedback,
            ROUND(total_ai_cost, 4) as total_ai_cost,
            ai_assisted_count
        FROM v_parts_recommendation_accuracy
        WHERE recommendation_date > NOW() - make_interval(days => $1)
        ORDER BY recommendation_date DESC
    `

	rows, err := h.db.QueryContext(r.Context(), accuracyQuery, days)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve analytics",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

    var dailyAnalytics []map[string]interface{}
    for rows.Next() {
        var (
            date            string
            total           int
            accurate        int
            accuracyRate    sql.NullFloat64
            pending         int
            cost            sql.NullFloat64
            aiAssistedCount int
        )

        if err := rows.Scan(&date, &total, &accurate, &accuracyRate, &pending, &cost, &aiAssistedCount); err != nil {
            continue
        }

        var accuracyPtr *float64
        if accuracyRate.Valid {
            v := accuracyRate.Float64
            accuracyPtr = &v
        }

        var costPtr *float64
        if cost.Valid {
            v := cost.Float64
            costPtr = &v
        }

        dailyAnalytics = append(dailyAnalytics, map[string]interface{}{
            "date":              date,
            "total":             total,
            "accurate":          accurate,
            "accuracy_rate":     accuracyPtr,
            "pending_feedback":  pending,
            "ai_cost":           costPtr,
            "ai_assisted_count": aiAssistedCount,
        })
    }

	// Get top parts usage
	partsQuery := `
		SELECT 
			part_number,
			part_name,
			category,
			times_used,
			times_recommended,
			avg_cost
		FROM v_parts_usage_analysis
		ORDER BY times_used DESC
		LIMIT 10
	`

    partsRows, err := h.db.QueryContext(r.Context(), partsQuery)
	if err == nil {
		defer partsRows.Close()
		var topParts []map[string]interface{}

		for partsRows.Next() {
			var (
				partNumber       string
				partName         string
				category         string
				timesUsed        int
				timesRecommended int
                avgCost          sql.NullFloat64
			)

            if err := partsRows.Scan(&partNumber, &partName, &category, &timesUsed, &timesRecommended, &avgCost); err != nil {
                continue
            }

            var avgCostPtr *float64
            if avgCost.Valid {
                v := avgCost.Float64
                avgCostPtr = &v
            }

            topParts = append(topParts, map[string]interface{}{
				"part_number":        partNumber,
				"part_name":          partName,
				"category":           category,
				"times_used":         timesUsed,
				"times_recommended":  timesRecommended,
                "avg_cost":           avgCostPtr,
			})
		}

		// Get accessory sales
		accessoryQuery := `
			SELECT 
				part_number,
				part_name,
				times_sold,
				unit_price,
				total_revenue
			FROM v_accessory_sales_analysis
			ORDER BY total_revenue DESC
			LIMIT 10
		`

        accessoryRows, err := h.db.QueryContext(r.Context(), accessoryQuery)
		var topAccessories []map[string]interface{}

		if err == nil {
			defer accessoryRows.Close()

			for accessoryRows.Next() {
				var (
					partNumber  string
					partName    string
					timesSold   int
					unitPrice   float64
					totalRevenue float64
				)

                if err := accessoryRows.Scan(&partNumber, &partName, &timesSold, &unitPrice, &totalRevenue); err != nil {
                    continue
                }

                topAccessories = append(topAccessories, map[string]interface{}{
					"part_number":   partNumber,
					"part_name":     partName,
					"times_sold":    timesSold,
					"unit_price":    unitPrice,
					"total_revenue": totalRevenue,
				})
			}
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"period_days":      days,
			"daily_analytics":  dailyAnalytics,
			"top_parts":        topParts,
			"top_accessories":  topAccessories,
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"period_days":     days,
		"daily_analytics": dailyAnalytics,
	})
}

// SearchCatalog handles GET /api/parts/catalog
func (h *PartsHandler) SearchCatalog(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	equipmentTypeIDStr := r.URL.Query().Get("equipment_type_id")

	if query == "" && category == "" && equipmentTypeIDStr == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Missing search criteria",
			Message: "Provide at least one of: q, category, or equipment_type_id",
		})
		return
	}

	sqlQuery := `
		SELECT DISTINCT
			pc.part_id,
			pc.part_number,
			pc.part_name,
			pc.description,
			pc.category,
			pc.subcategory,
			pc.part_type,
			pc.is_oem_part,
			pc.unit_price,
			pc.manufacturer_name,
			pi.quantity_available
		FROM parts_catalog pc
		LEFT JOIN parts_inventory pi ON pc.part_id = pi.part_id
		LEFT JOIN equipment_parts ep ON pc.part_id = ep.part_id
		WHERE pc.is_active = true
	`

	args := []interface{}{}
	argCount := 1

	if query != "" {
		sqlQuery += ` AND (pc.part_name ILIKE $` + strconv.Itoa(argCount) + 
			` OR pc.part_number ILIKE $` + strconv.Itoa(argCount) + 
			` OR pc.description ILIKE $` + strconv.Itoa(argCount) + `)`
		args = append(args, "%"+query+"%")
		argCount++
	}

	if category != "" {
		sqlQuery += ` AND pc.category ILIKE $` + strconv.Itoa(argCount)
		args = append(args, "%"+category+"%")
		argCount++
	}

	if equipmentTypeIDStr != "" {
		equipmentTypeID, err := strconv.ParseInt(equipmentTypeIDStr, 10, 64)
		if err == nil {
			sqlQuery += ` AND ep.equipment_type_id = $` + strconv.Itoa(argCount)
			args = append(args, equipmentTypeID)
		}
	}

	sqlQuery += ` ORDER BY pc.part_name LIMIT 50`

	rows, err := h.db.QueryContext(r.Context(), sqlQuery, args...)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Search failed",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

    var results []map[string]interface{}
    for rows.Next() {
        var (
            partID       int64
            partNumber   string
            partName     string
            description  string
            category     string
            subcategory  sql.NullString
            partType     string
            isOEM        bool
            unitPrice    sql.NullFloat64
            manufacturer sql.NullString
            qtyAvailable sql.NullInt64
        )

        if err := rows.Scan(&partID, &partNumber, &partName, &description, &category, &subcategory,
            &partType, &isOEM, &unitPrice, &manufacturer, &qtyAvailable); err != nil {
            continue
        }

        result := map[string]interface{}{
            "part_id":     partID,
            "part_number": partNumber,
            "part_name":   partName,
            "description": description,
            "category":    category,
            "part_type":   partType,
            "is_oem":      isOEM,
        }

        if subcategory.Valid {
            result["subcategory"] = subcategory.String
        }
        if unitPrice.Valid {
            result["unit_price"] = unitPrice.Float64
        }
        if manufacturer.Valid {
            result["manufacturer"] = manufacturer.String
        }
        if qtyAvailable.Valid {
            result["stock_available"] = qtyAvailable.Int64
        }

        results = append(results, result)
    }

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(results),
		"results": results,
	})
}

// GetTicketParts handles GET /api/tickets/{ticketId}/parts
func (h *PartsHandler) GetTicketParts(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    ticketID := vars["ticketId"]

    const q = `
        SELECT spare_part_id, part_number, part_name, unit_price, currency,
               is_critical, quantity_required, part_category, stock_status, lead_time_days
        FROM get_parts_for_ticket($1)
    `

    rows, err := h.db.QueryContext(r.Context(), q, ticketID)
    if err != nil {
        respondJSON(w, http.StatusInternalServerError, ErrorResponse{
            Error:   "Failed to fetch ticket parts",
            Message: err.Error(),
        })
        return
    }
    defer rows.Close()

    var parts []map[string]interface{}
    for rows.Next() {
        var (
            partID       string
            partNumber   string
            partName     string
            unitPrice    sql.NullFloat64
            currency     sql.NullString
            isCritical   sql.NullBool
            qtyRequired  sql.NullInt64
            category     sql.NullString
            stockStatus  sql.NullString
            leadTimeDays sql.NullInt64
        )

        if err := rows.Scan(&partID, &partNumber, &partName, &unitPrice, &currency,
            &isCritical, &qtyRequired, &category, &stockStatus, &leadTimeDays); err != nil {
            continue
        }

        item := map[string]interface{}{
            "spare_part_id": partID,
            "part_number":   partNumber,
            "part_name":     partName,
        }
        if unitPrice.Valid {
            item["unit_price"] = unitPrice.Float64
        }
        if currency.Valid {
            item["currency"] = currency.String
        }
        if isCritical.Valid {
            item["is_critical"] = isCritical.Bool
        }
        if qtyRequired.Valid {
            item["quantity_required"] = qtyRequired.Int64
        }
        if category.Valid {
            item["category"] = category.String
        }
        if stockStatus.Valid {
            item["stock_status"] = stockStatus.String
        }
        if leadTimeDays.Valid {
            item["lead_time_days"] = leadTimeDays.Int64
        }

        parts = append(parts, item)
    }

    respondJSON(w, http.StatusOK, map[string]interface{}{
        "ticket_id": ticketID,
        "count":     len(parts),
        "parts":     parts,
    })
}

