package infra

import (
    "context"
    "encoding/json"

    domain "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
    "github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct { pool *pgxpool.Pool }

func NewEventRepository(pool *pgxpool.Pool) *EventRepository { return &EventRepository{pool: pool} }

func (r *EventRepository) CreateEvent(ctx context.Context, eventType, aggregateType, aggregateID string, payload json.RawMessage) (string, error) {
    const q = `INSERT INTO service_events(event_type, aggregate_type, aggregate_id, payload)
               VALUES ($1,$2,$3,$4)
               RETURNING id`
    var id string
    if err := r.pool.QueryRow(ctx, q, eventType, aggregateType, aggregateID, payload).Scan(&id); err != nil { return "", err }
    return id, nil
}

func (r *EventRepository) EnqueueDeliveriesForEvent(ctx context.Context, eventID string, eventType string) error {
    // insert deliveries for all active subscriptions matching type or '*'
    const q = `INSERT INTO webhook_deliveries(event_id, subscription_id)
               SELECT $1, ws.id
               FROM webhook_subscriptions ws
               WHERE ws.active = true AND (ws.event_types @> ARRAY[$2]::text[] OR ws.event_types @> ARRAY['*']::text[])`
    _, err := r.pool.Exec(ctx, q, eventID, eventType)
    return err
}

var _ domain.EventRepository = (*EventRepository)(nil)
