package reports

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// DailyReportService generates daily reports for admins
type DailyReportService struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewDailyReportService creates a new daily report service
func NewDailyReportService(db *sql.DB, logger *slog.Logger) *DailyReportService {
	return &DailyReportService{
		db:     db,
		logger: logger,
	}
}

// DailyReportData contains all data for the daily report
type DailyReportData struct {
	ReportDate      time.Time
	ReportType      string // "morning" or "evening"
	
	// Ticket Statistics
	TotalTickets           int
	NewTicketsToday        int
	ResolvedTicketsToday   int
	PendingTickets         int
	InProgressTickets      int
	OnHoldTickets          int
	
	// Priority Breakdown
	CriticalTickets        int
	HighPriorityTickets    int
	MediumPriorityTickets  int
	LowPriorityTickets     int
	
	// Engineer Statistics
	TotalEngineers         int
	ActiveEngineers        int
	EngineersWithTickets   int
	AverageTicketsPerEngineer float64
	
	// Equipment Statistics
	TotalEquipment         int
	EquipmentWithIssues    int
	EquipmentServiced      int
	
	// Performance Metrics
	AverageResolutionTime  float64 // in hours
	TicketsSLA             int     // tickets within SLA
	TicketsOverdue         int     // tickets past SLA
	
	// Top Lists
	TopIssueTypes          []IssueTypeStat
	TopEngineers           []EngineerStat
	TopEquipment           []EquipmentStat
	
	// Recent Activity
	RecentTickets          []RecentTicket
	TicketsNeedingAttention []TicketAlert
}

// IssueTypeStat represents issue type statistics
type IssueTypeStat struct {
	IssueType string
	Count     int
}

// EngineerStat represents engineer performance statistics
type EngineerStat struct {
	EngineerName      string
	TicketsAssigned   int
	TicketsResolved   int
	AverageResolution float64 // in hours
}

// EquipmentStat represents equipment with most issues
type EquipmentStat struct {
	EquipmentName  string
	Manufacturer   string
	IssueCount     int
	LastServiceDate time.Time
}

// RecentTicket represents a recently created ticket
type RecentTicket struct {
	TicketNumber  string
	CustomerName  string
	EquipmentName string
	Priority      string
	Status        string
	CreatedAt     time.Time
}

// TicketAlert represents tickets needing attention
type TicketAlert struct {
	TicketNumber  string
	CustomerName  string
	EquipmentName string
	Priority      string
	Status        string
	DaysOpen      int
	Reason        string
}

// GenerateDailyReport generates a comprehensive daily report
func (s *DailyReportService) GenerateDailyReport(ctx context.Context, reportType string) (*DailyReportData, error) {
	s.logger.Info("Generating daily report", slog.String("type", reportType))
	
	report := &DailyReportData{
		ReportDate: time.Now(),
		ReportType: reportType,
	}
	
	// Get ticket statistics
	if err := s.getTicketStatistics(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to get ticket statistics: %w", err)
	}
	
	// Get engineer statistics
	if err := s.getEngineerStatistics(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to get engineer statistics: %w", err)
	}
	
	// Get equipment statistics
	if err := s.getEquipmentStatistics(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to get equipment statistics: %w", err)
	}
	
	// Get performance metrics
	if err := s.getPerformanceMetrics(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to get performance metrics: %w", err)
	}
	
	// Get top lists
	if err := s.getTopLists(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to get top lists: %w", err)
	}
	
	// Get recent activity
	if err := s.getRecentActivity(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	
	s.logger.Info("Daily report generated successfully",
		slog.String("type", reportType),
		slog.Int("total_tickets", report.TotalTickets),
		slog.Int("new_today", report.NewTicketsToday),
	)
	
	return report, nil
}

// getTicketStatistics retrieves ticket-related statistics
func (s *DailyReportService) getTicketStatistics(ctx context.Context, report *DailyReportData) error {
	// Total tickets
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets
	`).Scan(&report.TotalTickets)
	if err != nil {
		return err
	}
	
	// New tickets today
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.NewTicketsToday)
	if err != nil {
		return err
	}
	
	// Resolved tickets today
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets
		WHERE status = 'resolved'
		AND DATE(updated_at) = CURRENT_DATE
	`).Scan(&report.ResolvedTicketsToday)
	if err != nil {
		return err
	}
	
	// Status breakdown
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets WHERE status = 'new' OR status = 'assigned'
	`).Scan(&report.PendingTickets)
	if err != nil {
		return err
	}
	
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets WHERE status = 'in_progress'
	`).Scan(&report.InProgressTickets)
	if err != nil {
		return err
	}
	
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets WHERE status = 'on_hold'
	`).Scan(&report.OnHoldTickets)
	if err != nil {
		return err
	}
	
	// Priority breakdown
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets WHERE priority = 'critical' AND status NOT IN ('resolved', 'closed')
	`).Scan(&report.CriticalTickets)
	if err != nil {
		return err
	}
	
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets WHERE priority = 'high' AND status NOT IN ('resolved', 'closed')
	`).Scan(&report.HighPriorityTickets)
	if err != nil {
		return err
	}
	
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets WHERE priority = 'medium' AND status NOT IN ('resolved', 'closed')
	`).Scan(&report.MediumPriorityTickets)
	if err != nil {
		return err
	}
	
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM service_tickets WHERE priority = 'low' AND status NOT IN ('resolved', 'closed')
	`).Scan(&report.LowPriorityTickets)
	if err != nil {
		return err
	}
	
	return nil
}

// getEngineerStatistics retrieves engineer-related statistics
func (s *DailyReportService) getEngineerStatistics(ctx context.Context, report *DailyReportData) error {
	// Total engineers
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM engineers
	`).Scan(&report.TotalEngineers)
	if err != nil {
		return err
	}
	
	// Engineers with active assignments
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT engineer_id) 
		FROM ticket_assignments 
		WHERE status = 'active'
	`).Scan(&report.EngineersWithTickets)
	if err != nil {
		return err
	}
	
	// Average tickets per engineer
	err = s.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(ticket_count), 0)
		FROM (
			SELECT engineer_id, COUNT(*) as ticket_count
			FROM ticket_assignments
			WHERE status = 'active'
			GROUP BY engineer_id
		) as counts
	`).Scan(&report.AverageTicketsPerEngineer)
	if err != nil {
		return err
	}
	
	report.ActiveEngineers = report.EngineersWithTickets
	
	return nil
}

