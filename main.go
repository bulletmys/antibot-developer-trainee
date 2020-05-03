package main

import (
	"antibot-trainee/internal/app/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to export env vars: %v", err)
	}

	server.RunServer()
}
