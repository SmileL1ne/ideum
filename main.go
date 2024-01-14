package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:./internal/database/rabbit.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}

	setup, err := os.ReadFile("./migrations/001_initial_setup_up.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(setup))
	if err != nil {
		log.Fatal(err)
	}
}
