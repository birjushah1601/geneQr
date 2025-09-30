package infra

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB wraps the PostgreSQL connection pool
type PostgresDB struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(ctx context.Context, connString string, logger *slog.Logger) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to PostgreSQL database",
		slog.String("component", "postgres_db"))

	return &PostgresDB{
		pool:   pool,
		logger: logger,
	}, nil
}

// Close closes the database connection pool
func (db *PostgresDB) Close() {
	db.pool.Close()
	db.logger.Info("PostgreSQL connection pool closed")
}

// Pool returns the underlying connection pool
func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}
