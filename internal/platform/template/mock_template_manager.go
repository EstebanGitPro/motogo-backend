package template

type MockTemplateManager struct {
	VerificationEmailContent string
	RecoveryEmailContent     string
	VerificationError        error
	RecoveryError            error
}

func NewMockTemplateManager() *MockTemplateManager {
	return &MockTemplateManager{
		VerificationEmailContent: "<html><body>Verification email content</body></html>",
		RecoveryEmailContent:     "<html><body>Recovery email content</body></html>",
	}
}

func (m *MockTemplateManager) GetVerificationEmailContent(verificationLink string) (string, error) {
	if m.VerificationError != nil {
		return "", m.VerificationError
	}
	return m.VerificationEmailContent, nil
}

func (m *MockTemplateManager) GetRecoveryEmailContent(recoveryCode string) (string, error) {
	if m.RecoveryError != nil {
		return "", m.RecoveryError
	}
	return m.RecoveryEmailContent, nil
}
