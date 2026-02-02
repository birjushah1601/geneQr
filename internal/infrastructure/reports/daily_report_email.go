package reports

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendDailyReportEmail sends the daily report email to admins
func SendDailyReportEmail(ctx context.Context, apiKey, fromEmail, fromName string, report *DailyReportData, recipients []string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients provided for daily report")
	}

	from := mail.NewEmail(fromName, fromEmail)
	
	// Determine subject based on report type
	var subject string
	if report.ReportType == "morning" {
		subject = fmt.Sprintf("Morning Report - %s", report.ReportDate.Format("Jan 02, 2006"))
	} else {
		subject = fmt.Sprintf("Evening Report - %s", report.ReportDate.Format("Jan 02, 2006"))
	}

	// Generate email content
	plainText := generateDailyReportPlainText(report)
	htmlContent := generateDailyReportHTML(report)

	// Send to each recipient
	client := sendgrid.NewSendClient(apiKey)
	
	for _, recipientEmail := range recipients {
		to := mail.NewEmail("Admin", recipientEmail)
		message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)

		response, err := client.Send(message)
		if err != nil {
			return fmt.Errorf("failed to send email to %s: %w", recipientEmail, err)
		}

		if response.StatusCode >= 400 {
			return fmt.Errorf("sendgrid error for %s: status %d, body: %s", recipientEmail, response.StatusCode, response.Body)
		}
	}

	return nil
}

// generateDailyReportPlainText generates plain text version of the report
func generateDailyReportPlainText(report *DailyReportData) string {
	reportTime := "Morning"
	if report.ReportType == "evening" {
		reportTime = "Evening"
	}

	text := fmt.Sprintf(`
ServQR Platform - %s Daily Report
Date: %s
========================================

TICKET SUMMARY
========================================
Total Tickets: %d
New Today: %d
Resolved Today: %d
Pending: %d
In Progress: %d
On Hold: %d

PRIORITY BREAKDOWN (Open Tickets)
========================================
Critical: %d
High: %d
Medium: %d
Low: %d

ENGINEER STATISTICS
========================================
Total Engineers: %d
Active Engineers: %d
Engineers with Assignments: %d
Avg Tickets per Engineer: %.1f

EQUIPMENT STATISTICS
========================================
Total Equipment: %d
Equipment with Issues: %d
Equipment Serviced Today: %d

PERFORMANCE METRICS
========================================
Avg Resolution Time: %.1f hours
Tickets within SLA: %d
Overdue Tickets: %d

`, reportTime, report.ReportDate.Format("January 02, 2006 15:04"),
		report.TotalTickets, report.NewTicketsToday, report.ResolvedTicketsToday,
		report.PendingTickets, report.InProgressTickets, report.OnHoldTickets,
		report.CriticalTickets, report.HighPriorityTickets, report.MediumPriorityTickets, report.LowPriorityTickets,
		report.TotalEngineers, report.ActiveEngineers, report.EngineersWithTickets, report.AverageTicketsPerEngineer,
		report.TotalEquipment, report.EquipmentWithIssues, report.EquipmentServiced,
		report.AverageResolutionTime, report.TicketsSLA, report.TicketsOverdue)

	// Top Engineers
	if len(report.TopEngineers) > 0 {
		text += "TOP PERFORMING ENGINEERS\n========================================\n"
		for i, eng := range report.TopEngineers {
			text += fmt.Sprintf("%d. %s - %d resolved, %d assigned (Avg: %.1fh)\n",
				i+1, eng.EngineerName, eng.TicketsResolved, eng.TicketsAssigned, eng.AverageResolution)
		}
		text += "\n"
	}

	// Equipment with most issues
	if len(report.TopEquipment) > 0 {
		text += "EQUIPMENT WITH MOST ISSUES\n========================================\n"
		for i, eq := range report.TopEquipment {
			text += fmt.Sprintf("%d. %s (%s) - %d issues\n",
				i+1, eq.EquipmentName, eq.Manufacturer, eq.IssueCount)
		}
		text += "\n"
	}

	// Tickets needing attention
	if len(report.TicketsNeedingAttention) > 0 {
		text += "âš ï¸  TICKETS NEEDING ATTENTION\n========================================\n"
		for i, alert := range report.TicketsNeedingAttention {
			text += fmt.Sprintf("%d. %s - %s (%s) - %s [%d days open]\n",
				i+1, alert.TicketNumber, alert.EquipmentName, alert.Priority, alert.Reason, alert.DaysOpen)
		}
		text += "\n"
	}

	text += `
========================================
ServQR Admin System
`

	return text
}

