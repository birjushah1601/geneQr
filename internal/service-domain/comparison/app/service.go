package app

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/comparison/domain"
	"github.com/segmentio/ksuid"
)

// ComparisonService handles comparison business logic
type ComparisonService struct {
	repo   domain.Repository
	logger *slog.Logger
}

// NewComparisonService creates a new comparison service
func NewComparisonService(repo domain.Repository, logger *slog.Logger) *ComparisonService {
	return &ComparisonService{
		repo:   repo,
		logger: logger.With(slog.String("service", "comparison")),
	}
}

// CreateComparison creates a new comparison
func (s *ComparisonService) CreateComparison(ctx context.Context, tenantID, createdBy string, req CreateComparisonRequest) (*domain.Comparison, error) {
	comparison, err := domain.NewComparison(tenantID, req.RFQID, req.Title, createdBy, req.QuoteIDs)
	if err != nil {
		return nil, err
	}

	comparison.ID = ksuid.New().String()
	comparison.Description = req.Description

	if err := s.repo.Create(ctx, comparison); err != nil {
		return nil, err
	}

	s.logger.Info("Comparison created", slog.String("id", comparison.ID))
	return comparison, nil
}

// GetComparison retrieves a comparison by ID
func (s *ComparisonService) GetComparison(ctx context.Context, tenantID, id string) (*domain.Comparison, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetComparison sByRFQ retrieves all comparisons for an RFQ
func (s *ComparisonService) GetComparisonsByRFQ(ctx context.Context, tenantID, rfqID string) ([]*domain.Comparison, error) {
	return s.repo.GetByRFQ(ctx, tenantID, rfqID)
}

// ListComparisons lists comparisons with filtering
func (s *ComparisonService) ListComparisons(ctx context.Context, criteria domain.ListCriteria) (*domain.ListResult, error) {
	return s.repo.List(ctx, criteria)
}

// UpdateComparison updates a comparison
func (s *ComparisonService) UpdateComparison(ctx context.Context, tenantID, id string, req UpdateComparisonRequest) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		comparison.Title = *req.Title
	}
	if req.Description != nil {
		comparison.Description = *req.Description
	}
	if req.Notes != nil {
		comparison.Notes = *req.Notes
	}

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	return comparison, nil
}

// UpdateScoringCriteria updates the scoring weights
func (s *ComparisonService) UpdateScoringCriteria(ctx context.Context, tenantID, id string, req UpdateScoringCriteriaRequest) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	criteria := domain.ScoringCriteria{
		PriceWeight:      req.PriceWeight,
		QualityWeight:    req.QualityWeight,
		DeliveryWeight:   req.DeliveryWeight,
		ComplianceWeight: req.ComplianceWeight,
	}

	if err := comparison.UpdateScoringCriteria(criteria); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	return comparison, nil
}

// AddQuote adds a quote to the comparison
func (s *ComparisonService) AddQuote(ctx context.Context, tenantID, id string, req AddQuoteRequest) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := comparison.AddQuote(req.QuoteID); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	return comparison, nil
}

// RemoveQuote removes a quote from the comparison
func (s *ComparisonService) RemoveQuote(ctx context.Context, tenantID, id, quoteID string) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := comparison.RemoveQuote(quoteID); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	return comparison, nil
}

// CalculateScores calculates scores for all quotes in the comparison
func (s *ComparisonService) CalculateScores(ctx context.Context, tenantID, id string, quotes []Quote) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("no quotes provided for scoring")
	}

	// Calculate individual scores
	scores := s.calculateQuoteScores(quotes, comparison.ScoringCriteria)

	// Rank scores
	scores = s.rankScores(scores)

	// Generate recommendations
	for i := range scores {
		scores[i].Recommendation = s.generateRecommendation(&scores[i], len(quotes))
	}

	// Calculate price differences
	priceDiffs := s.calculatePriceDifferences(quotes)

	// Build item comparisons
	itemComps := s.buildItemComparisons(quotes)

	// Generate overall recommendation
	overallRec := s.generateOverallRecommendation(scores, comparison.ScoringCriteria)

	// Update comparison
	comparison.SetScores(scores)
	comparison.SetPriceDifferences(priceDiffs)
	comparison.SetItemComparisons(itemComps)
	comparison.SetRecommendation(overallRec)

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	s.logger.Info("Scores calculated", slog.String("comparison_id", id), slog.Int("quote_count", len(quotes)))
	return comparison, nil
}

