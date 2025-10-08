package app

import (
    "bytes"
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type WebhookDispatcher struct {
    pool   *pgxpool.Pool
    client *http.Client
    logger *slog.Logger
}

func NewWebhookDispatcher(pool *pgxpool.Pool, logger *slog.Logger) *WebhookDispatcher {
    return &WebhookDispatcher{
        pool:   pool,
        client: &http.Client{Timeout: 10 * time.Second},
        logger: logger.With(slog.String("component", "webhook_dispatcher")),
    }
}

func (d *WebhookDispatcher) Run(ctx context.Context) {
    if !enabled(os.Getenv("ENABLE_EVENT_DISPATCHER")) {
        d.logger.Info("Dispatcher disabled; skipping run")
        return
    }
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            d.logger.Info("Dispatcher stopping")
            return
        case <-ticker.C:
            if err := d.dispatchBatch(ctx, 25); err != nil {
                d.logger.Error("dispatchBatch error", slog.String("error", err.Error()))
            }
        }
    }
}

type deliveryJob struct {
    DeliveryID   string
    EventID      string
    EventType    string
    Payload      []byte
    EndpointURL  string
    Secret       *string
}

func (d *WebhookDispatcher) dispatchBatch(ctx context.Context, limit int) error {
    const q = `SELECT d.id, e.id, e.event_type, e.payload, s.endpoint_url, s.secret
               FROM webhook_deliveries d
               JOIN service_events e ON e.id = d.event_id
               JOIN webhook_subscriptions s ON s.id = d.subscription_id
               WHERE d.status = 'queued'
               ORDER BY d.created_at
               LIMIT $1`
    rows, err := d.pool.Query(ctx, q, limit)
    if err != nil { return err }
    defer rows.Close()
    var jobs []deliveryJob
    for rows.Next() {
        var j deliveryJob
        if err := rows.Scan(&j.DeliveryID, &j.EventID, &j.EventType, &j.Payload, &j.EndpointURL, &j.Secret); err != nil {
            return err
        }
        jobs = append(jobs, j)
    }
    for _, j := range jobs {
        if err := d.deliver(ctx, j); err != nil {
            d.markFailure(ctx, j.DeliveryID, err)
        } else {
            d.markSuccess(ctx, j.DeliveryID)
        }
    }
    return nil
}

func (d *WebhookDispatcher) deliver(ctx context.Context, j deliveryJob) error {
    ts := fmt.Sprintf("%d", time.Now().Unix())
    body := j.Payload
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, j.EndpointURL, bytes.NewReader(body))
    if err != nil { return err }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Webhook-Event", j.EventType)
    req.Header.Set("X-Webhook-Timestamp", ts)
    if j.Secret != nil && *j.Secret != "" {
        sig := sign(*j.Secret, ts, body)
        req.Header.Set("X-Webhook-Signature", fmt.Sprintf("t=%s,v1=%s", ts, sig))
    }
    resp, err := d.client.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        return nil
    }
    return fmt.Errorf("non-2xx status: %d", resp.StatusCode)
}

func (d *WebhookDispatcher) markSuccess(ctx context.Context, deliveryID string) {
    const q = `UPDATE webhook_deliveries SET status='delivered', delivered_at=NOW(), last_attempt_at=NOW(), attempt_count=attempt_count+1 WHERE id=$1` 
    _, _ = d.pool.Exec(ctx, q, deliveryID)
}

func (d *WebhookDispatcher) markFailure(ctx context.Context, deliveryID string, err error) {
    const q = `UPDATE webhook_deliveries SET status='queued', last_error=$2, last_attempt_at=NOW(), attempt_count=attempt_count+1 WHERE id=$1`
    _, _ = d.pool.Exec(ctx, q, deliveryID, err.Error())
}

func sign(secret, ts string, body []byte) string {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(ts))
    mac.Write([]byte("."))
    mac.Write(body)
    sum := mac.Sum(nil)
    return hex.EncodeToString(sum)
}

// optional helper to pretty-print payload for debug
func pretty(b []byte) string {
    var m map[string]any
    if json.Unmarshal(b, &m) == nil {
        if x, _ := json.MarshalIndent(m, "", "  "); x != nil { return string(x) }
    }
    return string(b)
}
