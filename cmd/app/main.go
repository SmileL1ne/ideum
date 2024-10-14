package main

import (
	"forum/config"
	"forum/pkg/database/sqlite3"
	"forum/pkg/sesm"
	"forum/pkg/sesm/sqlite3store"
	"log"
	"net"
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

	server := &http.Server{
		Addr:           net.JoinHostPort("", cfg.Addr),
		Handler:        routes.Register(),
		MaxHeaderBytes: 1 << 20, // 1 mb
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

	logger.Printf("Starting the server on port: %s", cfg.Addr)
	err = server.ListenAndServe()
	logger.Fatalf("Listen and serve error:%v", err)
}
