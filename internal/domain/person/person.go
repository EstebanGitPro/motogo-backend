package person

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	ID                  string `json:"id"`
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Password            string `json:"-"`
	Role                string `json:"role"`
}

func (u *Person) setID() {
	u.ID = uuid.New().String()
}

func (u *Person) hashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

const (
	TokenTypeEmailVerification = "email_verification"
	TokenTypePasswordRecovery  = "password_recovery"
)

type UserToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	Code      string    `json:"code" db:"code"`
	Type      string    `json:"type" db:"type"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	RawToken  string    `json:"-"`
}
