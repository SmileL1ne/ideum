package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	Config struct {
		App
		Http
		Database
		ExternalAuth
	}

	App struct {
		Name    string
		Version string
	}

	Http struct {
		Addr         string
		RateInterval int
		RateLimit    int
		RatePenalty  int
	}

	Database struct {
		DSN string
	}

	ExternalAuth struct {
		GoogleRedirectURL  string
		GoogleClientID     string
		GoogleClientSecret string
		GithubRedirectURL  string
		GithubClientID     string
		GithubClientSecret string
	}
)

// Load loads all required environments and returns ready config
func Load() *Config {
	rateInterval, err := strconv.Atoi(os.Getenv("HTTP_RATE_INTERVAL"))
	if err != nil {
		log.Fatal(err)
	}
	rateLimit, err := strconv.Atoi(os.Getenv("HTTP_RATE_LIMIT"))
	if err != nil {
		log.Fatal(err)
	}
	ratePenalty, err := strconv.Atoi(os.Getenv("HTTP_RATE_PENALTY"))
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		App{
			Name:    os.Getenv("APP_NAME"),
			Version: os.Getenv("APP_VERSION"),
		},
		Http{
			Addr:         os.Getenv("HTTP_ADDR"),
			RateInterval: rateInterval,
			RateLimit:    rateLimit,
			RatePenalty:  ratePenalty,
		},
		Database{
			DSN: os.Getenv("SQLITE3_DSN"),
		},
		ExternalAuth{
			GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),

			GithubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		},
	}
}

// Parse ".env" file and sets key-value pairs in it into system environment
func init() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("open env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				log.Print("ignoring .env line (invalid format):", line)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if err := os.Setenv(key, value); err != nil {
				log.Fatalf("set env variable: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("env parse scanner error: %v", err)
	}
}
