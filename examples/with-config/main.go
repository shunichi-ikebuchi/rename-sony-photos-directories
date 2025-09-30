// Package main demonstrates using a custom configuration file.
package main

import (
	"log"

	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/config"
	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/rename"
)

func main() {
	// Load configuration from specific file
	cfg, err := config.Load("./my-config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Use configured path
	log.Printf("Renaming directories in: %s", cfg.TargetPath)
	if err := rename.Directories(cfg.TargetPath); err != nil {
		log.Fatalf("Failed to rename directories: %v", err)
	}

	log.Println("Directory renaming completed!")
}
