package api

import (
	"html/template"

	domain "github.com/EstebanGitPro/motogo-backend/internal/domain/person"
)

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

type PersonUpdateRequest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	PhoneNumber    string `json:"phone_number"`
}

type ResponseEmail struct {
	Title   string
	Content template.HTML
}

func (P PersonUpdateRequest) ToDomain() domain.Person {
	return domain.Person{
		FirstName:      P.FirstName,
		LastName:       P.LastName,
		SecondLastName: P.SecondLastName,
		PhoneNumber:    P.PhoneNumber,
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
		SecondLastName:      p.SecondLastName,
		Email:               p.Email,
		PhoneNumber:         p.PhoneNumber,
		Password:            p.Password,
		EmailVerified:       p.EmailVerified,
		PhoneNumberVerified: p.PhoneNumberVerified,
		Role:                p.Role,
	}
}
