package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bensmile/go-api-tdd/pkg/domain"
	"github.com/bensmile/go-api-tdd/pkg/security"
)

func TestCreateUser(t *testing.T) {

	testCases := []struct {
		name         string
		expectedCode int
		body         string
	}{
		{
			name:         "OK",
			expectedCode: http.StatusCreated,
			body: `{
			"name" : "Smile",
			"email" : "smile1@gmail.com",
			"password" : "password"
		}`,
		}, {
			name:         "Bad Json",
			expectedCode: http.StatusBadRequest,
			body:         `{"name" : }`,
		},
		{
			name:         "Validation Error",
			expectedCode: http.StatusBadRequest,
			body: `{
			"name" : "",
			"email" : "",
			"password" : ""
		}`,
		},
	}

	newJWT, err := security.NewJWT(key)

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	srv := newServer(testStore, newJWT)
	srv.store.DeleteAllUsers()
	ts := newTestServer(srv.routes())

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			res, err := ts.Client().
				Post(ts.URL+"/api/v1/users",
					"application/json",
					strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.expectedCode {
				t.Errorf("want status code %d, got %d", tc.expectedCode, res.StatusCode)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	newJWT, err := security.NewJWT(key)

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	srv := newServer(testStore, newJWT)
	tstSrv := newTestServer(srv.routes())

	if err := srv.store.DeleteAllPosts(); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if err := srv.store.DeleteAllUsers(); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	_, err = srv.store.CreateUser(&domain.User{
		Email:    "smile1@gmail.com",
		Password: "password",
		Name:     "Smile",
	})

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	testCases := []struct {
		name         string
		expectedCode int
		body         string
		checkBody    func(t *testing.T, body []byte)
	}{
		{
			name:         "OK",
			expectedCode: http.StatusOK,
			body: `{
			"email" : "smile1@gmail.com",
			"password" : "password"
		}`, checkBody: func(t *testing.T, body []byte) {
				loginResponse := struct {
					Token string `json:"token"`
				}{}

				if err := json.Unmarshal(body, &loginResponse); err != nil {
					t.Fatal(err)
				}
				if loginResponse.Token == "" {
					t.Error("want token, got empty string")
				}
			},
		},
		{
			name:         "Invalid Json",
			expectedCode: http.StatusBadRequest,
			body: `{Invalid}"
		}`,
		}, {
			name:         "Missing username and password",
			expectedCode: http.StatusBadRequest,
			body: `{
			"email" : "",
			"password" : ""
		}`}, {
			name:         "User not found",
			expectedCode: http.StatusUnauthorized,
			body: `{
			"email" : "smile1@gmail",
			"password" : "password"
		}`}, {
			name:         "Invalid credentials",
			expectedCode: http.StatusUnauthorized,
			body: `{
			"email" : "smile1@gmail.com",
			"password" : "invalid password"
		}`},
	}

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			res, err := tstSrv.Client().
				Post(tstSrv.URL+"/api/v1/users/login",
					"application/json",
					strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.expectedCode {
				t.Errorf("want status code %d, got %d", tc.expectedCode, res.StatusCode)
			}

			if tc.checkBody != nil {
				defer res.Body.Close()

				bodyByteSlice, err := io.ReadAll(res.Body)

				if err != nil {
					t.Fatal(err)
				}
				tc.checkBody(t, bodyByteSlice)
			}
		})
	}
}
