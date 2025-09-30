package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/marketplace/catalog/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

// PostgresDB represents a PostgreSQL database connection pool
type PostgresDB struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewPostgresDB creates a new PostgreSQL database connection pool
func NewPostgresDB(ctx context.Context, dsn string, logger *slog.Logger) (*PostgresDB, error) {
	dbLogger := logger.With(slog.String("component", "postgres_db"))

	// Configure connection pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	// Set reasonable defaults for the connection pool
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbLogger.Info("Connected to PostgreSQL database")

	return &PostgresDB{
		pool:   pool,
		logger: dbLogger,
	}, nil
}

// Close closes the database connection pool
func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
		db.logger.Info("Database connection pool closed")
	}
}

// CatalogRepository implements the domain.CatalogRepository interface
type CatalogRepository struct {
	db     *PostgresDB
	logger *slog.Logger
}

// NewCatalogRepository creates a new catalog repository
func NewCatalogRepository(db *PostgresDB, logger *slog.Logger) *CatalogRepository {
	return &CatalogRepository{
		db:     db,
		logger: logger.With(slog.String("component", "catalog_repository")),
	}
}

// Create persists a new equipment to the database
func (r *CatalogRepository) Create(ctx context.Context, equipment *domain.Equipment) error {
	query := `
		INSERT INTO equipment (
			id, name, category_id, manufacturer_id, model, description, 
			specifications, price_amount, price_currency, sku, images, 
			is_active, created_at, updated_at, tenant_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	// Convert specifications to JSON
	specs, err := json.Marshal(equipment.Specifications)
	if err != nil {
		return fmt.Errorf("failed to convert specifications to JSON: %w", err)
	}

	// Execute the query
	_, err = r.db.pool.Exec(
		ctx,
		query,
		equipment.ID,
		equipment.Name,
		equipment.Category.ID,
		equipment.Manufacturer.ID,
		equipment.Model,
		equipment.Description,
		specs,
		equipment.Price.Amount,
		equipment.Price.Currency,
		equipment.SKU,
		equipment.Images,
		equipment.IsActive,
		equipment.CreatedAt,
		equipment.UpdatedAt,
		equipment.TenantID(),
	)

	if err != nil {
		r.logger.Error("Failed to create equipment",
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipment.ID))
		return fmt.Errorf("failed to create equipment: %w", err)
	}

	return nil
}

// GetByID retrieves equipment by ID
func (r *CatalogRepository) GetByID(ctx context.Context, id string, tenantID string) (*domain.Equipment, error) {
	query := `
		SELECT 
			e.id, e.name, e.model, e.description, e.specifications,
			e.price_amount, e.price_currency, e.sku, e.images, e.is_active,
			e.created_at, e.updated_at, e.tenant_id,
			c.id as category_id, c.name as category_name, c.parent_id as category_parent_id,
			m.id as manufacturer_id, m.name as manufacturer_name, 
			m.country as manufacturer_country, m.website as manufacturer_website
		FROM 
			equipment e
			JOIN categories c ON e.category_id = c.id
			JOIN manufacturers m ON e.manufacturer_id = m.id
		WHERE 
			e.id = $1 AND e.tenant_id = $2
	`

	row := r.db.pool.QueryRow(ctx, query, id, tenantID)

	var equipment domain.Equipment
	var categoryID, categoryName string
	var categoryParentID *string
	var manufacturerID, manufacturerName, manufacturerCountry, manufacturerWebsite string
	var specs []byte
	var images []string
	var equipmentTenantID string

	err := row.Scan(
		&equipment.ID,
		&equipment.Name,
		&equipment.Model,
		&equipment.Description,
		&specs,
		&equipment.Price.Amount,
		&equipment.Price.Currency,
		&equipment.SKU,
		&images,
		&equipment.IsActive,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
		&equipmentTenantID,
		&categoryID,
		&categoryName,
		&categoryParentID,
		&manufacturerID,
		&manufacturerName,
		&manufacturerCountry,
		&manufacturerWebsite,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEquipmentNotFound
		}
		r.logger.Error("Failed to get equipment by ID",
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		return nil, fmt.Errorf("failed to get equipment: %w", err)
	}

	// Set category and manufacturer
	equipment.Category = domain.Category{
		ID:       categoryID,
		Name:     categoryName,
		ParentID: categoryParentID,
	}

	equipment.Manufacturer = domain.Manufacturer{
		ID:      manufacturerID,
		Name:    manufacturerName,
		Country: manufacturerCountry,
		Website: manufacturerWebsite,
	}

	// Parse specifications JSON
	if err := json.Unmarshal(specs, &equipment.Specifications); err != nil {
		r.logger.Error("Failed to parse specifications JSON",
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		return nil, fmt.Errorf("failed to parse specifications: %w", err)
	}

	// assign tenant
	equipment.SetTenantID(equipmentTenantID)

	equipment.Images = images

	return &equipment, nil
}

// Update updates existing equipment
func (r *CatalogRepository) Update(ctx context.Context, equipment *domain.Equipment) error {
	query := `
		UPDATE equipment SET
			name = $1,
			category_id = $2,
			manufacturer_id = $3,
			model = $4,
			description = $5,
			specifications = $6,
			price_amount = $7,
			price_currency = $8,
			sku = $9,
			images = $10,
			is_active = $11,
			updated_at = $12
		WHERE
			id = $13 AND tenant_id = $14
	`

	// Convert specifications to JSON
	specs, err := json.Marshal(equipment.Specifications)
	if err != nil {
		return fmt.Errorf("failed to convert specifications to JSON: %w", err)
	}

	// Execute the query
	result, err := r.db.pool.Exec(
		ctx,
		query,
		equipment.Name,
		equipment.Category.ID,
		equipment.Manufacturer.ID,
		equipment.Model,
		equipment.Description,
		specs,
		equipment.Price.Amount,
		equipment.Price.Currency,
		equipment.SKU,
		equipment.Images,
		equipment.IsActive,
		equipment.UpdatedAt,
		equipment.ID,
		equipment.TenantID(),
	)

	if err != nil {
		r.logger.Error("Failed to update equipment",
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipment.ID))
		return fmt.Errorf("failed to update equipment: %w", err)
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		return domain.ErrEquipmentNotFound
	}

	return nil
}

// Delete removes equipment from the database
func (r *CatalogRepository) Delete(ctx context.Context, id string, tenantID string) error {
	query := `DELETE FROM equipment WHERE id = $1 AND tenant_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		r.logger.Error("Failed to delete equipment",
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		return fmt.Errorf("failed to delete equipment: %w", err)
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		return domain.ErrEquipmentNotFound
	}

	return nil
}

