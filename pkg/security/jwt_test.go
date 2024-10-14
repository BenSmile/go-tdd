package security

import (
	"github.com/bensmile/go-api-tdd/pkg/common"
	"github.com/bensmile/go-api-tdd/pkg/domain"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	_, err := NewJWT("")
	if err == nil {
		t.Error("NewJWT should return an error when key is too short")
	}
}

func TestJWTToken(t *testing.T) {
	key := common.RandomString(32)
	newJWT, err := NewJWT(key)

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	user := domain.User{
		Id:    1,
		Email: "test@test.com",
	}

	payload, err := newJWT.CreateToken(user, 1*time.Minute)

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if payload.UserId != user.Id {
		t.Errorf("expected userId to be %d, got %d", user.Id, payload.UserId)
	}

	if len(payload.Token) == 0 {
		t.Error("expected token to be set")
	}

	if payload.ExpiresAt.Before(time.Now()) {
		t.Error("token should return a not expired token")
	}

	_, err = newJWT.VerifyToken(payload.Token + "invalid")

	if err == nil {
		t.Error("VerifyToken should return error for invalid token")
	}

	expiredToken, err := newJWT.CreateToken(user, -2*time.Minute)

	if err != nil {
		t.Errorf("CreateToken should return error, got %v", err)
	}

	_, err = newJWT.VerifyToken(expiredToken.Token)

	if err == nil {
		t.Error("VerifyToken should return error for expired token")
	}

}
