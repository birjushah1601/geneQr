package config

import (
	"os"
	"strconv"
	"strings"
)

// FeatureFlags manages feature toggles for the application
type FeatureFlags struct {
	// Email Notifications
	EmailNotificationsEnabled       bool
	EmailTicketCreatedEnabled       bool
	EmailEngineerAssignedEnabled    bool
	EmailStatusChangedEnabled       bool
	
	// SMS Notifications (future)
	SMSNotificationsEnabled         bool
	SMSTicketCreatedEnabled         bool
	SMSEngineerAssignedEnabled      bool
	SMSStatusChangedEnabled         bool
	
	// WhatsApp Notifications (future)
	WhatsAppNotificationsEnabled    bool
	WhatsAppTicketCreatedEnabled    bool
	WhatsAppEngineerAssignedEnabled bool
	WhatsAppStatusChangedEnabled    bool
	
	// Daily Reports
	DailyReportsEnabled             bool
	DailyReportMorningEnabled       bool
	DailyReportEveningEnabled       bool
	
	// Other Features
	AIAnalysisEnabled               bool
	MultiModelAssignmentEnabled     bool
	AuditLoggingEnabled             bool
}

// LoadFeatureFlags loads feature flags from environment variables
func LoadFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		// Email Notifications - Master switch and individual toggles
		EmailNotificationsEnabled:       getBoolEnv("FEATURE_EMAIL_NOTIFICATIONS", false),
		EmailTicketCreatedEnabled:       getBoolEnv("FEATURE_EMAIL_TICKET_CREATED", false),
		EmailEngineerAssignedEnabled:    getBoolEnv("FEATURE_EMAIL_ENGINEER_ASSIGNED", false),
		EmailStatusChangedEnabled:       getBoolEnv("FEATURE_EMAIL_STATUS_CHANGED", false),
		
		// SMS Notifications - Future
		SMSNotificationsEnabled:         getBoolEnv("FEATURE_SMS_NOTIFICATIONS", false),
		SMSTicketCreatedEnabled:         getBoolEnv("FEATURE_SMS_TICKET_CREATED", false),
		SMSEngineerAssignedEnabled:      getBoolEnv("FEATURE_SMS_ENGINEER_ASSIGNED", false),
		SMSStatusChangedEnabled:         getBoolEnv("FEATURE_SMS_STATUS_CHANGED", false),
		
		// WhatsApp Notifications - Future
		WhatsAppNotificationsEnabled:    getBoolEnv("FEATURE_WHATSAPP_NOTIFICATIONS", false),
		WhatsAppTicketCreatedEnabled:    getBoolEnv("FEATURE_WHATSAPP_TICKET_CREATED", false),
		WhatsAppEngineerAssignedEnabled: getBoolEnv("FEATURE_WHATSAPP_ENGINEER_ASSIGNED", false),
		WhatsAppStatusChangedEnabled:    getBoolEnv("FEATURE_WHATSAPP_STATUS_CHANGED", false),
		
		// Daily Reports
		DailyReportsEnabled:             getBoolEnv("FEATURE_DAILY_REPORTS", false),
		DailyReportMorningEnabled:       getBoolEnv("FEATURE_DAILY_REPORT_MORNING", false),
		DailyReportEveningEnabled:       getBoolEnv("FEATURE_DAILY_REPORT_EVENING", false),
		
		// Other Features
		AIAnalysisEnabled:               getBoolEnv("FEATURE_AI_ANALYSIS", true),
		MultiModelAssignmentEnabled:     getBoolEnv("FEATURE_MULTI_MODEL_ASSIGNMENT", true),
		AuditLoggingEnabled:             getBoolEnv("FEATURE_AUDIT_LOGGING", true),
	}
}

// ShouldSendEmailNotification checks if email notification should be sent for a specific event
func (f *FeatureFlags) ShouldSendEmailNotification(eventType string) bool {
	// Master switch must be enabled
	if !f.EmailNotificationsEnabled {
		return false
	}
	
	// Check individual event type
	switch strings.ToLower(eventType) {
	case "ticket_created":
		return f.EmailTicketCreatedEnabled
	case "engineer_assigned":
		return f.EmailEngineerAssignedEnabled
	case "status_changed":
		return f.EmailStatusChangedEnabled
	default:
		return false
	}
}

// ShouldSendSMSNotification checks if SMS notification should be sent for a specific event
func (f *FeatureFlags) ShouldSendSMSNotification(eventType string) bool {
	// Master switch must be enabled
	if !f.SMSNotificationsEnabled {
		return false
	}
	
	// Check individual event type
	switch strings.ToLower(eventType) {
	case "ticket_created":
		return f.SMSTicketCreatedEnabled
	case "engineer_assigned":
		return f.SMSEngineerAssignedEnabled
	case "status_changed":
		return f.SMSStatusChangedEnabled
	default:
		return false
	}
}

// ShouldSendWhatsAppNotification checks if WhatsApp notification should be sent for a specific event
func (f *FeatureFlags) ShouldSendWhatsAppNotification(eventType string) bool {
	// Master switch must be enabled
	if !f.WhatsAppNotificationsEnabled {
		return false
	}
	
	// Check individual event type
	switch strings.ToLower(eventType) {
	case "ticket_created":
		return f.WhatsAppTicketCreatedEnabled
	case "engineer_assigned":
		return f.WhatsAppEngineerAssignedEnabled
	case "status_changed":
		return f.WhatsAppStatusChangedEnabled
	default:
		return false
	}
}

// GetFeatureFlagsStatus returns a map of all feature flags and their status
func (f *FeatureFlags) GetFeatureFlagsStatus() map[string]bool {
	return map[string]bool{
		// Email
		"email_notifications":        f.EmailNotificationsEnabled,
		"email_ticket_created":       f.EmailTicketCreatedEnabled,
		"email_engineer_assigned":    f.EmailEngineerAssignedEnabled,
		"email_status_changed":       f.EmailStatusChangedEnabled,
		
		// SMS
		"sms_notifications":          f.SMSNotificationsEnabled,
		"sms_ticket_created":         f.SMSTicketCreatedEnabled,
		"sms_engineer_assigned":      f.SMSEngineerAssignedEnabled,
		"sms_status_changed":         f.SMSStatusChangedEnabled,
		
		// WhatsApp
		"whatsapp_notifications":     f.WhatsAppNotificationsEnabled,
		"whatsapp_ticket_created":    f.WhatsAppTicketCreatedEnabled,
		"whatsapp_engineer_assigned": f.WhatsAppEngineerAssignedEnabled,
		"whatsapp_status_changed":    f.WhatsAppStatusChangedEnabled,
		
		// Daily Reports
		"daily_reports":              f.DailyReportsEnabled,
		"daily_report_morning":       f.DailyReportMorningEnabled,
		"daily_report_evening":       f.DailyReportEveningEnabled,
		
		// Other
		"ai_analysis":                f.AIAnalysisEnabled,
		"multi_model_assignment":     f.MultiModelAssignmentEnabled,
		"audit_logging":              f.AuditLoggingEnabled,
	}
}

// Helper function to get boolean environment variable
func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	// Parse boolean
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		// Try common string representations
		valueLower := strings.ToLower(strings.TrimSpace(value))
		switch valueLower {
		case "true", "yes", "y", "1", "on", "enabled":
			return true
		case "false", "no", "n", "0", "off", "disabled":
			return false
		default:
			return defaultValue
		}
	}
	
	return boolValue
}
