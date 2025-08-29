package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID       uuid.UUID `json:"uid"`
	Roles        []string  `json:"roles,omitempty"`
	Perms        []string  `json:"perms,omitempty"`
	TokenVersion int       `json:"tokenVersion,omitempty"`
	jwt.RegisteredClaims
}

func SignHS256(secret []byte, issuer string, sub uuid.UUID, roles, perms []string, tokenVersion int, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:       sub,
		Roles:        roles,
		Perms:        perms,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   sub.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(secret)
}

func ParseAndVerify(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
