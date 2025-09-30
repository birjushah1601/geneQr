package infra

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

// Pool returns the underlying connection pool
func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}
