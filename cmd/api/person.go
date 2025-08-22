package api

import (
	"html/template"

	domain "github.com/EstebanGitPro/motogo-backend/internal/domain/person"
)

type ResetPasswordRequest struct {
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type SendPasswordRecoveryRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ValidateRecoveryCodeRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}


type PersonLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PersonRequest struct {
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	Password            string `json:"password"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Role                string `json:"role"`
}

type PersonUpdateRequest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	PhoneNumber    string `json:"phone_number"`
}


type ValidateRecoveryCodeResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
}


type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}


type LoginResponse struct {
	ID                  string `json:"id"`
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Role                string `json:"role"`
	Token               string `json:"token"`
}

type PersonEmailVerifiedResponse struct {
	EmailVerified bool `json:"email_verified"`
}

type PersonResponse struct {
	ID                  string `json:"id"`
	IdentityNumber      string `json:"identity_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	SecondLastName      string `json:"second_last_name"`
	Email               string `json:"email"`
	PhoneNumber         string `json:"phone_number"`
	EmailVerified       bool   `json:"email_verified"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Role                string `json:"role"`
}

type ResponseEmail struct {
	Title   string
	Content template.HTML
}

func (p PersonUpdateRequest) ToDomain() domain.Person {
	return domain.Person{
		FirstName:      p.FirstName,
		LastName:       p.LastName,
		SecondLastName: &p.SecondLastName,
		PhoneNumber:    p.PhoneNumber,
	}
}

func (p PersonLogin) ToDomain() domain.Person {
	return domain.Person{
		Email:    p.Email,
		Password: p.Password,
	}
}

func (p PersonRequest) ToDomain() domain.Person {
	return domain.Person{
		IdentityNumber:      p.IdentityNumber,
		FirstName:           p.FirstName,
		LastName:            p.LastName,
		SecondLastName:      &p.SecondLastName,
		Email:               p.Email,
		PhoneNumber:         p.PhoneNumber,
		Password:            p.Password,
		EmailVerified:       p.EmailVerified,
		PhoneNumberVerified: p.PhoneNumberVerified,
		Role:                p.Role,
	}
}