// getEquipmentStatistics retrieves equipment-related statistics
func (s *DailyReportService) getEquipmentStatistics(ctx context.Context, report *DailyReportData) error {
	// Total equipment
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM equipment
	`).Scan(&report.TotalEquipment)
	if err != nil {
		return err
	}
	
	// Equipment with open issues
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT equipment_id)
		FROM service_tickets
		WHERE status NOT IN ('resolved', 'closed')
	`).Scan(&report.EquipmentWithIssues)
	if err != nil {
		return err
	}
	
	// Equipment serviced today
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT equipment_id)
		FROM service_tickets
		WHERE status = 'resolved'
		AND DATE(updated_at) = CURRENT_DATE
	`).Scan(&report.EquipmentServiced)
	if err != nil {
		return err
	}
	
	return nil
}

// getPerformanceMetrics retrieves performance-related metrics
func (s *DailyReportService) getPerformanceMetrics(ctx context.Context, report *DailyReportData) error {
	// Average resolution time (in hours)
	err := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600), 0)
		FROM service_tickets
		WHERE status = 'resolved'
		AND updated_at >= CURRENT_DATE - INTERVAL '7 days'
	`).Scan(&report.AverageResolutionTime)
	if err != nil {
		return err
	}
	
	// Tickets within SLA (resolved within 48 hours)
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM service_tickets
		WHERE status = 'resolved'
		AND (updated_at - created_at) <= INTERVAL '48 hours'
		AND DATE(updated_at) = CURRENT_DATE
	`).Scan(&report.TicketsSLA)
	if err != nil {
		return err
	}
	
	// Overdue tickets (open for more than 48 hours)
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM service_tickets
		WHERE status NOT IN ('resolved', 'closed')
		AND (CURRENT_TIMESTAMP - created_at) > INTERVAL '48 hours'
	`).Scan(&report.TicketsOverdue)
	if err != nil {
		return err
	}
	
	return nil
}

