package notification

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.com/tour/internal/config"
	"gitlab.com/tour/internal/pkg/security"
	"go.uber.org/zap"
)

type EmailNotificationService struct {
	config *config.EmailNotificationConfig
	redis  *redis.Client
	logger *zap.Logger
}

func NewEmailNotificationService(config *config.Config, redisClient *redis.Client, logger *zap.Logger) *EmailNotificationService {
	return &EmailNotificationService{
		config: &config.EmailNotificationConfig,
		redis:  redisClient,
		logger: logger,
	}
}

func (e *EmailNotificationService) SendVerificationCode(toAddr string) {
	randomCode, err := security.GenerateRandomStringOrdinary(5)
	if err != nil {
		e.logger.Error("Failed to generate random string", zap.Error(err))
	}

	codeLiveDuration, err := time.ParseDuration(e.config.RandomCodeLiveTime)
	if err != nil {
		e.logger.Error("Failed to parse duration", zap.Error(err))
	}

	err = e.redis.Set(context.Background(), toAddr, randomCode, codeLiveDuration).Err()
	if err != nil {
		e.logger.Error("Failed to set redis", zap.Error(err))
	}

	subject := fmt.Sprintf("%s - Email Verification Code", e.config.ProductName)
	body := fmt.Sprintf("Your verification code is %s", randomCode)
	e.SendEmail(toAddr, subject, body)
}

func (e *EmailNotificationService) SendEmail(toAddr, subject, body string) {
	// Receiver email address.
	to := []string{toAddr}

	message := "From: " + fmt.Sprintf("%s <%s>", e.config.ProductName, e.config.SenderEmail) + "\r\n"
	message += "To: " + toAddr + "\r\n"
	message += "Subject: " + subject + "\r\n\r\n"
	message += body

	// Message.
	messageBytes := []byte(message)

	// Authentication.
	auth := smtp.PlainAuth("", e.config.SenderEmail, e.config.SenderPassword, e.config.SmtpHost)

	// Sending email.
	err := smtp.SendMail(e.config.SmtpHost+":"+e.config.SmtpPort, auth, e.config.SenderEmail, to, messageBytes)
	if err != nil {
		e.logger.Error("Failed to send email", zap.Error(err))
	}
}
