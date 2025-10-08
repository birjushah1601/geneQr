package domain

import "context"

// PolicyRepository exposes access to service policy decisions
type PolicyRepository interface {
    // GetDefaultResponsibleOrg returns default_org_id from latest enabled policy, if configured
    GetDefaultResponsibleOrg(ctx context.Context) (*string, error)

    // GetSLARules returns SLA rules JSON parsed to a simple struct if available (org-scoped optional)
    GetSLARules(ctx context.Context, orgID *string) (*SLARules, error)
}

// SLARules holds response/resolution hours per priority
type SLARules struct {
    Critical struct{ Response int `json:"resp"`; Resolution int `json:"res"` } `json:"critical"`
    High     struct{ Response int `json:"resp"`; Resolution int `json:"res"` } `json:"high"`
    Medium   struct{ Response int `json:"resp"`; Resolution int `json:"res"` } `json:"medium"`
    Low      struct{ Response int `json:"resp"`; Resolution int `json:"res"` } `json:"low"`
}
