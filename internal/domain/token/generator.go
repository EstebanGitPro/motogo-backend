package token

import "time"

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type Generator interface {
	Generate(userID string, duration time.Duration) (string, error)
	Validate(tokenString string) (*Claims, error)
}
