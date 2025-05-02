package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	IP    string `json:"ip"`
	jwt.RegisteredClaims
}

func NewUserClaim(id uint, email string, ip string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &UserClaims{
		ID:    id,
		Email: email,
		IP:    ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID.String(),
			Subject:   email,
		},
	}, nil
}
