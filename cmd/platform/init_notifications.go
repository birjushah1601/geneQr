package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/aby-med/medical-platform/internal/infrastructure/config"
	"github.com/aby-med/medical-platform/internal/infrastructure/email"
	"github.com/aby-med/medical-platform/internal/infrastructure/notification"
	"github.com/aby-med/medical-platform/internal/infrastructure/reports"
	"github.com/jmoiron/sqlx"
)

// initNotificationsAndReports initializes email notifications and daily reports systems
func initNotificationsAndReports(ctx context.Context, db *sqlx.DB, logger *slog.Logger) (*notification.Manager, *reports.ReportScheduler, error) {
	logger.Info("Initializing Notifications and Reports Systems")

	// Load feature flags
	featureFlags := config.LoadFeatureFlags()

	// ========================================================================
	// INITIALIZE EMAIL NOTIFICATION SERVICE
	// ========================================================================
	
	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	sendgridFromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	sendgridFromName := os.Getenv("SENDGRID_FROM_NAME")

	if sendgridFromEmail == "" {
		sendgridFromEmail = "noreply@ServQR.com"
	}
	if sendgridFromName == "" {
		sendgridFromName = "ServQR Platform"
	}

	var notificationManager *notification.Manager
	if sendgridAPIKey == "" {
		logger.Warn("SendGrid API key not configured - email notifications disabled",
			slog.Bool("email_notifications_enabled", featureFlags.EmailNotificationsEnabled))
		// Create a nil notification manager to prevent crashes
		notificationManager = nil
	} else {
		// Initialize email service
		emailService := email.NewNotificationService(
			sendgridAPIKey,
			sendgridFromEmail,
			sendgridFromName,
		)

		// Get admin email
		adminEmail := os.Getenv("ADMIN_EMAIL")
		if adminEmail == "" {
			adminEmail = "admin@ServQR.com"
			logger.Warn("ADMIN_EMAIL not set, using default", slog.String("admin_email", adminEmail))
		}

		// Initialize notification manager
		notificationManager = notification.NewManager(emailService, featureFlags, logger, adminEmail)

		logger.Info("âœ… Email Notification Service initialized",
			slog.String("from_email", sendgridFromEmail),
			slog.String("admin_email", adminEmail),
			slog.Bool("email_enabled", featureFlags.EmailNotificationsEnabled),
			slog.Bool("ticket_created", featureFlags.EmailTicketCreatedEnabled),
			slog.Bool("engineer_assigned", featureFlags.EmailEngineerAssignedEnabled),
			slog.Bool("status_changed", featureFlags.EmailStatusChangedEnabled),
		)
	}

	// ========================================================================
	// INITIALIZE DAILY REPORTS SYSTEM
	// ========================================================================

	var reportScheduler *reports.ReportScheduler

	if !featureFlags.DailyReportsEnabled {
		logger.Info("Daily reports disabled by feature flag")
		reportScheduler = nil
	} else {
		// Get report configuration from environment
		morningTime := os.Getenv("DAILY_REPORT_MORNING_TIME")
		if morningTime == "" {
			morningTime = "09:00"
		}

		eveningTime := os.Getenv("DAILY_REPORT_EVENING_TIME")
		if eveningTime == "" {
			eveningTime = "18:00"
		}

		timezone := os.Getenv("DAILY_REPORT_TIMEZONE")
		if timezone == "" {
			timezone = "UTC"
		}

		// Parse recipients
		recipientsStr := os.Getenv("DAILY_REPORT_RECIPIENTS")
		var recipients []string
		if recipientsStr != "" {
			recipients = strings.Split(recipientsStr, ",")
			for i := range recipients {
				recipients[i] = strings.TrimSpace(recipients[i])
			}
		}

		if len(recipients) == 0 {
			logger.Warn("No recipients configured for daily reports")
			reportScheduler = nil
		} else if sendgridAPIKey == "" {
			logger.Warn("SendGrid API key not configured - daily reports disabled")
			reportScheduler = nil
		} else {
			// Initialize report service (convert sqlx.DB to sql.DB)
			reportService := reports.NewDailyReportService(db.DB, logger)

			// Create report scheduler
			scheduler, err := reports.NewReportScheduler(
				reportService,
				sendgridAPIKey,
				sendgridFromEmail,
				sendgridFromName,
				featureFlags,
				logger,
				morningTime,
				eveningTime,
				recipients,
				timezone,
			)

			if err != nil {
				logger.Error("Failed to create report scheduler", slog.String("error", err.Error()))
				reportScheduler = nil
			} else {
				reportScheduler = scheduler

				logger.Info("âœ… Daily Reports System initialized",
					slog.String("morning_time", morningTime),
					slog.String("evening_time", eveningTime),
					slog.String("timezone", timezone),
					slog.Int("recipients", len(recipients)),
					slog.Bool("morning_enabled", featureFlags.DailyReportMorningEnabled),
					slog.Bool("evening_enabled", featureFlags.DailyReportEveningEnabled),
				)
			}
		}
	}

	// ========================================================================
	// START REPORT SCHEDULER
	// ========================================================================

	if reportScheduler != nil {
		if err := reportScheduler.Start(); err != nil {
			logger.Error("Failed to start report scheduler", slog.String("error", err.Error()))
			// Don't return error - continue without reports
			reportScheduler = nil
		} else {
			logger.Info("âœ… Report scheduler started successfully")
		}
	}

	return notificationManager, reportScheduler, nil
}
