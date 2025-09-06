package main

import (
	"log"
	"subcalc/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Config: %+v\n", cfg)
	log.Println("Starting service...")
}
