package notification

import (
	"fmt"

	"github.com/resend/resend-go/v2"
	"github.com/EstebanGitPro/motogo-backend/internal/domain/person"
	errs "github.com/EstebanGitPro/motogo-backend/internal/domain/errors"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/template"
)

type ResendNotifier struct {
	client          *resend.Client
	fromEmail       string
	log             logger.Logger
	templateManager template.Manager
}


func NewResendNotifier(apiKey, fromEmail string, log logger.Logger, templateManager template.Manager) (person.Notifier, error) {
	if apiKey == "" {
		return nil, errs.ErrAPIKeyRequired
	}
	if fromEmail == "" {
		return nil, errs.ErrFromEmailRequired
	}

	client := resend.NewClient(apiKey)

	return &ResendNotifier{
		client:          client,
		fromEmail:       fromEmail,
		log:             log,
		templateManager: templateManager,
	}, nil
}


func (r *ResendNotifier) sendEmail(to, subject, html string) (string, error) {
	params := &resend.SendEmailRequest{
		From:    r.fromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	sent, err := r.client.Emails.Send(params)
	if err != nil {
		return "", err
	}

	return sent.Id, nil
}

func (r *ResendNotifier) SendVerificationEmail(email, verificationLink string) error {
	if email == "" {
		return errs.ErrEmailRequired
	}
	if verificationLink == "" {
		return errs.ErrVerificationLinkReq
	}

	html, err := r.templateManager.GetVerificationEmailContent(verificationLink)
	if err != nil {
		r.log.Error("generate verification email content", "email", email, "error", err)
		return fmt.Errorf("%w: %v", errs.ErrGenerateEmailContent, err)
	}

	id, err := r.sendEmail(email, SubjectVerificationEmail, html)
	if err != nil {
		r.log.Error("send verification email", "email", email, "error", err)
		return fmt.Errorf("%w: %v", errs.ErrSendEmail, err)
	}

	r.log.Info("verification email sent", "email", email, "id", id)
	return nil
}

func (r *ResendNotifier) SendPasswordRecoveryEmail(email, recoveryCode string) error {
	if email == "" {
		return errs.ErrEmailRequired
	}
	if recoveryCode == "" {
		return errs.ErrRecoveryCodeReq
	}

	html, err := r.templateManager.GetRecoveryEmailContent(recoveryCode)
	if err != nil {
		r.log.Error("generate recovery email content", "email", email, "error", err)
		return fmt.Errorf("%w: %v", errs.ErrGenerateEmailContent, err)
	}

	id, err := r.sendEmail(email, SubjectRecoveryEmail, html)
	if err != nil {
		r.log.Error("send recovery email", "email", email, "error", err)
		return fmt.Errorf("%w: %v", errs.ErrSendEmail, err)
	}

	r.log.Info("recovery email sent", "email", email, "id", id)
	return nil
}