// Search searches for equipment based on criteria
func (r *CatalogRepository) Search(ctx context.Context, criteria domain.SearchCriteria) ([]*domain.Equipment, int, error) {
	// Build the query dynamically based on search criteria
	queryParts := []string{`
		SELECT 
			e.id, e.name, e.model, e.description, e.specifications,
			e.price_amount, e.price_currency, e.sku, e.images, e.is_active,
			e.created_at, e.updated_at, e.tenant_id,
			c.id as category_id, c.name as category_name, c.parent_id as category_parent_id,
			m.id as manufacturer_id, m.name as manufacturer_name, 
			m.country as manufacturer_country, m.website as manufacturer_website
		FROM 
			equipment e
			JOIN categories c ON e.category_id = c.id
			JOIN manufacturers m ON e.manufacturer_id = m.id
		WHERE 
			e.tenant_id = $1
	`}

	// Parameters for the query
	params := []interface{}{criteria.TenantID}
	paramIndex := 2 // Start from 2 since we already used $1

	// Add search conditions
	if criteria.Query != "" {
		queryParts = append(queryParts, fmt.Sprintf(`AND (
			e.name ILIKE $%d OR
			e.description ILIKE $%d OR
			e.model ILIKE $%d
		)`, paramIndex, paramIndex, paramIndex))
		params = append(params, "%"+criteria.Query+"%")
		paramIndex++
	}

	// Add category filter
	if len(criteria.CategoryIDs) > 0 {
		placeholders := make([]string, len(criteria.CategoryIDs))
		for i := range criteria.CategoryIDs {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, criteria.CategoryIDs[i])
			paramIndex++
		}
		queryParts = append(queryParts, fmt.Sprintf("AND e.category_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add manufacturer filter
	if len(criteria.ManufacturerIDs) > 0 {
		placeholders := make([]string, len(criteria.ManufacturerIDs))
		for i := range criteria.ManufacturerIDs {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, criteria.ManufacturerIDs[i])
			paramIndex++
		}
		queryParts = append(queryParts, fmt.Sprintf("AND e.manufacturer_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add price range filter
	if criteria.PriceMin != nil {
		queryParts = append(queryParts, fmt.Sprintf("AND e.price_amount >= $%d", paramIndex))
		params = append(params, *criteria.PriceMin)
		paramIndex++
	}

	if criteria.PriceMax != nil {
		queryParts = append(queryParts, fmt.Sprintf("AND e.price_amount <= $%d", paramIndex))
		params = append(params, *criteria.PriceMax)
		paramIndex++
	}

	// Add active filter
	if criteria.IsActive != nil {
		queryParts = append(queryParts, fmt.Sprintf("AND e.is_active = $%d", paramIndex))
		params = append(params, *criteria.IsActive)
		paramIndex++
	}

	// Add sorting
	sortColumn := "e.name"
	sortDirection := "ASC"

	if criteria.SortBy != "" {
		// Validate sort column to prevent SQL injection
		validSortColumns := map[string]string{
			"name":       "e.name",
			"price":      "e.price_amount",
			"created_at": "e.created_at",
		}

		if column, valid := validSortColumns[criteria.SortBy]; valid {
			sortColumn = column
		}
	}

	if criteria.SortDirection == "desc" {
		sortDirection = "DESC"
	}

	queryParts = append(queryParts, fmt.Sprintf("ORDER BY %s %s", sortColumn, sortDirection))

	// Add pagination
	offset := (criteria.Page - 1) * criteria.PageSize
	queryParts = append(queryParts, fmt.Sprintf("LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
	params = append(params, criteria.PageSize, offset)

	// Combine all parts into a single query
	query := strings.Join(queryParts, " ")

	// Execute the query
	rows, err := r.db.pool.Query(ctx, query, params...)
	if err != nil {
		r.logger.Error("Failed to search equipment",
			slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to search equipment: %w", err)
	}
	defer rows.Close()

	// Parse the results
	var equipment []*domain.Equipment
	for rows.Next() {
		var e domain.Equipment
		var categoryID, categoryName string
		var categoryParentID *string
		var manufacturerID, manufacturerName, manufacturerCountry, manufacturerWebsite string
		var specs []byte
		var images []string
		var equipmentTenantID string

		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Model,
			&e.Description,
			&specs,
			&e.Price.Amount,
			&e.Price.Currency,
			&e.SKU,
			&images,
			&e.IsActive,
			&e.CreatedAt,
			&e.UpdatedAt,
			&equipmentTenantID,
			&categoryID,
			&categoryName,
			&categoryParentID,
			&manufacturerID,
			&manufacturerName,
			&manufacturerCountry,
			&manufacturerWebsite,
		)

		if err != nil {
			r.logger.Error("Failed to scan equipment row",
				slog.String("error", err.Error()))
			return nil, 0, fmt.Errorf("failed to scan equipment row: %w", err)
		}

		// Set category and manufacturer
		e.Category = domain.Category{
			ID:       categoryID,
			Name:     categoryName,
			ParentID: categoryParentID,
		}

		e.Manufacturer = domain.Manufacturer{
			ID:      manufacturerID,
			Name:    manufacturerName,
			Country: manufacturerCountry,
			Website: manufacturerWebsite,
		}

		// Parse specifications JSON
		if err := json.Unmarshal(specs, &e.Specifications); err != nil {
			r.logger.Error("Failed to parse specifications JSON",
				slog.String("error", err.Error()),
				slog.String("equipment_id", e.ID))
			return nil, 0, fmt.Errorf("failed to parse specifications: %w", err)
		}

		e.Images = images
		e.SetTenantID(equipmentTenantID)
		equipment = append(equipment, &e)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over equipment rows",
			slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("error iterating over equipment rows: %w", err)
	}

	// Count total items for pagination
	// Reuse the WHERE conditions from the main query
	countQueryParts := strings.Split(query, "ORDER BY")[0]
	countQueryParts = strings.Replace(countQueryParts, `
		SELECT 
			e.id, e.name, e.model, e.description, e.specifications,
			e.price_amount, e.price_currency, e.sku, e.images, e.is_active,
			e.created_at, e.updated_at, e.tenant_id,
			c.id as category_id, c.name as category_name, c.parent_id as category_parent_id,
			m.id as manufacturer_id, m.name as manufacturer_name, 
			m.country as manufacturer_country, m.website as manufacturer_website
		FROM `, "SELECT COUNT(*) FROM ", 1)

	var total int
	err = r.db.pool.QueryRow(ctx, countQueryParts, params[:len(params)-2]...).Scan(&total)
	if err != nil {
		r.logger.Error("Failed to count total equipment",
			slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to count total equipment: %w", err)
	}

	return equipment, total, nil
}

// ListByCategory retrieves equipment by category
func (r *CatalogRepository) ListByCategory(ctx context.Context, categoryID string, tenantID string, page, pageSize int) ([]*domain.Equipment, int, error) {
	// Create search criteria with category filter
	criteria := domain.SearchCriteria{
		CategoryIDs: []string{categoryID},
		Page:        page,
		PageSize:    pageSize,
		TenantID:    tenantID,
	}

	// Reuse the search method
	return r.Search(ctx, criteria)
}

// ListByManufacturer retrieves equipment by manufacturer
func (r *CatalogRepository) ListByManufacturer(ctx context.Context, manufacturerID string, tenantID string, page, pageSize int) ([]*domain.Equipment, int, error) {
	// Create search criteria with manufacturer filter
	criteria := domain.SearchCriteria{
		ManufacturerIDs: []string{manufacturerID},
		Page:            page,
		PageSize:        pageSize,
		TenantID:        tenantID,
	}

	// Reuse the search method
	return r.Search(ctx, criteria)
}

// ListCategories retrieves all categories
func (r *CatalogRepository) ListCategories(ctx context.Context, tenantID string) ([]domain.Category, error) {
	query := `
		SELECT id, name, parent_id
		FROM categories
		WHERE tenant_id = $1
		ORDER BY name
	`

	rows, err := r.db.pool.Query(ctx, query, tenantID)
	if err != nil {
		r.logger.Error("Failed to list categories",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID); err != nil {
			r.logger.Error("Failed to scan category row",
				slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan category row: %w", err)
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over category rows",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("error iterating over category rows: %w", err)
	}

	return categories, nil
}

// ListManufacturers retrieves all manufacturers
func (r *CatalogRepository) ListManufacturers(ctx context.Context, tenantID string) ([]domain.Manufacturer, error) {
	query := `
		SELECT id, name, country, website
		FROM manufacturers
		WHERE tenant_id = $1
		ORDER BY name
	`

	rows, err := r.db.pool.Query(ctx, query, tenantID)
	if err != nil {
		r.logger.Error("Failed to list manufacturers",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to list manufacturers: %w", err)
	}
	defer rows.Close()

	var manufacturers []domain.Manufacturer
	for rows.Next() {
		var m domain.Manufacturer
		if err := rows.Scan(&m.ID, &m.Name, &m.Country, &m.Website); err != nil {
			r.logger.Error("Failed to scan manufacturer row",
				slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan manufacturer row: %w", err)
		}
		manufacturers = append(manufacturers, m)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over manufacturer rows",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("error iterating over manufacturer rows: %w", err)
	}

	return manufacturers, nil
}

// HealthCheck checks if the database is accessible
func (r *CatalogRepository) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.db.pool.Ping(ctx); err != nil {
		r.logger.Error("Database health check failed", slog.String("error", err.Error()))
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// KafkaEventPublisher implements the domain.EventPublisher interface
type KafkaEventPublisher struct {
	writer *kafka.Writer
	logger *slog.Logger
}

// NewKafkaEventPublisher creates a new Kafka event publisher
func NewKafkaEventPublisher(brokers []string, logger *slog.Logger) *KafkaEventPublisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        "medical-platform.catalog-events",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	return &KafkaEventPublisher{
		writer: writer,
		logger: logger.With(slog.String("component", "kafka_event_publisher")),
	}
}

// Publish publishes a domain event to Kafka
func (p *KafkaEventPublisher) Publish(ctx context.Context, event domain.DomainEvent) error {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to convert event to JSON",
			slog.String("error", err.Error()),
			slog.String("event_type", event.EventType()))
		return fmt.Errorf("failed to convert event to JSON: %w", err)
	}

	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(event.AggregateID()),
		Value: eventJSON,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.EventType())},
			{Key: "tenant_id", Value: []byte(event.TenantID())},
			{Key: "occurred_at", Value: []byte(event.OccurredAt().Format(time.RFC3339))},
		},
		Time: time.Now(),
	}

	// Publish the message
	if err := p.writer.WriteMessages(ctx, message); err != nil {
		p.logger.Error("Failed to publish event to Kafka",
			slog.String("error", err.Error()),
			slog.String("event_type", event.EventType()))
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("Event published successfully",
		slog.String("event_type", event.EventType()),
		slog.String("aggregate_id", event.AggregateID()))

	return nil
}

// Close closes the Kafka writer
func (p *KafkaEventPublisher) Close() error {
	if err := p.writer.Close(); err != nil {
		p.logger.Error("Failed to close Kafka writer", slog.String("error", err.Error()))
		return fmt.Errorf("failed to close Kafka writer: %w", err)
	}
	return nil
}
