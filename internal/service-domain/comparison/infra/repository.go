package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aby-med/medical-platform/internal/service-domain/comparison/domain"
	"github.com/jackc/pgx/v5"
)

// ComparisonRepository implements the domain.Repository interface
type ComparisonRepository struct {
	db     *PostgresDB
	logger *slog.Logger
}

// NewComparisonRepository creates a new comparison repository
func NewComparisonRepository(db *PostgresDB, logger *slog.Logger) *ComparisonRepository {
	return &ComparisonRepository{
		db:     db,
		logger: logger.With(slog.String("component", "comparison_repository")),
	}
}

// Create creates a new comparison
func (r *ComparisonRepository) Create(ctx context.Context, comparison *domain.Comparison) error {
	scoringCriteria, err := json.Marshal(comparison.ScoringCriteria)
	if err != nil {
		return fmt.Errorf("failed to marshal scoring criteria: %w", err)
	}

	quoteScores, err := json.Marshal(comparison.QuoteScores)
	if err != nil {
		return fmt.Errorf("failed to marshal quote scores: %w", err)
	}

	priceDifferences, err := json.Marshal(comparison.PriceDifferences)
	if err != nil {
		return fmt.Errorf("failed to marshal price differences: %w", err)
	}

	itemComparisons, err := json.Marshal(comparison.ItemComparisons)
	if err != nil {
		return fmt.Errorf("failed to marshal item comparisons: %w", err)
	}

	query := `
		INSERT INTO comparisons (
			id, tenant_id, rfq_id, title, description, quote_ids, status,
			scoring_criteria, quote_scores, price_differences, item_comparisons,
			best_overall_quote, best_price_quote, recommendation, notes,
			created_by, created_at, updated_at, completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	_, err = r.db.Pool().Exec(ctx, query,
		comparison.ID,
		comparison.TenantID,
		comparison.RFQID,
		comparison.Title,
		comparison.Description,
		comparison.QuoteIDs,
		comparison.Status,
		scoringCriteria,
		quoteScores,
		priceDifferences,
		itemComparisons,
		comparison.BestOverallQuote,
		comparison.BestPriceQuote,
		comparison.Recommendation,
		comparison.Notes,
		comparison.CreatedBy,
		comparison.CreatedAt,
		comparison.UpdatedAt,
		comparison.CompletedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create comparison", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create comparison: %w", err)
	}

	return nil
}

// GetByID retrieves a comparison by ID
func (r *ComparisonRepository) GetByID(ctx context.Context, tenantID, id string) (*domain.Comparison, error) {
	query := `
		SELECT id, tenant_id, rfq_id, title, description, quote_ids, status,
			   scoring_criteria, quote_scores, price_differences, item_comparisons,
			   best_overall_quote, best_price_quote, recommendation, notes,
			   created_by, created_at, updated_at, completed_at
		FROM comparisons
		WHERE tenant_id = $1 AND id = $2
	`

	row := r.db.Pool().QueryRow(ctx, query, tenantID, id)

	comparison := &domain.Comparison{}
	var scoringCriteriaJSON, quoteScoresJSON, priceDifferencesJSON, itemComparisonsJSON []byte

	err := row.Scan(
		&comparison.ID,
		&comparison.TenantID,
		&comparison.RFQID,
		&comparison.Title,
		&comparison.Description,
		&comparison.QuoteIDs,
		&comparison.Status,
		&scoringCriteriaJSON,
		&quoteScoresJSON,
		&priceDifferencesJSON,
		&itemComparisonsJSON,
		&comparison.BestOverallQuote,
		&comparison.BestPriceQuote,
		&comparison.Recommendation,
		&comparison.Notes,
		&comparison.CreatedBy,
		&comparison.CreatedAt,
		&comparison.UpdatedAt,
		&comparison.CompletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrComparisonNotFound
		}
		r.logger.Error("Failed to get comparison", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get comparison: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(scoringCriteriaJSON, &comparison.ScoringCriteria); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scoring criteria: %w", err)
	}
	if err := json.Unmarshal(quoteScoresJSON, &comparison.QuoteScores); err != nil {
		return nil, fmt.Errorf("failed to unmarshal quote scores: %w", err)
	}
	if err := json.Unmarshal(priceDifferencesJSON, &comparison.PriceDifferences); err != nil {
		return nil, fmt.Errorf("failed to unmarshal price differences: %w", err)
	}
	if err := json.Unmarshal(itemComparisonsJSON, &comparison.ItemComparisons); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item comparisons: %w", err)
	}

	return comparison, nil
}

// GetByRFQ retrieves all comparisons for an RFQ
func (r *ComparisonRepository) GetByRFQ(ctx context.Context, tenantID, rfqID string) ([]*domain.Comparison, error) {
	query := `
		SELECT id, tenant_id, rfq_id, title, description, quote_ids, status,
			   scoring_criteria, quote_scores, price_differences, item_comparisons,
			   best_overall_quote, best_price_quote, recommendation, notes,
			   created_by, created_at, updated_at, completed_at
		FROM comparisons
		WHERE tenant_id = $1 AND rfq_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool().Query(ctx, query, tenantID, rfqID)
	if err != nil {
		r.logger.Error("Failed to get comparisons by RFQ", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get comparisons by RFQ: %w", err)
	}
	defer rows.Close()

	comparisons := []*domain.Comparison{}
	for rows.Next() {
		comparison := &domain.Comparison{}
		var scoringCriteriaJSON, quoteScoresJSON, priceDifferencesJSON, itemComparisonsJSON []byte

		err := rows.Scan(
			&comparison.ID,
			&comparison.TenantID,
			&comparison.RFQID,
			&comparison.Title,
			&comparison.Description,
			&comparison.QuoteIDs,
			&comparison.Status,
			&scoringCriteriaJSON,
			&quoteScoresJSON,
			&priceDifferencesJSON,
			&itemComparisonsJSON,
			&comparison.BestOverallQuote,
			&comparison.BestPriceQuote,
			&comparison.Recommendation,
			&comparison.Notes,
			&comparison.CreatedBy,
			&comparison.CreatedAt,
			&comparison.UpdatedAt,
			&comparison.CompletedAt,
		)

		if err != nil {
			r.logger.Error("Failed to scan comparison", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan comparison: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(scoringCriteriaJSON, &comparison.ScoringCriteria); err != nil {
			return nil, fmt.Errorf("failed to unmarshal scoring criteria: %w", err)
		}
		if err := json.Unmarshal(quoteScoresJSON, &comparison.QuoteScores); err != nil {
			return nil, fmt.Errorf("failed to unmarshal quote scores: %w", err)
		}
		if err := json.Unmarshal(priceDifferencesJSON, &comparison.PriceDifferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal price differences: %w", err)
		}
		if err := json.Unmarshal(itemComparisonsJSON, &comparison.ItemComparisons); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item comparisons: %w", err)
		}

		comparisons = append(comparisons, comparison)
	}

	return comparisons, nil
}

// List retrieves comparisons with filtering
func (r *ComparisonRepository) List(ctx context.Context, criteria domain.ListCriteria) (*domain.ListResult, error) {
	// Build query with filters
	var conditions []string
	var args []interface{}
	argCount := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argCount))
	args = append(args, criteria.TenantID)
	argCount++

	if criteria.RFQID != "" {
		conditions = append(conditions, fmt.Sprintf("rfq_id = $%d", argCount))
		args = append(args, criteria.RFQID)
		argCount++
	}

	if len(criteria.Status) > 0 {
		statusPlaceholders := []string{}
		for _, status := range criteria.Status {
			statusPlaceholders = append(statusPlaceholders, fmt.Sprintf("$%d", argCount))
			args = append(args, status)
			argCount++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(statusPlaceholders, ",")))
	}

	if criteria.CreatedBy != "" {
		conditions = append(conditions, fmt.Sprintf("created_by = $%d", argCount))
		args = append(args, criteria.CreatedBy)
		argCount++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM comparisons %s", whereClause)
	var total int
	if err := r.db.Pool().QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count comparisons: %w", err)
	}

	// Build ORDER BY
	sortBy := "created_at"
	if criteria.SortBy != "" {
		sortBy = criteria.SortBy
	}
	sortDirection := "DESC"
	if criteria.SortDirection != "" {
		sortDirection = strings.ToUpper(criteria.SortDirection)
	}

	// Pagination
	page := criteria.Page
	if page < 1 {
		page = 1
	}
	pageSize := criteria.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(`
		SELECT id, tenant_id, rfq_id, title, description, quote_ids, status,
			   scoring_criteria, quote_scores, price_differences, item_comparisons,
			   best_overall_quote, best_price_quote, recommendation, notes,
			   created_by, created_at, updated_at, completed_at
		FROM comparisons
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, sortDirection, argCount, argCount+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.Pool().Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to list comparisons", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to list comparisons: %w", err)
	}
	defer rows.Close()

	comparisons := []*domain.Comparison{}
	for rows.Next() {
		comparison := &domain.Comparison{}
		var scoringCriteriaJSON, quoteScoresJSON, priceDifferencesJSON, itemComparisonsJSON []byte

		err := rows.Scan(
			&comparison.ID,
			&comparison.TenantID,
			&comparison.RFQID,
			&comparison.Title,
			&comparison.Description,
			&comparison.QuoteIDs,
			&comparison.Status,
			&scoringCriteriaJSON,
			&quoteScoresJSON,
			&priceDifferencesJSON,
			&itemComparisonsJSON,
			&comparison.BestOverallQuote,
			&comparison.BestPriceQuote,
			&comparison.Recommendation,
			&comparison.Notes,
			&comparison.CreatedBy,
			&comparison.CreatedAt,
			&comparison.UpdatedAt,
			&comparison.CompletedAt,
		)

		if err != nil {
			r.logger.Error("Failed to scan comparison", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan comparison: %w", err)
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(scoringCriteriaJSON, &comparison.ScoringCriteria); err != nil {
			return nil, fmt.Errorf("failed to unmarshal scoring criteria: %w", err)
		}
		if err := json.Unmarshal(quoteScoresJSON, &comparison.QuoteScores); err != nil {
			return nil, fmt.Errorf("failed to unmarshal quote scores: %w", err)
		}
		if err := json.Unmarshal(priceDifferencesJSON, &comparison.PriceDifferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal price differences: %w", err)
		}
		if err := json.Unmarshal(itemComparisonsJSON, &comparison.ItemComparisons); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item comparisons: %w", err)
		}

		comparisons = append(comparisons, comparison)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &domain.ListResult{
		Comparisons: comparisons,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

// Update updates a comparison
func (r *ComparisonRepository) Update(ctx context.Context, comparison *domain.Comparison) error {
	scoringCriteria, err := json.Marshal(comparison.ScoringCriteria)
	if err != nil {
		return fmt.Errorf("failed to marshal scoring criteria: %w", err)
	}

	quoteScores, err := json.Marshal(comparison.QuoteScores)
	if err != nil {
		return fmt.Errorf("failed to marshal quote scores: %w", err)
	}

	priceDifferences, err := json.Marshal(comparison.PriceDifferences)
	if err != nil {
		return fmt.Errorf("failed to marshal price differences: %w", err)
	}

	itemComparisons, err := json.Marshal(comparison.ItemComparisons)
	if err != nil {
		return fmt.Errorf("failed to marshal item comparisons: %w", err)
	}

	query := `
		UPDATE comparisons
		SET title = $3, description = $4, quote_ids = $5, status = $6,
			scoring_criteria = $7, quote_scores = $8, price_differences = $9, item_comparisons = $10,
			best_overall_quote = $11, best_price_quote = $12, recommendation = $13, notes = $14,
			updated_at = $15, completed_at = $16
		WHERE tenant_id = $1 AND id = $2
	`

	result, err := r.db.Pool().Exec(ctx, query,
		comparison.TenantID,
		comparison.ID,
		comparison.Title,
		comparison.Description,
		comparison.QuoteIDs,
		comparison.Status,
		scoringCriteria,
		quoteScores,
		priceDifferences,
		itemComparisons,
		comparison.BestOverallQuote,
		comparison.BestPriceQuote,
		comparison.Recommendation,
		comparison.Notes,
		comparison.UpdatedAt,
		comparison.CompletedAt,
	)

	if err != nil {
		r.logger.Error("Failed to update comparison", slog.String("error", err.Error()))
		return fmt.Errorf("failed to update comparison: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrComparisonNotFound
	}

	return nil
}

// Delete deletes a comparison
func (r *ComparisonRepository) Delete(ctx context.Context, tenantID, id string) error {
	query := "DELETE FROM comparisons WHERE tenant_id = $1 AND id = $2"

	result, err := r.db.Pool().Exec(ctx, query, tenantID, id)
	if err != nil {
		r.logger.Error("Failed to delete comparison", slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete comparison: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrComparisonNotFound
	}

	return nil
}
