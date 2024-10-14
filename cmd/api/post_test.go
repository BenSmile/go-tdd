package main

import (
	"encoding/json"
	"fmt"
	"github.com/bensmile/go-api-tdd/pkg/domain"
	"github.com/bensmile/go-api-tdd/pkg/security"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {

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

	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}
	createdUser, err := srv.store.CreateUser(user)

	if err != nil {
		t.Fatal(err)
	}

	jwtPayload, err := srv.jwt.CreateToken(*createdUser, 2*time.Minute)

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	testCases := []struct {
		name         string
		body         string
		expectedCode int
		setupHeader  func(r *http.Request)
		checkBody    func(t *testing.T, res *http.Response)
	}{
		{
			name:         "Valid",
			expectedCode: http.StatusCreated,
			body:         `{"title" : "test blog", "body" : "test content"}`,
		},
		{
			name:         "Invalid request body",
			expectedCode: http.StatusBadRequest,
			body:         `{"title" : invalid}`,
		},
		{
			name:         "Validation Error",
			expectedCode: http.StatusBadRequest,
			body:         `{"title" : "", "body" : ""}`,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/blog", strings.NewReader(tc.body))
			req.Header.Set(AuthorizationHeader, fmt.Sprintf("%s %s", AuthorizationType, jwtPayload.Token))
			req.RequestURI = ""

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

func TestGetUserPosts(t *testing.T) {

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

	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}
	createdUser, err := srv.store.CreateUser(user)

	if err != nil {
		t.Fatal(err)
	}

	jwtPayload, err := srv.jwt.CreateToken(*createdUser, 2*time.Minute)

	if err != nil {
		t.Errorf("expected error to be nil when creating token, got %v", err)
	}

	_, err = srv.store.CreatePost(&domain.Post{
		UserId: createdUser.Id,
		Title:  "test post",
		Body:   "test content",
	})

	if err != nil {
		t.Errorf("expected error to be nil when creating post, got %v", err)
	}

	testCases := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "Valid",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ts.URL+"/api/v1/blog", nil)
			req.Header.Set(AuthorizationHeader, fmt.Sprintf("%s %s", AuthorizationType, jwtPayload.Token))
			req.RequestURI = ""

			res, err := ts.Client().Do(req)

			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			if tc.expectedCode != res.StatusCode {
				t.Errorf("want %d; got %d", tc.expectedCode, res.StatusCode)
			}

			var postsRes []domain.Post

			resBodyBytes, err := io.ReadAll(res.Body)

			if err != nil {
				t.Fatal(err)
			}

			if err := json.Unmarshal(resBodyBytes, &postsRes); err != nil {
				t.Fatal(err)
			}

			if len(postsRes) != 1 {
				t.Errorf("expected 1 post; got %d", len(postsRes))
			}

		})
	}

}
