package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/config"
	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/rename"
	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/workflow"
)

func handleCreateConfig() error {
	defaultConfig := config.Default()
	configFile := "config.yaml"
	if err := config.Save(defaultConfig, configFile); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	log.Printf("Default configuration file created: %s", configFile)
	return nil
}

func loadConfiguration(configPath string) (*config.Config, error) {
	if configPath != "" {
		return config.Load(configPath)
	}
	return config.LoadOrDefault()
}

func runRenameOnly(cfg *config.Config, targetPath string, dryRun bool) error {
	path := cfg.TargetPath
	if targetPath != "" {
		path = targetPath
	}

	if path == "" {
		path = "."
	}

	if dryRun {
		log.Println("=== DRY RUN MODE ===")
		log.Println("No actual changes will be made")
		log.Printf("Would rename directories in: %s", path)
		return nil
	}

	log.Printf("Renaming directories in: %s", path)
	if err := rename.Directories(path); err != nil {
		return fmt.Errorf("failed to rename directories: %w", err)
	}

	log.Println("Directory renaming completed successfully")
	return nil
}

func main() {
	// Command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	targetPath := flag.String("path", "", "Target path to rename directories (overrides config)")
	createConfig := flag.Bool("create-config", false, "Create a default configuration file")
	workflowFlag := flag.Bool("workflow", false, "Run full workflow: copy, rename, and delete")
	backupCleanup := flag.Bool("backup-cleanup", false, "Delete files from backup SD card and eject")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without making any changes")
	flag.Parse()

	// Create default config if requested
	if *createConfig {
		if err := handleCreateConfig(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Load configuration
	cfg, err := loadConfiguration(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Execute based on command flags
	if *workflowFlag {
		if *dryRun {
			log.Println("=== DRY RUN MODE ===")
			log.Println("No actual changes will be made")
		}
		log.Println("Starting workflow: copy, rename, and delete")
		if err := workflow.Run(cfg, *dryRun); err != nil {
			log.Fatalf("Workflow failed: %v", err)
		}
		log.Println("Workflow completed successfully")
	} else if *backupCleanup {
		if *dryRun {
			log.Println("=== DRY RUN MODE ===")
			log.Println("No actual changes will be made")
		}
		log.Println("Starting backup cleanup")
		if err := workflow.RunBackupCleanup(cfg, *dryRun); err != nil {
			log.Fatalf("Backup cleanup failed: %v", err)
		}
		log.Println("Backup cleanup completed successfully")
	} else {
		if err := runRenameOnly(cfg, *targetPath, *dryRun); err != nil {
			log.Fatal(err)
		}
	}
}
