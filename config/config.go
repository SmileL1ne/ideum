package config

// TODO: parse config info from .yml config file -> change NewConfig() to use 'cleanenv' package to parse config file

type (
	// Project config
	Config struct {
		App
		Http
		Database
	}

	// Information about the app
	App struct {
		Name    string
		Version string
	}

	// Http related info
	Http struct {
		Addr      string
		StaticDir string
	}

	Database struct {
		DSN string
	}
)

// NewConfig returns config
func NewConfig() *Config {
	cfg := &Config{
		App{
			Name:    "Idearoom",
			Version: "1.0.0",
		},
		Http{
			Addr:      ":5000",
			StaticDir: "./web/static",
		},
		Database{
			DSN: "file:./internal/database/rabbit.db?foreign_keys=on",
		},
	}
	return cfg
}
