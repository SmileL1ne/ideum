package main

import (
	"forum/config"
	"forum/internal/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