// calculateQuoteScores calculates scores for each quote
func (s *ComparisonService) calculateQuoteScores(quotes []Quote, criteria domain.ScoringCriteria) []domain.QuoteScore {
	scores := make([]domain.QuoteScore, len(quotes))
	now := time.Now()

	for i, quote := range quotes {
		priceScore := s.calculatePriceScore(quote, quotes)
		qualityScore := s.calculateQualityScore(quote)
		deliveryScore := s.calculateDeliveryScore(quote)
		complianceScore := s.calculateComplianceScore(quote)

		overallScore := (priceScore * criteria.PriceWeight / 100.0) +
			(qualityScore * criteria.QualityWeight / 100.0) +
			(deliveryScore * criteria.DeliveryWeight / 100.0) +
			(complianceScore * criteria.ComplianceWeight / 100.0)

		strengths, weaknesses := s.analyzeQuoteStrengthsWeaknesses(priceScore, qualityScore, deliveryScore, complianceScore)

		scores[i] = domain.QuoteScore{
			QuoteID:         quote.ID,
			QuoteNumber:     quote.QuoteNumber,
			SupplierID:      quote.SupplierID,
			SupplierName:    quote.SupplierName,
			TotalAmount:     quote.TotalAmount,
			PriceScore:      priceScore,
			QualityScore:    qualityScore,
			DeliveryScore:   deliveryScore,
			ComplianceScore: complianceScore,
			OverallScore:    overallScore,
			Strengths:       strengths,
			Weaknesses:      weaknesses,
			CalculatedAt:    now,
		}
	}

	return scores
}

// calculatePriceScore calculates price competitiveness (0-100, higher is better)
func (s *ComparisonService) calculatePriceScore(quote Quote, allQuotes []Quote) float64 {
	if len(allQuotes) == 0 {
		return 50.0 // Default if no comparison
	}

	// Find min and max prices
	minPrice := math.MaxFloat64
	maxPrice := 0.0
	for _, q := range allQuotes {
		if q.TotalAmount < minPrice {
			minPrice = q.TotalAmount
		}
		if q.TotalAmount > maxPrice {
			maxPrice = q.TotalAmount
		}
	}

	// Avoid division by zero
	if maxPrice == minPrice {
		return 100.0
	}

	// Lower price gets higher score (inverted scale)
	// Best price gets 100, worst gets 0
	normalized := (maxPrice - quote.TotalAmount) / (maxPrice - minPrice)
	return normalized * 100.0
}

// calculateQualityScore calculates quality indicators (0-100)
func (s *ComparisonService) calculateQualityScore(quote Quote) float64 {
	score := 0.0
	factors := 0.0

	// Warranty terms (40 points)
	if quote.WarrantyTerms != "" {
		warrantyScore := s.scoreWarranty(quote.WarrantyTerms)
		score += warrantyScore * 0.4
	}
	factors += 0.4

	// Manufacturer reputation (30 points) - based on known manufacturers
	for _, item := range quote.Items {
		if s.isReputableManufacturer(item.ManufacturerName) {
			score += 30.0
			break
		}
	}
	factors += 0.3

	// Specifications completeness (30 points)
	specsCompleteness := s.scoreSpecificationsCompleteness(quote.Items)
	score += specsCompleteness * 0.3
	factors += 0.3

	return (score / factors) * 100.0
}

// calculateDeliveryScore calculates delivery speed score (0-100)
func (s *ComparisonService) calculateDeliveryScore(quote Quote) float64 {
	if len(quote.Items) == 0 {
		return 50.0 // Default
	}

	// Average delivery timeframe across items
	totalDays := 0.0
	validItems := 0

	for _, item := range quote.Items {
		days := s.parseDeliveryDays(item.DeliveryTimeframe)
		if days > 0 {
			totalDays += days
			validItems++
		}
	}

	if validItems == 0 {
		return 50.0 // Default if no delivery info
	}

	avgDays := totalDays / float64(validItems)

	// Score based on delivery speed
	// 0-30 days: 90-100
	// 31-60 days: 70-89
	// 61-90 days: 50-69
	// 91+ days: 0-49
	switch {
	case avgDays <= 30:
		return 90.0 + (30.0-avgDays)/30.0*10.0
	case avgDays <= 60:
		return 70.0 + (60.0-avgDays)/30.0*20.0
	case avgDays <= 90:
		return 50.0 + (90.0-avgDays)/30.0*20.0
	default:
		return math.Max(0, 50.0-(avgDays-90.0)/90.0*50.0)
	}
}

// calculateComplianceScore calculates compliance and certification score (0-100)
func (s *ComparisonService) calculateComplianceScore(quote Quote) float64 {
	if len(quote.Items) == 0 {
		return 0.0
	}

	totalScore := 0.0
	for _, item := range quote.Items {
		itemScore := s.scoreCompliance(item.ComplianceCerts)
		totalScore += itemScore
	}

	return (totalScore / float64(len(quote.Items))) * 100.0
}

