# Webhooks and SLA Standards

This document defines the standards for outbound webhooks (event delivery) and the SLA policy DSL used by the service ticketing domain.

## Webhooks

- Transport: HTTPS POST only
- Content-Type: application/json
- Idempotency: Receivers must be idempotent by `event_id`
- Security: HMAC SHA-256 signature with a per-subscription secret
- Replay protection: 5-minute window based on timestamp

Headers
- `X-Webhook-Event`: event type (e.g., `ticket.created`)
- `X-Webhook-Timestamp`: Unix seconds when signed
- `X-Webhook-Signature`: `t=<timestamp>,v1=<hex(hmac_sha256(secret, t "." body))>`

Verification (receiver side)
1. Parse `t` and ensure |now - t| <= 300s
2. Compute `expected = HMAC_SHA256(secret, t + "." + raw_body)`
3. Compare `hex(expected)` to `v1` (constant-time compare)

Retries and Backoff
- Delivery statuses: `queued | delivered | failed`
- Exponential backoff on non-2xx responses (e.g., 1m, 5m, 15m, 60m)
- Max attempts default: 10 (configurable)
- Dead-letter strategy: mark as `failed` after max attempts; operators can requeue

Event Payload (envelope)
```
{
  "id": "<event_uuid>",
  "type": "ticket.created",
  "occurred_at": "<rfc3339>",
  "data": { ... domain-specific ... }
}
```

Event Types (initial)
- `ticket.created`, `ticket.assigned`, `ticket.acknowledged`, `ticket.started`, `ticket.on_hold`, `ticket.resumed`, `ticket.resolved`, `ticket.closed`, `ticket.cancelled`, `ticket.commented`

Subscription Model
- `webhook_subscriptions(name, endpoint_url, event_types, secret, active)`
- `event_types` supports exact match or `*`

## SLA Policy DSL

Location: `sla_policies.rules` JSON

Two supported shapes (both valid):
1) Flat
```
{
  "critical": { "resp": 1,  "res": 4  },
  "high":     { "resp": 2,  "res": 8  },
  "medium":   { "resp": 4,  "res": 24 },
  "low":      { "resp": 8,  "res": 48 }
}
```

2) Nested under `priority`
```
{
  "priority": {
    "critical": { "resp": 1,  "res": 4  },
    "high":     { "resp": 2,  "res": 8  },
    "medium":   { "resp": 4,  "res": 24 },
    "low":      { "resp": 8,  "res": 48 }
  }
}
```

Resolution Algorithm (CreateTicket)
1. Load active org-scoped policy (if any), else global
2. Select pair for `ticket.priority`
3. If missing, fallback to defaults

Future Extensions
- Business hours calendars and holidays per org
- Severity-based multipliers and customer tier adjustments
- Pausing SLA during `on_hold` with reason codes
- Region-based overrides and daytime windows
