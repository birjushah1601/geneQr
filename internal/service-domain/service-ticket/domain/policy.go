package domain

import "context"

// PolicyRepository exposes access to service policy decisions
type PolicyRepository interface {
    // GetDefaultResponsibleOrg returns default_org_id from latest enabled policy, if configured
    GetDefaultResponsibleOrg(ctx context.Context) (*string, error)
}
