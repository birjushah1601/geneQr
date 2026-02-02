package sms

import (
	"context"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// TwilioSender implements SMS/WhatsApp sending via Twilio
type TwilioSender struct {
	client           *twilio.RestClient
	phoneNumber      string
	whatsappNumber   string
}

// NewTwilioSender creates a new Twilio SMS sender
func NewTwilioSender(accountSID, authToken, phoneNumber, whatsappNumber string) *TwilioSender {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &TwilioSender{
		client:         client,
		phoneNumber:    phoneNumber,
		whatsappNumber: whatsappNumber,
	}
}

// SendOTP sends an OTP code via SMS
func (s *TwilioSender) SendOTP(ctx context.Context, to, otp string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.phoneNumber)
	params.SetBody(fmt.Sprintf("Your ServQR verification code is: %s\n\nThis code will expire in 5 minutes.", otp))

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}

// SendWhatsAppOTP sends an OTP code via WhatsApp
func (s *TwilioSender) SendWhatsAppOTP(ctx context.Context, to, otp string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(fmt.Sprintf("whatsapp:%s", to))
	params.SetFrom(fmt.Sprintf("whatsapp:%s", s.whatsappNumber))
	params.SetBody(fmt.Sprintf("*ServQR Verification Code*\n\nYour code: *%s*\n\nThis code will expire in 5 minutes.\n\nIf you didn't request this, please ignore.", otp))

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	return nil
}
