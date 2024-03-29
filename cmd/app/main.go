package main

import (
	"crypto/tls"
	"forum/config"
	"forum/pkg/database/sqlite3"
	"forum/pkg/sesm"
	"forum/pkg/sesm/sqlite3store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum/internal/handlers"
	"forum/internal/repository"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.Load()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := sqlite3.OpenDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Error opening database connection:%v", err)
	}

	r := repository.New(db)
	s := service.New(r)

	sesm := sesm.New()
	sesm.Store = sqlite3store.New(db)

	routes := handlers.NewRouter(
		s,
		sesm,
		logger,
		cfg,
	)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}

	tlsCert := "./tls/cert.pem"
	tlsKey := "./tls/key.pem"

	server := &http.Server{
		Addr:           "0.0.0.0" + cfg.Http.Addr,
		Handler:        routes.Register(),
		MaxHeaderBytes: 1 << 20, // 1 mb
		TLSConfig:      tlsConfig,
		IdleTimeout:    time.Minute,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   35 * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sigCh
		logger.Printf("signal received:%s", sig.String())
		db.Close()

		os.Exit(0)
	}()

	logger.Printf("starting the server on address - https://localhost%s", cfg.Addr)
	err = server.ListenAndServeTLS(tlsCert, tlsKey)
	logger.Fatalf("Listen and serve error:%v", err)
}
