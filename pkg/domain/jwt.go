package domain

import "time"

type JWTPayload struct {
	UserId    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
	Token     string    `json:"token"`
}

func (p *JWTPayload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return ErrExpiredToken
	}

	if p.UserId == 0 {
		return ErrUserNotFound
	}

	return nil
}
