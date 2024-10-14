package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var (
	testDB *sql.DB
)

const (
	postgresDNS = "postgres://root:secret@localhost:5432/blog_test?sslmode=disable"

	test
)

func TestMain(m *testing.M) {
	dbConn, err := sql.Open("postgres", postgresDNS)
	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()

	if err != nil {
		log.Fatal(err)
	}

	testDB = dbConn

	code := m.Run()

	err = testDB.Close()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)

}
