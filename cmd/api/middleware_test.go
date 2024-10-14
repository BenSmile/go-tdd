package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bensmile/go-api-tdd/pkg/domain"
	"github.com/bensmile/go-api-tdd/pkg/security"
	"github.com/gin-gonic/gin"
)

func TestApplyAuthenticationMiddleware(t *testing.T) {

	newJWT, err := security.NewJWT(key)

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	srv := newServer(testStore, newJWT)
	if err := srv.store.DeleteAllPosts(); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	if err := srv.store.DeleteAllUsers(); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	ts := newTestServer(srv.routes())

	user, _ := srv.store.CreateUser(&domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "OJ",
	})

	testCases := []struct {
		name         string
		expectedCode int
		setupHeader  func(r *http.Request)
		checkBody    func(t *testing.T, res *http.Response)
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			setupHeader: func(r *http.Request) {
				jwtPayload, _ := srv.jwt.CreateToken(*user, 2*time.Minute)
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtPayload.Token))
			},
		},
		{
			name:         "Missing auth header",
			expectedCode: http.StatusUnauthorized,
			setupHeader:  func(r *http.Request) {},
		}, {
			name:         "Invalid token type",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "invalid token type")
			},
		}, {
			name:         "Expired token",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				jwtPayload, _ := srv.jwt.CreateToken(*user, -2*time.Minute)
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtPayload.Token))
			},
		},
		{
			name:         "Invalid token",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid_token")
			},
		}, {
			name:         "Invalid token | No existing User",
			expectedCode: http.StatusUnauthorized,
			setupHeader: func(r *http.Request) {
				jwtPayload, _ := srv.jwt.CreateToken(domain.User{
					Id: -1,
				}, 1*time.Minute)
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtPayload.Token))
			},
		},
	}

	srv.router.GET("/auth", srv.applyAuthentication(), func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ts.URL+"/auth", nil)
			req.RequestURI = ""
			tc.setupHeader(req)

			res, err := ts.Client().Do(req)

			if err != nil {
				t.Fatal(err)
			}

			if tc.expectedCode != res.StatusCode {
				t.Errorf("want %d; got %d", tc.expectedCode, res.StatusCode)
			}

			if tc.checkBody != nil {
				tc.checkBody(t, res)
			}
		})
	}
}
