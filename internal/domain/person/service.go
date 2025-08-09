package person

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
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
	SaveVerificationToken(token *EmailVerificationToken) error
	GetVerificationTokenByHash(hashedToken string) error
	MarkTokenAsUsedTx(tx *sql.Tx, tokenID string) error
	MarkEmailAsVerifiedTx(tx *sql.Tx, userID string) error
	CleanupExpiredTokens() error
	Update(id string, person Person) error
}

type Service interface {
	GetByID(id string) (*Person, error)
	GetPersonByEmail(email string) (*Person, error)
	Save(person Person) (Person, error)
	VerifyEmailByToken(tokenString string) error
	CleanupExpiredTokens() error
	StartCleanupScheduler()
	Login(person Person) (*Person, string, error)
	Update(id string, person Person) error
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
		ExpiresAt: time.Now().Add(1 * time.Minute),
		Used:      false,
		CreatedAt: time.Now(),
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
		verificationToken, err := s.generateSecureVerificationToken(personFound.ID)
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
