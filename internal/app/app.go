package app

import (
	"database/sql"
	"forum/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	h "forum/internal/controller/http"
	"forum/internal/repository"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

/*
	TODO: Add graceful shutdown
*/

func Run(cfg *config.Config) {
	// Database connection
	db, err := OpenDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Error opening database connection:%v", err)
	}

	// Service
	r := repository.New(db)
	s := service.New(r)

	// Server creation
	server := &http.Server{
		Addr:    "127.0.0.1" + cfg.Http.Addr,
		Handler: h.NewRouter(s),
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sigCh
		log.Printf("signal received:%s", sig.String())
		db.Close()

		os.Exit(0)
	}()

	// Starting the server
	log.Printf("starting the server on address:%s", cfg.Addr)
	err = server.ListenAndServe()
	log.Fatalf("Listen and serve error:%v", err)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
