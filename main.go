package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:./internal/database/rabbit.db?_foreign_keys=on")
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	setup, err := os.ReadFile("./migrations/002_add_indexes_up.sql")
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	_, err = db.Exec(string(setup))
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
}