// getTopLists retrieves top performers and problem areas
func (s *DailyReportService) getTopLists(ctx context.Context, report *DailyReportData) error {
	// Top 5 engineers by resolved tickets
	rows, err := s.db.QueryContext(ctx, `
		SELECT 
			e.name,
			COUNT(CASE WHEN st.status = 'resolved' AND DATE(st.updated_at) >= CURRENT_DATE - INTERVAL '7 days' THEN 1 END) as resolved,
			COUNT(*) as assigned,
			COALESCE(AVG(CASE WHEN st.status = 'resolved' THEN EXTRACT(EPOCH FROM (st.updated_at - st.created_at))/3600 END), 0) as avg_resolution
		FROM engineers e
		LEFT JOIN ticket_assignments ta ON e.id = ta.engineer_id
		LEFT JOIN service_tickets st ON ta.ticket_id = st.id
		WHERE ta.status = 'active' OR ta.status = 'completed'
		GROUP BY e.id, e.name
		ORDER BY resolved DESC
		LIMIT 5
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	report.TopEngineers = []EngineerStat{}
	for rows.Next() {
		var stat EngineerStat
		if err := rows.Scan(&stat.EngineerName, &stat.TicketsResolved, &stat.TicketsAssigned, &stat.AverageResolution); err != nil {
			return err
		}
		report.TopEngineers = append(report.TopEngineers, stat)
	}
	
	// Top 5 equipment with most issues
	rows, err = s.db.QueryContext(ctx, `
		SELECT 
			eq.name,
			COALESCE(eq.manufacturer, 'Unknown') as manufacturer,
			COUNT(st.id) as issue_count,
			COALESCE(MAX(st.created_at), eq.created_at) as last_service
		FROM equipment_registry eq
		LEFT JOIN service_tickets st ON eq.id = st.equipment_id
		WHERE st.created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY eq.id, eq.name, eq.manufacturer, eq.created_at
		ORDER BY issue_count DESC
		LIMIT 5
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	report.TopEquipment = []EquipmentStat{}
	for rows.Next() {
		var stat EquipmentStat
		if err := rows.Scan(&stat.EquipmentName, &stat.Manufacturer, &stat.IssueCount, &stat.LastServiceDate); err != nil {
			return err
		}
		report.TopEquipment = append(report.TopEquipment, stat)
	}
	
	return nil
}

// getRecentActivity retrieves recent tickets and alerts
func (s *DailyReportService) getRecentActivity(ctx context.Context, report *DailyReportData) error {
	// Recent tickets (last 10)
	rows, err := s.db.QueryContext(ctx, `
		SELECT 
			ticket_number,
			customer_name,
			equipment_name,
			priority,
			status,
			created_at
		FROM service_tickets
		WHERE DATE(created_at) = CURRENT_DATE
		ORDER BY created_at DESC
		LIMIT 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	report.RecentTickets = []RecentTicket{}
	for rows.Next() {
		var ticket RecentTicket
		if err := rows.Scan(&ticket.TicketNumber, &ticket.CustomerName, &ticket.EquipmentName, &ticket.Priority, &ticket.Status, &ticket.CreatedAt); err != nil {
			return err
		}
		report.RecentTickets = append(report.RecentTickets, ticket)
	}
	
	// Tickets needing attention
	rows, err = s.db.QueryContext(ctx, `
		SELECT 
			ticket_number,
			customer_name,
			equipment_name,
			priority,
			status,
			EXTRACT(DAY FROM (CURRENT_TIMESTAMP - created_at))::INT as days_open,
			CASE 
				WHEN priority = 'critical' AND status = 'new' THEN 'Critical ticket unassigned'
				WHEN (CURRENT_TIMESTAMP - created_at) > INTERVAL '48 hours' AND status NOT IN ('resolved', 'closed') THEN 'Overdue (>48 hours)'
				WHEN status = 'on_hold' AND (CURRENT_TIMESTAMP - updated_at) > INTERVAL '24 hours' THEN 'On hold too long'
				ELSE 'Needs attention'
			END as reason
		FROM service_tickets
		WHERE 
			(priority = 'critical' AND status = 'new')
			OR ((CURRENT_TIMESTAMP - created_at) > INTERVAL '48 hours' AND status NOT IN ('resolved', 'closed'))
			OR (status = 'on_hold' AND (CURRENT_TIMESTAMP - updated_at) > INTERVAL '24 hours')
		ORDER BY 
			CASE priority 
				WHEN 'critical' THEN 1 
				WHEN 'high' THEN 2 
				WHEN 'medium' THEN 3 
				ELSE 4 
			END,
			created_at
		LIMIT 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	report.TicketsNeedingAttention = []TicketAlert{}
	for rows.Next() {
		var alert TicketAlert
		if err := rows.Scan(&alert.TicketNumber, &alert.CustomerName, &alert.EquipmentName, &alert.Priority, &alert.Status, &alert.DaysOpen, &alert.Reason); err != nil {
			return err
		}
		report.TicketsNeedingAttention = append(report.TicketsNeedingAttention, alert)
	}
	
	return nil
}
