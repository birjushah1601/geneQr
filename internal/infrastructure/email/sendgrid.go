package email

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridSender implements email sending via SendGrid
type SendGridSender struct {
	apiKey    string
	fromEmail string
	fromName  string
}

// NewSendGridSender creates a new SendGrid email sender
func NewSendGridSender(apiKey, fromEmail, fromName string) *SendGridSender {
	return &SendGridSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// SendOTP sends an OTP code via email
func (s *SendGridSender) SendOTP(ctx context.Context, to, otp string) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	subject := "Your Verification Code"
	toEmail := mail.NewEmail("", to)
	
	// Plain text content
	plainTextContent := fmt.Sprintf(`
Your verification code is: %s

This code will expire in 5 minutes.

If you didn't request this code, please ignore this email.

Best regards,
ServQR Platform
`, otp)

	// HTML content
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 5px 5px; }
        .otp-code { background-color: #fff; border: 2px solid #2563eb; padding: 20px; font-size: 32px; font-weight: bold; text-align: center; letter-spacing: 8px; margin: 20px 0; border-radius: 5px; color: #2563eb; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>ServQR Platform</h1>
        </div>
        <div class="content">
            <h2>Your Verification Code</h2>
            <p>Use the following code to complete your authentication:</p>
            <div class="otp-code">%s</div>
            <p><strong>This code will expire in 5 minutes.</strong></p>
            <p>If you didn't request this code, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>This is an automated message from ServQR Platform</p>
            <p>&copy; 2025 ServQR. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, otp)

	message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
