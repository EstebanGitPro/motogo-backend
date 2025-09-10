package notification

const (

	SubjectVerificationEmail = "Verify your email - MotoGo"
	SubjectRecoveryEmail     = "Password recovery - MotoGo"
	SubjectWelcomeEmail      = "Welcome to MotoGo!"

	
	TemplateVerification = "templates/verification_email.html"
	TemplateRecovery     = "templates/recovery_email.html"


	LinkExpirationMinutes = 15
	CodeExpirationMinutes = 15
)
