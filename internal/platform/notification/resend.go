
package notification

import (
	"errors"
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
)

type ResendNotifier struct {
	client    *resend.Client
	fromEmail string
}

func NewResendNotifier(apiKey, fromEmail string) (*ResendNotifier, error) {
	if apiKey == "" {
		return nil, errors.New("la clave de API de Resend es requerida")
	}
	if fromEmail == "" {
		return nil, errors.New("el correo del remitente es requerido")
	}

	client := resend.NewClient(apiKey)

	return &ResendNotifier{
		client:    client,
		fromEmail: fromEmail,
	}, nil
}

func (r *ResendNotifier) SendVerificationEmail(email, verificationLink string) error {
	
	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">¡Bienvenido a MotoGo!</h2>
			<p>Gracias por registrarte en nuestra plataforma.</p>
			<p>Para completar tu registro, por favor haz clic en el siguiente enlace:</p>
			<div style="margin: 20px 0;">
				<a href="%s" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">
					Verificar mi email
				</a>
			</div>
			<p style="color: #666; font-size: 14px;">
				Si no puedes hacer clic en el botón, copia y pega este enlace en tu navegador:
			</p>
			<p style="color: #666; font-size: 14px; word-break: break-all;">
				%s
			</p>
			<p style="color: #999; font-size: 12px; margin-top: 30px;">
				Este enlace expirará en 15 minutos por seguridad.
			</p>
		</div>
	`, verificationLink, verificationLink)

	params := &resend.SendEmailRequest{
		From:    r.fromEmail,
		To:      []string{email},
		Subject: "Verifica tu email - MotoGo",
		Html:    htmlContent,
	}

	sent, err := r.client.Emails.Send(params)
	if err != nil {
		log.Printf("Error al enviar correo a %s: %v", email, err)
		return fmt.Errorf("no se pudo enviar el correo de verificación: %w", err)
	}

	log.Printf("Correo de verificación enviado exitosamente a %s (ID: %s)", email, sent.Id)
	return nil
}

