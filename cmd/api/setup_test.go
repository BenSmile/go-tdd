package main

import (
	"github.com/bensmile/go-api-tdd/pkg/domain"
	"github.com/bensmile/go-api-tdd/pkg/store/sqlstore/postgres"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	testStore domain.Store
)

type testServer struct {
	*httptest.Server
}

func newTestServer(h http.Handler) *testServer {
	return &testServer{httptest.NewServer(h)}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	db, err := connectToDB(driver)
	if err != nil {
		log.Fatal(err)
	}

	testStore = postgres.NewPostgresStore(db)

	code := m.Run()

	_ = testStore.DeleteAllUsers()

	os.Exit(code)
}
