package jwt

import (
	"errors"
	"time"

	"github.com/EstebanGitPro/motogo-backend/internal/domain/token"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type jwtGenerator struct {
	secretKey []byte
}

func New(secret string) token.Generator {
	return &jwtGenerator{
		secretKey: []byte(secret),
	}
}

func (j *jwtGenerator) Generate(userID string, duration time.Duration) (string, error) {
	claims := &Claims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwtGenerator) Validate(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.ID, nil
	}

	return "", jwt.ErrTokenInvalidClaims
}
