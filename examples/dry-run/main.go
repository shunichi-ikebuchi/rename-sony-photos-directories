// Package main demonstrates dry-run mode.
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

	// Run workflow in dry-run mode (no actual changes)
	log.Println("=== DRY RUN MODE ===")
	log.Println("Showing what would be done without making changes...")

	if err := workflow.Run(cfg, true); err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	log.Println("Dry-run completed!")
}
