package token
import (
	"time"
)
type Generator interface {
	Generate(userID string, duration time.Duration) (string, error)
	Validate(tokenString string) (userID string, err error)
}