// Helper scoring functions

func (s *ComparisonService) scoreWarranty(warrantyTerms string) float64 {
	lower := strings.ToLower(warrantyTerms)
	if strings.Contains(lower, "5 year") || strings.Contains(lower, "5-year") {
		return 1.0
	} else if strings.Contains(lower, "3 year") || strings.Contains(lower, "3-year") {
		return 0.8
	} else if strings.Contains(lower, "2 year") || strings.Contains(lower, "2-year") {
		return 0.6
	} else if strings.Contains(lower, "1 year") || strings.Contains(lower, "1-year") {
		return 0.4
	}
	return 0.2
}

func (s *ComparisonService) isReputableManufacturer(name string) bool {
	reputable := []string{
		"siemens", "ge", "philips", "medtronic", "stryker",
		"boston scientific", "abbott", "johnson & johnson", "roche",
		"baxter", "becton dickinson", "cardinal health",
	}
	lower := strings.ToLower(name)
	for _, rep := range reputable {
		if strings.Contains(lower, rep) {
			return true
		}
	}
	return false
}

func (s *ComparisonService) scoreSpecificationsCompleteness(items []QuoteItem) float64 {
	if len(items) == 0 {
		return 0.0
	}

	completeCount := 0
	for _, item := range items {
		if item.Specifications != "" && item.ModelNumber != "" {
			completeCount++
		}
	}

	return float64(completeCount) / float64(len(items))
}

func (s *ComparisonService) parseDeliveryDays(timeframe string) float64 {
	lower := strings.ToLower(timeframe)

	// Try to extract number
	var days float64
	if strings.Contains(lower, "day") {
		fmt.Sscanf(lower, "%f", &days)
		return days
	}
	if strings.Contains(lower, "week") {
		fmt.Sscanf(lower, "%f", &days)
		return days * 7
	}
	if strings.Contains(lower, "month") {
		fmt.Sscanf(lower, "%f", &days)
		return days * 30
	}

	return 0.0
}

func (s *ComparisonService) scoreCompliance(certs string) float64 {
	if certs == "" {
		return 0.0
	}

	score := 0.0
	lower := strings.ToLower(certs)

	if strings.Contains(lower, "fda") {
		score += 0.4
	}
	if strings.Contains(lower, "ce") || strings.Contains(lower, "ce mark") {
		score += 0.3
	}
	if strings.Contains(lower, "iso") {
		score += 0.2
	}
	if strings.Contains(lower, "ul") || strings.Contains(lower, "csa") {
		score += 0.1
	}

	return math.Min(score, 1.0)
}

// analyzeQuoteStrengthsWeaknesses identifies strengths and weaknesses
func (s *ComparisonService) analyzeQuoteStrengthsWeaknesses(price, quality, delivery, compliance float64) ([]string, []string) {
	strengths := []string{}
	weaknesses := []string{}

	if price >= 80 {
		strengths = append(strengths, "Highly competitive pricing")
	} else if price < 50 {
		weaknesses = append(weaknesses, "Higher price compared to alternatives")
	}

	if quality >= 80 {
		strengths = append(strengths, "Excellent quality indicators")
	} else if quality < 50 {
		weaknesses = append(weaknesses, "Quality indicators below average")
	}

	if delivery >= 80 {
		strengths = append(strengths, "Fast delivery timeframe")
	} else if delivery < 50 {
		weaknesses = append(weaknesses, "Longer delivery timeframe")
	}

	if compliance >= 80 {
		strengths = append(strengths, "Strong compliance and certifications")
	} else if compliance < 50 {
		weaknesses = append(weaknesses, "Limited compliance certifications")
	}

	return strengths, weaknesses
}

// rankScores assigns ranks to scores (1 = best)
func (s *ComparisonService) rankScores(scores []domain.QuoteScore) []domain.QuoteScore {
	// Sort by overall score descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].OverallScore > scores[j].OverallScore
	})

	// Assign ranks
	for i := range scores {
		scores[i].Rank = i + 1
	}

	return scores
}

// generateRecommendation generates a recommendation for a quote
func (s *ComparisonService) generateRecommendation(score *domain.QuoteScore, totalQuotes int) string {
	switch score.Rank {
	case 1:
		return fmt.Sprintf("Best overall choice with a score of %.1f. Recommended for award.", score.OverallScore)
	case 2:
		return fmt.Sprintf("Strong alternative with a score of %.1f. Consider as backup option.", score.OverallScore)
	default:
		if score.Rank <= totalQuotes/2 {
			return fmt.Sprintf("Good option with a score of %.1f.", score.OverallScore)
		}
		return fmt.Sprintf("Lower-ranked option with a score of %.1f. May not be the best choice.", score.OverallScore)
	}
}

