package main

import (
	"database/sql"
	"github.com/bensmile/go-api-tdd/pkg/store/sqlstore/postgres"
	_ "github.com/lib/pq"
	"log"
)

const (
	postgresDNS = "postgres://root:secret@localhost:5432/blog?sslmode=disable"
	driver      = "postgres"
	key         = "a6c0a92eceb23257a7fe647f35616ce4e6d2720da278ec0e9fd0433b9a8e21c2"
)

func main() {

	srv, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	_, err = connectToDB(driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = srv.run(":8080"); err != nil {
		log.Fatal(err)
	}
	log.Println("Hello World")
}

func setup() (*server, error) {
	db, err := connectToDB(driver)
	if err != nil {
		return nil, err
	}
	pStore := postgres.NewPostgresStore(db)
	srv := newServer(pStore, nil)
	srv.routes()
	return srv, nil
}

func connectToDB(driver string) (*sql.DB, error) {

	db, err := sql.Open(driver, postgresDNS)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, nil
}
