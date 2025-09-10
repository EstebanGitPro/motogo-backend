package person

import "database/sql"

type Repository interface {
	Save(person Person) error
	GetPersonByEmail(email string) (*Person, error)
	GetByID(id string) (*Person, error)
	SaveVerificationToken(token *UserToken) error
	GetVerificationTokenByHash(hashedToken string) error
	MarkTokenAsUsedTx(tx *sql.Tx, tokenID string) error
	MarkEmailAsVerifiedTx(tx *sql.Tx, userID string) error
	CleanupExpiredTokens() error
	Update(id string, person Person) error
	UpdatePassword(userID, hashedPassword string) error
	GetTokenByHash(hashedToken, tokenType string) (*UserToken, error)
	ConsumePasswordRecoveryToken(codeString string) (string, error)
}

type Service interface {
	GetByID(id string) (*Person, error)
	GetPersonByEmail(email string) (*Person, error)
	Save(person Person) (Person, error)
	VerifyEmailByToken(tokenString string) error
	CheckPasswordRecoveryByCode(codeString string) error
	CleanupExpiredTokens() error
	StartCleanupScheduler()
	Login(person Person) (*Person, string, error)
	Update(id string, person Person) error
	SendPasswordRecoveryEmail(email string) error
	RecoveryPassword(userID, newPassword string) error
	GetUserIDFromRecoveryCode(codeString string) (string, error)
}

type Notifier interface {
	SendVerificationEmail(email string, verificationLink string) error
	SendPasswordRecoveryEmail(email string, recoveryCode string) error
}