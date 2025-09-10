package person

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/internal/domain/token"
	"golang.org/x/crypto/bcrypt"
)

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
