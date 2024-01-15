package app

import (
	"database/sql"
	"forum/config"
	"log/slog"
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
	// Logger init
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Database connection
	db, err := OpenDB(cfg.Database.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Service
	r := repository.New(db)
	s := service.New(r)

	// Server creation
	server := &http.Server{
		Addr:     "127.0.0.1" + cfg.Http.Addr,
		Handler:  h.NewRouter(logger, s),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sigCh
		logger.Info("signal received", "signal", sig.String())
		db.Close()

		os.Exit(0)
	}()

	// Starting the server
	logger.Info("starting the server", slog.String("addr", cfg.Addr))
	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
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
