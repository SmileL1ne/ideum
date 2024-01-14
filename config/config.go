package config

// TODO: parse config info from .yml config file -> change NewConfig() to use 'cleanenv' package to parse config file

type (
	// Project config
	Config struct {
		App
		Http
		Sqlite
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

	Sqlite struct {
		Dsn string
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
		Sqlite{
			Dsn: "file:./forum.db?cache=shared&mode=rwc",
		},
	}
	return cfg
}
