package person

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/internal/domain/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func (s service) generateSecureVerificationToken(userID, tokenType string) (*UserToken, error) {
	var rawToken string
	var code string
	var hashedToken string

	if tokenType == TokenTypePasswordRecovery {
		generatedCode, err := generateSecureCode(6)
		if err != nil {
			return nil, err
		}
		rawToken = generatedCode

		codeHasher := sha256.New()
		codeHasher.Write([]byte(generatedCode))
		hashedCode := hex.EncodeToString(codeHasher.Sum(nil))
		code = hashedCode

		hashedToken = "password_recovery_placeholder_token_not_used_for_verification"

	} else {
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			return nil, err
		}
		rawToken = base64.URLEncoding.EncodeToString(tokenBytes)

		hasher := sha256.New()
		hasher.Write([]byte(rawToken))
		hashedToken = hex.EncodeToString(hasher.Sum(nil))
	}

	var expiration time.Duration
	if tokenType == TokenTypePasswordRecovery {
		expiration = 15 * time.Minute
	} else {
		expiration = 10 * time.Minute
	}

	verificationToken := &UserToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     hashedToken,
		Code:      code,
		Type:      tokenType,
		ExpiresAt: time.Now().UTC().Add(expiration),
		Used:      false,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repository.SaveVerificationToken(verificationToken); err != nil {
		return nil, err
	}

	verificationToken.RawToken = rawToken

	return verificationToken, nil
}

func (p Person) comparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	if err != nil {
		return ErrValidationUser
	}
	return err
}

func (s service) CleanupExpiredTokens() error {
	return s.repository.CleanupExpiredTokens()
}

func (s service) GetByID(id string) (*Person, error) {
	return s.repository.GetByID(id)
}

func (s service) GetPersonByEmail(email string) (*Person, error) {
	return s.repository.GetPersonByEmail(email)
}

func generateSecureCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
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
		return Person{}, err
	}

	verificationToken, err := s.generateSecureVerificationToken(person.ID, TokenTypeEmailVerification)
	if err != nil {
		return Person{}, err
	}

	verificationLink := fmt.Sprintf("%s/v1/motogo/auth/verify-email/%s",
		s.config.Verification.BaseURL,
		verificationToken.RawToken)

	err = s.notifier.SendVerificationEmail(person.Email, verificationLink)
	if err != nil {

		log.Printf("Error sending verification email to %s: %v", person.Email, err)
	}

	return person, nil
}

func (s service) VerifyEmailByToken(tokenString string) error {
	hasher := sha256.New()
	hasher.Write([]byte(tokenString))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	err := s.repository.GetVerificationTokenByHash(hashedToken)
	if err != nil {
		return err
	}

	return err
}

func (s *service) CheckPasswordRecoveryByCode(codeString string) error {
	hasher := sha256.New()
	hasher.Write([]byte(codeString))
	hashedCode := hex.EncodeToString(hasher.Sum(nil))

	_, err := s.repository.ConsumePasswordRecoveryToken(hashedCode)
	return err
}

func (s *service) Login(person Person) (*Person, string, error) {
	personFound, err := s.repository.GetPersonByEmail(person.Email)
	if err != nil {
		return nil, "", ErrValidationUser
	}

	err = personFound.comparePassword(person.Password)
	if err != nil {
		return nil, "", ErrValidationUser
	}

	if !personFound.EmailVerified {
		verificationToken, err := s.generateSecureVerificationToken(personFound.ID, TokenTypeEmailVerification)
		if err != nil {
			return nil, "", err
		}

		verificationLink := fmt.Sprintf("%s/v1/motogo/auth/verify-email/%s",
			s.config.Verification.BaseURL,
			verificationToken.RawToken)

		err = s.notifier.SendVerificationEmail(personFound.Email, verificationLink)
		if err != nil {
			log.Printf("Error sending verification email to %s: %v", personFound.Email, err)
		}
		return nil, "", ErrorEmailNotVerified
	}

	token, err := s.tokenGenerator.Generate(personFound.ID, 15*time.Minute)
	if err != nil {
		return nil, "", err
	}

	return personFound, token, nil
}

func (s service) StartCleanupScheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			if err := s.CleanupExpiredTokens(); err != nil {
				log.Printf("Error limpiando tokens expirados: %v", err)
			}
		}
	}()
}

func (s *service) Update(id string, person Person) error {
	existingPerson, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}
	if existingPerson == nil {
		return ErrUserCannotFound
	}

	existingPerson.FirstName = person.FirstName
	existingPerson.LastName = person.LastName
	existingPerson.SecondLastName = person.SecondLastName
	existingPerson.PhoneNumber = person.PhoneNumber

	return s.repository.Update(id, *existingPerson)
}

func (s service) SendPasswordRecoveryEmail(email string) error {
	personFound, err := s.repository.GetPersonByEmail(email)
	if err != nil {
		return err
	}
	if personFound == nil {
		return ErrUserCannotFound
	}

	verificationToken, err := s.generateSecureVerificationToken(personFound.ID, TokenTypePasswordRecovery)
	if err != nil {
		return err
	}

	err = s.notifier.SendPasswordRecoveryEmail(email, verificationToken.RawToken)
	if err != nil {
		log.Printf("Error sending password recovery email to %s: %v", email, err)
		return err
	}

	return nil
}

func (s *service) RecoveryPassword(userID string, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repository.UpdatePassword(userID, string(hashedPassword))
}

func (s *service) GetUserIDFromRecoveryCode(codeString string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(codeString))
	hashedCode := hex.EncodeToString(hasher.Sum(nil))

	token, err := s.repository.GetTokenByHash(hashedCode, TokenTypePasswordRecovery)
	if err != nil {
		return "", err
	}

	return token.UserID, nil
}
