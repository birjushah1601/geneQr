package app

import (
    "context"
    "encoding/json"
    "log/slog"
    "testing"
    "time"

    equipmentDomain "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
    ticketDomain "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
)

// --- fakes ---
type fakeTicketRepo struct{ m map[string]*ticketDomain.ServiceTicket }
func (f *fakeTicketRepo) Create(ctx context.Context, t *ticketDomain.ServiceTicket) error { if f.m==nil{f.m=map[string]*ticketDomain.ServiceTicket{}}; f.m[t.ID]=t; return nil }
func (f *fakeTicketRepo) Update(ctx context.Context, t *ticketDomain.ServiceTicket) error { f.m[t.ID]=t; return nil }
func (f *fakeTicketRepo) UpdateResponsibility(ctx context.Context, id string, orgID *string, prov json.RawMessage) error { return nil }
func (f *fakeTicketRepo) GetByID(ctx context.Context, id string) (*ticketDomain.ServiceTicket, error) { return f.m[id], nil }
func (f *fakeTicketRepo) GetByTicketNumber(ctx context.Context, n string) (*ticketDomain.ServiceTicket, error) { return nil, nil }
func (f *fakeTicketRepo) List(ctx context.Context, c ticketDomain.ListCriteria) (*ticketDomain.TicketListResult, error) { return nil, nil }
func (f *fakeTicketRepo) GetByEquipment(ctx context.Context, equipmentID string) ([]*ticketDomain.ServiceTicket, error) { return nil, nil }
func (f *fakeTicketRepo) GetByCustomer(ctx context.Context, customerID string) ([]*ticketDomain.ServiceTicket, error) { return nil, nil }
func (f *fakeTicketRepo) GetByEngineer(ctx context.Context, engineerID string) ([]*ticketDomain.ServiceTicket, error) { return nil, nil }
func (f *fakeTicketRepo) GetBySource(ctx context.Context, source ticketDomain.TicketSource) ([]*ticketDomain.ServiceTicket, error) { return nil, nil }
func (f *fakeTicketRepo) AddComment(ctx context.Context, c *ticketDomain.TicketComment) error { return nil }
func (f *fakeTicketRepo) GetComments(ctx context.Context, id string) ([]*ticketDomain.TicketComment, error) { return nil, nil }
func (f *fakeTicketRepo) AddStatusHistory(ctx context.Context, h *ticketDomain.StatusHistory) error { return nil }
func (f *fakeTicketRepo) GetStatusHistory(ctx context.Context, id string) ([]*ticketDomain.StatusHistory, error) { return nil, nil }

type fakeEquipRepo struct{}
func (f *fakeEquipRepo) Create(ctx context.Context, e *equipmentDomain.Equipment) error { return nil }
func (f *fakeEquipRepo) GetByID(ctx context.Context, id string) (*equipmentDomain.Equipment, error) { return nil, nil }
func (f *fakeEquipRepo) GetByQRCode(ctx context.Context, code string) (*equipmentDomain.Equipment, error) { return nil, nil }
func (f *fakeEquipRepo) GetBySerialNumber(ctx context.Context, sn string) (*equipmentDomain.Equipment, error) { return nil, nil }
func (f *fakeEquipRepo) List(ctx context.Context, c equipmentDomain.ListCriteria) (*equipmentDomain.ListResult, error) { return nil, nil }
func (f *fakeEquipRepo) Update(ctx context.Context, e *equipmentDomain.Equipment) error { return nil }
func (f *fakeEquipRepo) Delete(ctx context.Context, id string) error { return nil }
func (f *fakeEquipRepo) BulkCreate(ctx context.Context, equipment []*equipmentDomain.Equipment) error { return nil }
func (f *fakeEquipRepo) UpdateQRCode(ctx context.Context, equipmentID string, qrImage []byte, format string) error { return nil }

type fakePolicyRepo struct{ rules *ticketDomain.SLARules; respOrg *string }
func (f *fakePolicyRepo) GetDefaultResponsibleOrg(ctx context.Context) (*string, error) { return f.respOrg, nil }
func (f *fakePolicyRepo) GetSLARules(ctx context.Context, orgID *string) (*ticketDomain.SLARules, error) { return f.rules, nil }

type fakeEventRepo struct{ created bool; enqueued bool }
func (f *fakeEventRepo) CreateEvent(ctx context.Context, eventType, aggregateType, aggregateID string, payload json.RawMessage) (string, error) { f.created=true; return "evt1", nil }
func (f *fakeEventRepo) EnqueueDeliveriesForEvent(ctx context.Context, eventID string, eventType string) error { f.enqueued=true; return nil }

// --- tests ---
func TestCreateTicket_SLAFromPolicy(t *testing.T) {
    repo := &fakeTicketRepo{}
    equip := &fakeEquipRepo{}
    rules := &ticketDomain.SLARules{}
    rules.High.Response, rules.High.Resolution = 2, 8
    s := NewTicketService(repo, equip, &fakePolicyRepo{rules: rules}, &fakeEventRepo{}, testLogger())

    now := time.Now()
    req := CreateTicketRequest{
        EquipmentID: "eq1", SerialNumber: "SN", EquipmentName: "EQ", CustomerName: "C",
        IssueDescription: "desc", Priority: ticketDomain.PriorityHigh, Source: ticketDomain.SourceWeb, CreatedBy: "u",
    }
    ticket, err := s.CreateTicket(context.Background(), req)
    if err != nil { t.Fatalf("CreateTicket error: %v", err) }
    if ticket.SLAResponseDue.Before(now.Add(2*time.Hour)) { t.Fatalf("response due not applied from policy") }
    if ticket.SLAResolutionDue.Before(now.Add(8*time.Hour)) { t.Fatalf("resolution due not applied from policy") }
}

func TestAssignTicket_EmitsEvent(t *testing.T) {
    repo := &fakeTicketRepo{}
    equip := &fakeEquipRepo{}
    ev := &fakeEventRepo{}
    s := NewTicketService(repo, equip, &fakePolicyRepo{}, ev, testLogger())
    // create a ticket first
    ticket, _ := s.CreateTicket(context.Background(), CreateTicketRequest{
        EquipmentID: "eq1", SerialNumber: "SN", EquipmentName: "EQ", CustomerName: "C",
        IssueDescription: "desc", Priority: ticketDomain.PriorityMedium, Source: ticketDomain.SourceWeb, CreatedBy: "u",
    })
    if err := s.AssignTicket(context.Background(), ticket.ID, "eng1", "Eng One", "u"); err != nil {
        t.Fatalf("AssignTicket error: %v", err)
    }
    if !ev.created || !ev.enqueued { t.Fatalf("expected event created and enqueued") }
}

// testLogger returns a no-op slog.Logger
func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelError})) }
