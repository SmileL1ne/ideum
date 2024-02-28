package main

import (
	"crypto/tls"
	"database/sql"
	"forum/config"
	"forum/pkg/sesm"
	"forum/pkg/sesm/sqlite3store"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum/internal/repository"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

type routes struct {
	service   *service.Services
	tempCache map[string]*template.Template
	sesm      *sesm.SessionManager
	logger    *log.Logger
}

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

	// Temporary cache for one-time template initialization and subsequent
	// storage in templates map
	tempCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("Error creating cached templates:%v", err)
	}

	// Repos and Services init
	r := repository.New(db)
	s := service.New(r)

	// Session Manager creation
	sesm := sesm.New()
	sesm.Store = sqlite3store.New(db)

	// Routes init
	routes := &routes{
		service:   s,
		tempCache: tempCache,
		sesm:      sesm,
		logger:    logger,
	}

	// tls config
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}

	// Server creation
	server := &http.Server{
		Addr:           "0.0.0.0" + cfg.Http.Addr,
		Handler:        routes.newRouter(),
		MaxHeaderBytes: 1 << 20, // 1 mb
		TLSConfig:      tlsConfig,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sigCh
		routes.logger.Printf("signal received:%s", sig.String())
		db.Close()

		os.Exit(0)
	}()

	// Starting the server
	routes.logger.Printf("starting the server on address - https://localhost%s", cfg.Addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	routes.logger.Fatalf("Listen and serve error:%v", err)
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
