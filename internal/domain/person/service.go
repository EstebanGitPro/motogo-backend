package person

import (
	//security "github.com/andresh296/go-crud/internal/platform/token"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"log"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/internal/domain/token"
	"github.com/google/uuid"
)

type Repository interface {
	Save(person Person) error
	GetPersonByEmail(email string) (*Person, error)
	GetByID(id string) (*Person, error)
	SaveVerificationToken(token *EmailVerificationToken) error
	GetVerificationTokenByHash(hashedToken string) error
	MarkTokenAsUsedTx(tx *sql.Tx, tokenID string) error
	MarkEmailAsVerifiedTx(tx *sql.Tx, userID string) error
	CleanupExpiredTokens() error
}

type Service interface {
	GetByID(id string) (*Person, error)
	GetPersonByEmail(email string) (*Person, error)
	Save(person Person) (Person, error)
	VerifyEmailByToken(tokenString string) error
}

type Notifier interface {
	SendVerificationEmail(email string, verificationLink string) error
}

type service struct {
	repository     Repository
	notifier       Notifier
	tokenGenerator token.Generator
	config         *config.Config
}

func NewService(repo Repository, notifier Notifier, tokenGenerator token.Generator, cfg *config.Config) Service {
	return &service{
		repository:     repo,
		notifier:       notifier,
		tokenGenerator: tokenGenerator,
		config:         cfg,
	}
}

func (s service) GetByID(id string) (*Person, error) {
	return s.repository.GetByID(id)
}

func (s service) GetPersonByEmail(email string) (*Person, error) {
	return s.repository.GetPersonByEmail(email)
}

func (s service) Save(person Person) (Person, error) {

	existingPerson, err := s.repository.GetPersonByEmail(person.Email)
	if err == nil && existingPerson != nil {
		return Person{}, ErrDuplicateUser
	}

	person.setID()
	if err := person.hashPassword(); err != nil {
		return Person{}, err
	}

	err = s.repository.Save(person)
	if err != nil {
		return Person{}, ErrUserCannotSave
	}

	verificationToken, err := s.generateSecureVerificationToken(person.ID)
	if err != nil {
		return Person{}, err
	}

	verificationLink := fmt.Sprintf("%s/v1/auth/verify-email/%s",
		s.config.Verification.BaseURL,
		verificationToken.RawToken)

	err = s.notifier.SendVerificationEmail(person.Email, verificationLink)
	if err != nil {
		log.Printf("Error sending verification email to %s: %v", person.Email, err)
	}

	return person, nil
}

func (s service) generateSecureVerificationToken(userID string) (*EmailVerificationToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}

	rawToken := base64.URLEncoding.EncodeToString(tokenBytes)

	hasher := sha256.New()
	hasher.Write([]byte(rawToken))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	verificationToken := &EmailVerificationToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}

	if err := s.repository.SaveVerificationToken(verificationToken); err != nil {
		return nil, err
	}

	verificationToken.RawToken = rawToken

	return verificationToken, nil
}

func (s service) VerifyEmailByToken(tokenstring string) error {

	hasher := sha256.New()
	hasher.Write([]byte(tokenstring))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	err := s.repository.GetVerificationTokenByHash(hashedToken)
	if err != nil {
		return ErrTokenNotFound
	}

	return err
}

/*
func (s service) VerifyEmail(tokenString string) error {
	// Validar token (retorna identityNumber)
	identityNumber, err := s.tokenGenerator.Validate(tokenString)
	if err != nil {
		return err
	}

	// Iniciar transacción
	tx, err := s.repository.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Marcar el email como verificado
	if err := s.repository.MarkEmailAsVerifiedTx(tx, identityNumber); err != nil {
		return err
	}

	// Confirmar transacción
	return tx.Commit()
}
*/

func (s service) CleanupExpiredTokens() error {
	return s.repository.CleanupExpiredTokens()
}

func (s service) StartCleanupScheduler() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			if err := s.CleanupExpiredTokens(); err != nil {
				log.Printf("Error limpiando tokens expirados: %v", err)
			}
		}
	}()
}
