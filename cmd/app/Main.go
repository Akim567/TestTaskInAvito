package main

import (
	"log"

	"TestTaskInAvito/internal/app"
	"TestTaskInAvito/internal/config"
)

func main() {
	cfg := config.Load()

	a, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