// calculatePriceDifferences calculates price differences from the lowest quote
func (s *ComparisonService) calculatePriceDifferences(quotes []Quote) []domain.PriceDifference {
	if len(quotes) == 0 {
		return []domain.PriceDifference{}
	}

	// Find minimum price
	minPrice := math.MaxFloat64
	for _, q := range quotes {
		if q.TotalAmount < minPrice {
			minPrice = q.TotalAmount
		}
	}

	diffs := make([]domain.PriceDifference, len(quotes))
	for i, q := range quotes {
		diff := q.TotalAmount - minPrice
		percentage := 0.0
		if minPrice > 0 {
			percentage = (diff / minPrice) * 100.0
		}

		diffs[i] = domain.PriceDifference{
			QuoteID:              q.ID,
			QuoteNumber:          q.QuoteNumber,
			TotalAmount:          q.TotalAmount,
			DifferenceFromLowest: diff,
			PercentageFromLowest: percentage,
		}
	}

	return diffs
}

// buildItemComparisons builds item-level comparisons
func (s *ComparisonService) buildItemComparisons(quotes []Quote) []domain.ItemComparison {
	// Group items by equipment
	itemMap := make(map[string]*domain.ItemComparison)

	for _, quote := range quotes {
		for _, item := range quote.Items {
			key := item.EquipmentID
			if key == "" {
				key = item.EquipmentName
			}

			if _, exists := itemMap[key]; !exists {
				itemMap[key] = &domain.ItemComparison{
					EquipmentID:   item.EquipmentID,
					EquipmentName: item.EquipmentName,
					Quotes:        make(map[string]domain.ItemDetails),
				}
			}

			itemMap[key].Quotes[quote.ID] = domain.ItemDetails{
				Quantity:          item.Quantity,
				UnitPrice:         item.UnitPrice,
				TotalPrice:        item.TotalPrice,
				DeliveryTimeframe: item.DeliveryTimeframe,
				ManufacturerName:  item.ManufacturerName,
				ModelNumber:       item.ModelNumber,
				Specifications:    item.Specifications,
				ComplianceCerts:   item.ComplianceCerts,
			}
		}
	}

	// Convert map to slice
	comparisons := make([]domain.ItemComparison, 0, len(itemMap))
	for _, comp := range itemMap {
		comparisons = append(comparisons, *comp)
	}

	return comparisons
}

// generateOverallRecommendation generates an overall recommendation
func (s *ComparisonService) generateOverallRecommendation(scores []domain.QuoteScore, criteria domain.ScoringCriteria) string {
	if len(scores) == 0 {
		return "No quotes to compare."
	}

	best := scores[0]
	recommendation := fmt.Sprintf("Based on the analysis of %d quotes, Quote %s from %s is recommended with an overall score of %.1f/100. ",
		len(scores), best.QuoteNumber, best.SupplierName, best.OverallScore)

	// Add key factors
	factors := []string{}
	if criteria.PriceWeight >= 40 {
		factors = append(factors, fmt.Sprintf("price competitiveness (%.1f/100)", best.PriceScore))
	}
	if criteria.QualityWeight >= 30 {
		factors = append(factors, fmt.Sprintf("quality indicators (%.1f/100)", best.QualityScore))
	}
	if criteria.DeliveryWeight >= 20 {
		factors = append(factors, fmt.Sprintf("delivery timeframe (%.1f/100)", best.DeliveryScore))
	}

	if len(factors) > 0 {
		recommendation += "Key factors: " + strings.Join(factors, ", ") + "."
	}

	return recommendation
}

// ActivateComparison activates a comparison
func (s *ComparisonService) ActivateComparison(ctx context.Context, tenantID, id string) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := comparison.Activate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	return comparison, nil
}

// CompleteComparison completes a comparison
func (s *ComparisonService) CompleteComparison(ctx context.Context, tenantID, id string) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := comparison.Complete(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	return comparison, nil
}

// ArchiveComparison archives a comparison
func (s *ComparisonService) ArchiveComparison(ctx context.Context, tenantID, id string) (*domain.Comparison, error) {
	comparison, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := comparison.Archive(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, comparison); err != nil {
		return nil, err
	}

	return comparison, nil
}

// DeleteComparison deletes a comparison
func (s *ComparisonService) DeleteComparison(ctx context.Context, tenantID, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