// generateDailyReportHTML generates HTML version of the report
func generateDailyReportHTML(report *DailyReportData) string {
	reportTime := "Morning"
	headerColor := "#f59e0b" // Orange for morning
	if report.ReportType == "evening" {
		reportTime = "Evening"
		headerColor = "#6366f1" // Purple for evening
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; background-color: #f3f4f6; margin: 0; padding: 20px; }
        .container { max-width: 800px; margin: 0 auto; background-color: #fff; }
        .header { background-color: %s; color: white; padding: 30px; text-align: center; }
        .header h1 { margin: 0; font-size: 28px; }
        .header p { margin: 5px 0 0 0; opacity: 0.9; }
        .content { padding: 30px; }
        .section { margin-bottom: 30px; }
        .section-title { font-size: 18px; font-weight: bold; color: #1f2937; margin-bottom: 15px; padding-bottom: 5px; border-bottom: 2px solid #e5e7eb; }
        .stats-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 15px; margin-bottom: 20px; }
        .stat-card { background-color: #f9fafb; border-left: 4px solid %s; padding: 15px; border-radius: 4px; }
        .stat-label { font-size: 12px; color: #6b7280; text-transform: uppercase; margin-bottom: 5px; }
        .stat-value { font-size: 24px; font-weight: bold; color: #111827; }
        .stat-subtext { font-size: 12px; color: #6b7280; margin-top: 5px; }
        .priority-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; }
        .priority-card { padding: 10px; border-radius: 4px; text-align: center; }
        .priority-critical { background-color: #fef2f2; border: 2px solid #dc2626; }
        .priority-high { background-color: #fff7ed; border: 2px solid #f59e0b; }
        .priority-medium { background-color: #fffbeb; border: 2px solid #fbbf24; }
        .priority-low { background-color: #f0fdf4; border: 2px solid: #10b981; }
        .priority-label { font-size: 11px; text-transform: uppercase; font-weight: bold; }
        .priority-value { font-size: 20px; font-weight: bold; margin-top: 5px; }
        .table { width: 100%%; border-collapse: collapse; margin-top: 10px; }
        .table th { background-color: #f3f4f6; padding: 10px; text-align: left; font-size: 12px; font-weight: bold; border-bottom: 2px solid #e5e7eb; }
        .table td { padding: 10px; border-bottom: 1px solid #e5e7eb; font-size: 13px; }
        .badge { display: inline-block; padding: 3px 8px; border-radius: 3px; font-size: 11px; font-weight: bold; }
        .badge-critical { background-color: #fee2e2; color: #dc2626; }
        .badge-high { background-color: #fed7aa; color: #f59e0b; }
        .badge-medium { background-color: #fef3c7; color: #f59e0b; }
        .badge-low { background-color: #d1fae5; color: #10b981; }
        .alert-box { background-color: #fef2f2; border-left: 4px solid #dc2626; padding: 15px; border-radius: 4px; margin-top: 10px; }
        .alert-item { padding: 8px 0; border-bottom: 1px solid #fecaca; }
        .alert-item:last-child { border-bottom: none; }
        .footer { background-color: #f9fafb; padding: 20px; text-align: center; font-size: 12px; color: #6b7280; border-top: 1px solid #e5e7eb; }
        .metric-good { color: #10b981; font-weight: bold; }
        .metric-warning { color: #f59e0b; font-weight: bold; }
        .metric-bad { color: #dc2626; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ðŸ“Š %s Daily Report</h1>
            <p>%s</p>
        </div>
        
        <div class="content">
            <!-- Ticket Summary -->
            <div class="section">
                <div class="section-title">ðŸ“‹ Ticket Summary</div>
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-label">Total Tickets</div>
                        <div class="stat-value">%d</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">New Today</div>
                        <div class="stat-value" style="color: #3b82f6;">%d</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Resolved Today</div>
                        <div class="stat-value" style="color: #10b981;">%d</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Pending</div>
                        <div class="stat-value">%d</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">In Progress</div>
                        <div class="stat-value" style="color: #f59e0b;">%d</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">On Hold</div>
                        <div class="stat-value" style="color: #6b7280;">%d</div>
                    </div>
                </div>
            </div>

            <!-- Priority Breakdown -->
            <div class="section">
                <div class="section-title">ðŸ”¥ Priority Breakdown (Open Tickets)</div>
                <div class="priority-grid">
                    <div class="priority-card priority-critical">
                        <div class="priority-label" style="color: #dc2626;">Critical</div>
                        <div class="priority-value" style="color: #dc2626;">%d</div>
                    </div>
                    <div class="priority-card priority-high">
                        <div class="priority-label" style="color: #f59e0b;">High</div>
                        <div class="priority-value" style="color: #f59e0b;">%d</div>
                    </div>
                    <div class="priority-card priority-medium">
                        <div class="priority-label" style="color: #fbbf24;">Medium</div>
                        <div class="priority-value" style="color: #fbbf24;">%d</div>
                    </div>
                    <div class="priority-card priority-low">
                        <div class="priority-label" style="color: #10b981;">Low</div>
                        <div class="priority-value" style="color: #10b981;">%d</div>
                    </div>
                </div>
            </div>

            <!-- Performance Metrics -->
            <div class="section">
                <div class="section-title">âš¡ Performance Metrics</div>
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-label">Avg Resolution Time</div>
                        <div class="stat-value">%.1fh</div>
                        <div class="stat-subtext">Last 7 days</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Within SLA</div>
                        <div class="stat-value metric-good">%d</div>
                        <div class="stat-subtext">Resolved today</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Overdue Tickets</div>
                        <div class="stat-value metric-bad">%d</div>
                        <div class="stat-subtext">&gt;48 hours</div>
                    </div>
                </div>
            </div>

            <!-- Engineer Statistics -->
            <div class="section">
                <div class="section-title">ðŸ‘· Engineer Statistics</div>
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-label">Total Engineers</div>
                        <div class="stat-value">%d</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Active Engineers</div>
                        <div class="stat-value">%d</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Avg Tickets/Engineer</div>
                        <div class="stat-value">%.1f</div>
                    </div>
                </div>
            </div>
`, headerColor, headerColor, reportTime, report.ReportDate.Format("January 02, 2006 at 3:04 PM"),
		report.TotalTickets, report.NewTicketsToday, report.ResolvedTicketsToday,
		report.PendingTickets, report.InProgressTickets, report.OnHoldTickets,
		report.CriticalTickets, report.HighPriorityTickets, report.MediumPriorityTickets, report.LowPriorityTickets,
		report.AverageResolutionTime, report.TicketsSLA, report.TicketsOverdue,
		report.TotalEngineers, report.ActiveEngineers, report.AverageTicketsPerEngineer)

	// Top Engineers
	if len(report.TopEngineers) > 0 {
		html += `
            <div class="section">
                <div class="section-title">ðŸ† Top Performing Engineers</div>
                <table class="table">
                    <thead>
                        <tr>
                            <th>Rank</th>
                            <th>Engineer</th>
                            <th>Resolved</th>
                            <th>Assigned</th>
                            <th>Avg Resolution</th>
                        </tr>
                    </thead>
                    <tbody>`
		
		for i, eng := range report.TopEngineers {
			html += fmt.Sprintf(`
                        <tr>
                            <td><strong>%d</strong></td>
                            <td>%s</td>
                            <td><span class="metric-good">%d</span></td>
                            <td>%d</td>
                            <td>%.1f hours</td>
                        </tr>`, i+1, eng.EngineerName, eng.TicketsResolved, eng.TicketsAssigned, eng.AverageResolution)
		}
		
		html += `
                    </tbody>
                </table>
            </div>`
	}

	// Equipment with most issues
	if len(report.TopEquipment) > 0 {
		html += `
            <div class="section">
                <div class="section-title">ðŸ”§ Equipment with Most Issues</div>
                <table class="table">
                    <thead>
                        <tr>
                            <th>Rank</th>
                            <th>Equipment</th>
                            <th>Manufacturer</th>
                            <th>Issues (30 days)</th>
                        </tr>
                    </thead>
                    <tbody>`
		
		for i, eq := range report.TopEquipment {
			html += fmt.Sprintf(`
                        <tr>
                            <td><strong>%d</strong></td>
                            <td>%s</td>
                            <td>%s</td>
                            <td><span class="metric-warning">%d</span></td>
                        </tr>`, i+1, eq.EquipmentName, eq.Manufacturer, eq.IssueCount)
		}
		
		html += `
                    </tbody>
                </table>
            </div>`
	}

	// Tickets needing attention
	if len(report.TicketsNeedingAttention) > 0 {
		html += `
            <div class="section">
                <div class="section-title">âš ï¸ Tickets Needing Attention</div>`
		
		for _, alert := range report.TicketsNeedingAttention {
			priorityClass := "badge-medium"
			switch alert.Priority {
			case "critical":
				priorityClass = "badge-critical"
			case "high":
				priorityClass = "badge-high"
			case "low":
				priorityClass = "badge-low"
			}
			
			html += fmt.Sprintf(`
                <div class="alert-box">
                    <div class="alert-item">
                        <strong>%s</strong> - %s
                        <span class="badge %s">%s</span>
                        <div style="font-size: 12px; color: #6b7280; margin-top: 5px;">
                            Customer: %s | %s | %d days open
                        </div>
                    </div>
                </div>`, alert.TicketNumber, alert.EquipmentName, priorityClass, alert.Priority, 
				alert.CustomerName, alert.Reason, alert.DaysOpen)
		}
		
		html += `
            </div>`
	}

	html += `
        </div>
        
        <div class="footer">
            <p><strong>ServQR Admin Notification System</strong></p>
            <p>This is an automated daily report | &copy; 2025 ServQR. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	return html
}
