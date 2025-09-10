package template

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// Manager define la interfaz para el manejo de templates
type Manager interface {
	GetVerificationEmailContent(verificationLink string) (string, error)
	GetRecoveryEmailContent(recoveryCode string) (string, error)
}

// TemplateManager implementa Manager
type TemplateManager struct {
	verificationTemplate *template.Template
	recoveryTemplate     *template.Template
}

// TemplateData contiene los datos para los templates
type TemplateData struct {
	VerificationLink string
	RecoveryCode     string
	UserName         string
}

// NewTemplateManager crea un nuevo TemplateManager cargando todos los templates al inicio
func NewTemplateManager() (*TemplateManager, error) {
	// Cargar template de verificación
	verificationTmpl, err := template.ParseFS(templateFS, "templates/verification_email.html")
	if err != nil {
		return nil, fmt.Errorf("error al cargar template de verificación: %w", err)
	}

	// Cargar template de recuperación
	recoveryTmpl, err := template.ParseFS(templateFS, "templates/recovery_email.html")
	if err != nil {
		return nil, fmt.Errorf("error al cargar template de recuperación: %w", err)
	}

	return &TemplateManager{
		verificationTemplate: verificationTmpl,
		recoveryTemplate:     recoveryTmpl,
	}, nil
}

// GetVerificationEmailContent genera el contenido del email de verificación
func (tm *TemplateManager) GetVerificationEmailContent(verificationLink string) (string, error) {
	data := TemplateData{
		VerificationLink: verificationLink,
	}

	var buf bytes.Buffer
	if err := tm.verificationTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error al ejecutar template de verificación: %w", err)
	}

	return buf.String(), nil
}

// GetRecoveryEmailContent genera el contenido del email de recuperación
func (tm *TemplateManager) GetRecoveryEmailContent(recoveryCode string) (string, error) {
	data := TemplateData{
		RecoveryCode: recoveryCode,
	}

	var buf bytes.Buffer
	if err := tm.recoveryTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error al ejecutar template de recuperación: %w", err)
	}

	return buf.String(), nil
}
