package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewToken(uid, role string) *claims {
	return &claims{
		Id:   uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Musalabel",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
}

func (c *claims) Generate() (string, error) {
	screts := os.Getenv("JWT_KEYS")
	if c == nil || c.Id == "" || c.Role == "" {
		return "", errors.New("claim struct must not be nil or empty")
	}

	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tokens.SignedString([]byte(screts))
}

func VerifyToken(token string) (*claims, error) {
	screts := os.Getenv("JWT_KEYS")
	data, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(screts), nil
	})

	if err != nil {
		return nil, err
	}

	claimData := data.Claims.(*claims)
	return claimData, nil
}
