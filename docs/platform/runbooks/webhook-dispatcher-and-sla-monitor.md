# Runbook: Webhook Dispatcher and SLA Monitor

This runbook covers operations for the webhook dispatcher and SLA monitor workers.

## Components
- Webhook Dispatcher (ENABLE_EVENT_DISPATCHER)
- SLA Monitor (ENABLE_SLA_MONITOR)

## Configuration
- ENABLE_EVENT_DISPATCHER: enable/disable dispatcher
- ENABLE_SLA_MONITOR: enable/disable monitor
- WEBHOOK_MAX_ATTEMPTS: max delivery attempts (default 10)
- Delivery cadence (fixed): 1m, 5m, 15m, 60m, 180m (then capped)

## Dispatcher Behavior
- Polls queued deliveries; signs requests:
  - Headers: X-Webhook-Event, X-Webhook-Timestamp, X-Webhook-Signature
  - Signature: sha256 HMAC of `${timestamp}.${body}` with per-subscription secret
- Success: marks delivered
- Failure: increments attempt_count and requeues with backoff; marks failed (DLQ) after max attempts

## SLA Monitor Behavior
- Every minute:
  - Marks `sla_breached=true` when `acknowledged_at` is nil and now > `sla_response_due`
  - Marks `sla_breached=true` when `resolved_at` is nil and now > `sla_resolution_due`

## Operations
- Requeue failed deliveries: set `status='queued'`, reset `last_error` for selected rows
- Adjust sensitivity: change WEBHOOK_MAX_ATTEMPTS; cadence is code-defined
- Pause workers: unset flags or stop service

## Troubleshooting
- Check recent errors: `SELECT id,last_error,attempt_count FROM webhook_deliveries WHERE status!='delivered' ORDER BY last_attempt_at DESC LIMIT 50;`
- Verify subscriptions: endpoint URL reachable, secret correct, active=true, event_types matches
- Receiver validation: clock skew < 300s, matching signature

## Observability
- Add dashboards for: queue depth, delivery rate, failure rate, SLA breach count
- Alerting: delivery failures above threshold; sustained SLA breaches