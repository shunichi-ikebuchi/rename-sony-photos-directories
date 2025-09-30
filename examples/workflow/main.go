// Package main demonstrates the full workflow.
package main

import (
	"log"

	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/config"
	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/workflow"
)

func main() {
	// Load configuration
	cfg, err := config.LoadOrDefault()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run full workflow (copy, rename, delete)
	log.Println("Starting workflow...")
	if err := workflow.Run(cfg, false); err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	log.Println("Workflow completed successfully!")
}
