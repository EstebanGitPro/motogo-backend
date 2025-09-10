package notification

import (
	"testing"

	errs "github.com/EstebanGitPro/motogo-backend/internal/domain/errors"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/template"
)

func TestNewResendNotifier(t *testing.T) {
	mockLogger := logger.NewMockLogger()
	mockTemplateManager := template.NewMockTemplateManager()

	tests := []struct {
		name      string
		apiKey    string
		fromEmail string
		wantErr   error
	}{
		{
			name:      "valid parameters",
			apiKey:    "test-key",
			fromEmail: "test@example.com",
			wantErr:   nil,
		},
		{
			name:      "missing api key",
			apiKey:    "",
			fromEmail: "test@example.com",
			wantErr:   errs.ErrAPIKeyRequired,
		},
		{
			name:      "missing from email",
			apiKey:    "test-key",
			fromEmail: "",
			wantErr:   errs.ErrFromEmailRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewResendNotifier(tt.apiKey, tt.fromEmail, mockLogger, mockTemplateManager)
			if err != tt.wantErr {
				t.Errorf("NewResendNotifier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResendNotifier_SendVerificationEmail_Validation(t *testing.T) {
	mockLogger := logger.NewMockLogger()
	mockTemplateManager := template.NewMockTemplateManager()
	notifier, _ := NewResendNotifier("test-key", "test@example.com", mockLogger, mockTemplateManager)

	tests := []struct {
		name             string
		email            string
		verificationLink string
		wantErr          error
	}{
		{
			name:             "missing email",
			email:            "",
			verificationLink: "http://example.com/verify",
			wantErr:          errs.ErrEmailRequired,
		},
		{
			name:             "missing verification link",
			email:            "user@example.com",
			verificationLink: "",
			wantErr:          errs.ErrVerificationLinkReq,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.SendVerificationEmail(tt.email, tt.verificationLink)
			if err != tt.wantErr {
				t.Errorf("SendVerificationEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResendNotifier_SendPasswordRecoveryEmail_Validation(t *testing.T) {
	mockLogger := logger.NewMockLogger()
	mockTemplateManager := template.NewMockTemplateManager()
	notifier, _ := NewResendNotifier("test-key", "test@example.com", mockLogger, mockTemplateManager)

	tests := []struct {
		name         string
		email        string
		recoveryCode string
		wantErr      error
	}{
		{
			name:         "missing email",
			email:        "",
			recoveryCode: "123456",
			wantErr:      errs.ErrEmailRequired,
		},
		{
			name:         "missing recovery code",
			email:        "user@example.com",
			recoveryCode: "",
			wantErr:      errs.ErrRecoveryCodeReq,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.SendPasswordRecoveryEmail(tt.email, tt.recoveryCode)
			if err != tt.wantErr {
				t.Errorf("SendPasswordRecoveryEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockLogger_Functionality(t *testing.T) {
	mockLogger := logger.NewMockLogger()

	// Test Info logging
	mockLogger.Info("test info", "key", "value")
	if len(mockLogger.InfoCalls) != 1 {
		t.Errorf("Expected 1 info call, got %d", len(mockLogger.InfoCalls))
	}

	lastInfo := mockLogger.GetLastInfo()
	if lastInfo == nil || lastInfo.Message != "test info" {
		t.Errorf("Expected info message 'test info', got %v", lastInfo)
	}

	// Test Error logging
	mockLogger.Error("test error", "error", "something went wrong")
	if len(mockLogger.ErrorCalls) != 1 {
		t.Errorf("Expected 1 error call, got %d", len(mockLogger.ErrorCalls))
	}

	lastError := mockLogger.GetLastError()
	if lastError == nil || lastError.Message != "test error" {
		t.Errorf("Expected error message 'test error', got %v", lastError)
	}

	// Test Reset functionality
	mockLogger.Reset()
	if len(mockLogger.InfoCalls) != 0 || len(mockLogger.ErrorCalls) != 0 {
		t.Errorf("Expected empty calls after reset, got info: %d, error: %d", 
			len(mockLogger.InfoCalls), len(mockLogger.ErrorCalls))
	}
}
