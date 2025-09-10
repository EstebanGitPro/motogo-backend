package domain

import "errors"

var (
	ErrAPIKeyRequired       = errors.New("api key is required")
	ErrFromEmailRequired    = errors.New("from email is required")
	ErrEmailRequired        = errors.New("email is required")
	ErrVerificationLinkReq  = errors.New("verification link is required")
	ErrRecoveryCodeReq      = errors.New("recovery code is required")
	ErrGenerateEmailContent = errors.New("failed to generate email content")
	ErrSendEmail           = errors.New("failed to send email")
)
