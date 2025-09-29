package main

import (
	"context"
	"log"

	"github.com/0xm0-v1/sik6/internal/app"
	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/0xm0-v1/sik6/internal/httpserver"
)

func main() {
	if err := config.LoadDevDotEnv(); err != nil {
		log.Printf("warning: could not load .env.development: %v", err)
	}
	cfg := config.New()
	log.Printf("env loaded successfully")

	handler := app.NewHTTPHandler(cfg)

	if err := httpserver.Run(context.Background(), cfg, handler); err != nil {
		log.Fatalf("application run error: %v", err)
	}
}
