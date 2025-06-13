package main

import (
	"log"
	"renhance-email-detector/internal/api"
)

func main() {
	if err := api.Start(); err != nil {
		log.Fatalf("failed to start API: %v", err)
	}
}
