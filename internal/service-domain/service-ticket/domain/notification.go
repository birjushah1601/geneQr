package domain

import (
	"time"
)

// NotificationType represents the type of notification sent
type NotificationType string

const (
	NotificationTypeManual        NotificationType = "manual"
	NotificationTypeTicketCreated NotificationType = "ticket_created"
	NotificationTypeDailyDigest   NotificationType = "daily_digest"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusSent   NotificationStatus = "sent"
	NotificationStatusFailed NotificationStatus = "failed"
)

// TrackingToken represents a public tracking token for a ticket
type TrackingToken struct {
	ID        string    `json:"id" db:"id"`
	TicketID  string    `json:"ticket_id" db:"ticket_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// NotificationLog represents a record of sent notifications
type NotificationLog struct {
	ID               string             `json:"id" db:"id"`
	TicketID         string             `json:"ticket_id" db:"ticket_id"`
	NotificationType NotificationType   `json:"notification_type" db:"notification_type"`
	RecipientEmail   string             `json:"recipient_email" db:"recipient_email"`
	SentAt           time.Time          `json:"sent_at" db:"sent_at"`
	Status           NotificationStatus `json:"status" db:"status"`
	ErrorMessage     *string            `json:"error_message,omitempty" db:"error_message"`
}

// PublicTicketView represents the public-facing view of a ticket
type PublicTicketView struct {
	TicketNumber      string              `json:"ticket_number"`
	Status            string              `json:"status"`
	Priority          string              `json:"priority"`
	EquipmentName     string              `json:"equipment_name"`
	IssueDescription  string              `json:"issue_description"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	PublicComments    []PublicComment     `json:"comments"`
	StatusHistory     []PublicStatusEvent `json:"status_history"`
	AssignedEngineer  string              `json:"assigned_engineer,omitempty"`
}

// PublicComment represents a comment visible to customers
type PublicComment struct {
	Comment    string    `json:"comment"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	AuthorRole string    `json:"author_role"` // 'system', 'engineer', 'admin'
}

// PublicStatusEvent represents a status change event
type PublicStatusEvent struct {
	FromStatus string    `json:"from_status,omitempty"`
	ToStatus   string    `json:"to_status"`
	ChangedBy  string    `json:"changed_by"`
	ChangedAt  time.Time `json:"changed_at"`
	Comment    string    `json:"comment,omitempty"`
}

// EmailTemplateData contains all data for email templates
type EmailTemplateData struct {
	// Ticket Info
	TicketNumber          string    `json:"ticket_number"`
	TicketID              string    `json:"ticket_id"`
	Status                string    `json:"status"`
	Priority              string    `json:"priority"`
	
	// Customer Info
	CustomerName          string    `json:"customer_name"`
	CustomerEmail         string    `json:"customer_email"`
	
	// Equipment Info
	EquipmentName         string    `json:"equipment_name"`
	SerialNumber          string    `json:"serial_number,omitempty"`
	
	// Issue Details
	IssueDescription      string    `json:"issue_description"`
	IssueCategory         string    `json:"issue_category,omitempty"`
	
	// Assignment Info
	AssignedEngineerName  string    `json:"assigned_engineer_name,omitempty"`
	AssignedAt            *time.Time `json:"assigned_at,omitempty"`
	
	// Dates
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	
	// Comments
	Comments              []CommentData `json:"comments"`
	
	// Tracking
	TrackingURL           string    `json:"tracking_url"`
	TrackingToken         string    `json:"tracking_token"`
	
	// Status Updates
	StatusHistory         []StatusUpdate `json:"status_history,omitempty"`
}

// CommentData represents comment data for templates
type CommentData struct {
	Comment   string    `json:"comment"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	IsPublic  bool      `json:"is_public"` // Only public comments in emails
}

// StatusUpdate represents a status change
type StatusUpdate struct {
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	ChangedAt  time.Time `json:"changed_at"`
	ChangedBy  string    `json:"changed_by"`
}

// NotificationRepository defines methods for notification persistence
type NotificationRepository interface {
	// Tracking Tokens
	CreateTrackingToken(token *TrackingToken) error
	GetTrackingToken(token string) (*TrackingToken, error)
	DeleteExpiredTokens() error
	
	// Notification Log
	LogNotification(log *NotificationLog) error
	GetNotificationHistory(ticketID string) ([]*NotificationLog, error)
	GetFailedNotifications(since time.Time) ([]*NotificationLog, error)
}
