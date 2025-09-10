package person

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

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