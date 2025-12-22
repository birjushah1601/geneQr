package reports

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/infrastructure/config"
	"github.com/robfig/cron/v3"
)

// ReportScheduler handles scheduling of daily reports
type ReportScheduler struct {
	reportService *DailyReportService
	sendgridAPIKey string
	fromEmail      string
	fromName       string
	featureFlags  *config.FeatureFlags
	logger        *slog.Logger
	cron          *cron.Cron
	
	// Configuration
	morningTime   string   // e.g., "09:00"
	eveningTime   string   // e.g., "18:00"
	recipients    []string // Admin email addresses
	timezone      *time.Location
}

// NewReportScheduler creates a new report scheduler
func NewReportScheduler(
	reportService *DailyReportService,
	sendgridAPIKey, fromEmail, fromName string,
	featureFlags *config.FeatureFlags,
	logger *slog.Logger,
	morningTime, eveningTime string,
	recipients []string,
	timezone string,
) (*ReportScheduler, error) {
	// Parse timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Warn("Failed to load timezone, using UTC", "timezone", timezone, "error", err)
		loc = time.UTC
	}

	// Create cron scheduler
	cronScheduler := cron.New(cron.WithLocation(loc))

	return &ReportScheduler{
		reportService:  reportService,
		sendgridAPIKey: sendgridAPIKey,
		fromEmail:      fromEmail,
		fromName:       fromName,
		featureFlags:   featureFlags,
		logger:         logger,
		cron:           cronScheduler,
		morningTime:    morningTime,
		eveningTime:    eveningTime,
		recipients:     recipients,
		timezone:       loc,
	}, nil
}

// Start starts the report scheduler
func (s *ReportScheduler) Start() error {
	if !s.featureFlags.DailyReportsEnabled {
		s.logger.Info("Daily reports disabled by feature flag")
		return nil
	}

	if len(s.recipients) == 0 {
		s.logger.Warn("No recipients configured for daily reports")
		return nil
	}

	// Parse morning time
	morningHour, morningMin, err := parseTime(s.morningTime)
	if err != nil {
		return fmt.Errorf("invalid morning time: %w", err)
	}

	// Parse evening time
	eveningHour, eveningMin, err := parseTime(s.eveningTime)
	if err != nil {
		return fmt.Errorf("invalid evening time: %w", err)
	}

	// Schedule morning report
	morningCron := fmt.Sprintf("%d %d * * *", morningMin, morningHour)
	_, err = s.cron.AddFunc(morningCron, func() {
		s.sendScheduledReport("morning")
	})
	if err != nil {
		return fmt.Errorf("failed to schedule morning report: %w", err)
	}
	s.logger.Info("Morning report scheduled",
		slog.String("time", s.morningTime),
		slog.String("cron", morningCron),
		slog.String("timezone", s.timezone.String()),
	)

	// Schedule evening report
	eveningCron := fmt.Sprintf("%d %d * * *", eveningMin, eveningHour)
	_, err = s.cron.AddFunc(eveningCron, func() {
		s.sendScheduledReport("evening")
	})
	if err != nil {
		return fmt.Errorf("failed to schedule evening report: %w", err)
	}
	s.logger.Info("Evening report scheduled",
		slog.String("time", s.eveningTime),
		slog.String("cron", eveningCron),
		slog.String("timezone", s.timezone.String()),
	)

	// Start the cron scheduler
	s.cron.Start()
	s.logger.Info("Report scheduler started",
		slog.Int("recipients", len(s.recipients)),
	)

	return nil
}

// Stop stops the report scheduler
func (s *ReportScheduler) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		s.logger.Info("Report scheduler stopped")
	}
}

// sendScheduledReport generates and sends a scheduled report
func (s *ReportScheduler) sendScheduledReport(reportType string) {
	s.logger.Info("Sending scheduled report",
		slog.String("type", reportType),
		slog.Int("recipients", len(s.recipients)),
	)

	ctx := context.Background()

	// Generate report
	report, err := s.reportService.GenerateDailyReport(ctx, reportType)
	if err != nil {
		s.logger.Error("Failed to generate daily report",
			slog.String("type", reportType),
			slog.String("error", err.Error()),
		)
		return
	}

	// Send email
	err = SendDailyReportEmail(ctx, s.sendgridAPIKey, s.fromEmail, s.fromName, report, s.recipients)
	if err != nil {
		s.logger.Error("Failed to send daily report email",
			slog.String("type", reportType),
			slog.String("error", err.Error()),
		)
		return
	}

	s.logger.Info("Daily report sent successfully",
		slog.String("type", reportType),
		slog.Int("recipients", len(s.recipients)),
		slog.Int("total_tickets", report.TotalTickets),
		slog.Int("new_today", report.NewTicketsToday),
	)
}

// SendNow sends a report immediately (for testing or manual trigger)
func (s *ReportScheduler) SendNow(reportType string) error {
	s.logger.Info("Sending report on demand",
		slog.String("type", reportType),
	)

	ctx := context.Background()

	// Generate report
	report, err := s.reportService.GenerateDailyReport(ctx, reportType)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Send email
	err = SendDailyReportEmail(ctx, s.sendgridAPIKey, s.fromEmail, s.fromName, report, s.recipients)
	if err != nil {
		return fmt.Errorf("failed to send report email: %w", err)
	}

	s.logger.Info("On-demand report sent successfully",
		slog.String("type", reportType),
	)

	return nil
}

// UpdateSchedule updates the report schedule
func (s *ReportScheduler) UpdateSchedule(morningTime, eveningTime string) error {
	s.Stop()
	s.morningTime = morningTime
	s.eveningTime = eveningTime
	return s.Start()
}

// UpdateRecipients updates the recipient list
func (s *ReportScheduler) UpdateRecipients(recipients []string) {
	s.recipients = recipients
	s.logger.Info("Report recipients updated",
		slog.Int("count", len(recipients)),
	)
}

// GetScheduleInfo returns current schedule information
func (s *ReportScheduler) GetScheduleInfo() map[string]interface{} {
	return map[string]interface{}{
		"enabled":      s.featureFlags.DailyReportsEnabled,
		"morning_time": s.morningTime,
		"evening_time": s.eveningTime,
		"timezone":     s.timezone.String(),
		"recipients":   len(s.recipients),
	}
}

// Helper function to parse time string (HH:MM)
func parseTime(timeStr string) (hour, minute int, err error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format: %s (expected HH:MM)", timeStr)
	}

	_, err = fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse time: %w", err)
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("invalid time values: %d:%d", hour, minute)
	}

	return hour, minute, nil
}
