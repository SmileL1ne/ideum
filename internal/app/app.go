package app

import (
	"database/sql"
	"forum/config"
	postRepository "forum/internal/repository/post"
	postService "forum/internal/service/post"
	"log/slog"
	"net/http"
	"os"

	h "forum/internal/controller/http"

	_ "github.com/mattn/go-sqlite3"
)

/*
	TODO: Add graceful shutdown
*/

func Run(cfg *config.Config) {
	// Logger init
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Database connection
	db, err := OpenDB(cfg.Sqlite.Dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Services
	postservice := postService.New(postRepository.New(db))

	// Server creation
	server := &http.Server{
		Addr:     "127.0.0.1" + cfg.Http.Addr,
		Handler:  h.NewRouter(logger, postservice),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

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
