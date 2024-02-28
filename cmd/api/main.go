package main

import (
	"database/sql"
	"forum/config"
	"forum/pkg/sesm"
	"forum/pkg/sesm/sqlite3store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"forum/internal/handlers"
	"forum/internal/repository"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Parse config
	cfg := config.NewConfig()

	// Logger init
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Database connection
	db, err := OpenDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Error opening database connection:%v", err)
	}

	// Repos and Services init
	r := repository.New(db)
	s := service.New(r)

	// Session Manager creation
	sesm := sesm.New()
	sesm.Store = sqlite3store.New(db)

	// Routes init
	routes := handlers.NewRouter(
		s,
		sesm,
		logger,
	)

	// Server creation
	server := &http.Server{
		Addr:    "0.0.0.0" + cfg.Http.Addr,
		Handler: routes.Register(),
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sigCh
		logger.Printf("signal received:%s", sig.String())
		db.Close()

		os.Exit(0)
	}()

	// Starting the server
	logger.Printf("starting the server on address - http://localhost%s", cfg.Addr)
	err = server.ListenAndServe()
	logger.Fatalf("Listen and serve error:%v", err)
}

// OpenDB opens connection to the database using standard sql library
// with given Data Source Name (DSN)
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

	// Enable foreign keys (they are disabled by default for backwards compatibility)
	query := "PRAGMA foreign_keys = ON;"
	_, err = db.Exec(query)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
