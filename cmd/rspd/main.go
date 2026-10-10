package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/config"
	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/rename"
	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/workflow"
	"github.com/spf13/cobra"
)

var (
	configPath string
	dryRun     bool
	targetPath string
)

func main() {
	os.Exit(run())
}

// run は main の中身 (os.Exit を呼ばないので test から実行できる。exit は main が担う)。
func run() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

var rootCmd = &cobra.Command{
	Use:   "rspd",
	Short: "Rename Sony photo directories",
	Long:  `A tool to rename Sony camera photo directories with proper date formatting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfiguration()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return runRenameOnly(cfg)
	},
}

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Run full workflow: copy, rename, and delete",
	Long:  `Execute the complete workflow: copy files from source, rename directories, and delete originals.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfiguration()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if dryRun {
			log.Println("=== DRY RUN MODE ===")
			log.Println("No actual changes will be made")
		}
		log.Println("Starting workflow: copy, rename, and delete")
		if err := workflow.Run(cfg, dryRun); err != nil {
			return fmt.Errorf("workflow failed: %w", err)
		}
		log.Println("Workflow completed successfully")
		return nil
	},
}

var backupCleanupCmd = &cobra.Command{
	Use:   "backup-cleanup",
	Short: "Delete files from backup SD card and eject",
	Long:  `Clean up the backup SD card by deleting copied files and ejecting the card.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfiguration()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if dryRun {
			log.Println("=== DRY RUN MODE ===")
			log.Println("No actual changes will be made")
		}
		log.Println("Starting backup cleanup")
		if err := workflow.RunBackupCleanup(cfg, dryRun); err != nil {
			return fmt.Errorf("backup cleanup failed: %w", err)
		}
		log.Println("Backup cleanup completed successfully")
		return nil
	},
}

var createConfigCmd = &cobra.Command{
	Use:   "create-config",
	Short: "Create a default configuration file",
	Long:  `Generate a default configuration file (config.yaml) in the current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleCreateConfig()
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be done without making any changes")

	// Root command flags
	rootCmd.Flags().StringVarP(&targetPath, "path", "p", "", "Target path to rename directories (overrides config)")

	// Add subcommands
	rootCmd.AddCommand(workflowCmd)
	rootCmd.AddCommand(backupCleanupCmd)
	rootCmd.AddCommand(createConfigCmd)
}

func handleCreateConfig() error {
	defaultConfig := config.Default()
	configFile := "config.yaml"
	if err := config.Save(defaultConfig, configFile); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	log.Printf("Default configuration file created: %s", configFile)
	return nil
}

func loadConfiguration() (*config.Config, error) {
	if configPath != "" {
		return config.Load(configPath)
	}
	return config.LoadOrDefault()
}

func runRenameOnly(cfg *config.Config) error {
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
