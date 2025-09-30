// Package workflow provides workflow operations for photo management.
package workflow

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/config"
	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/rename"
)

// copyDir recursively copies a directory tree
func CopyDir(src, dst string, dryRun bool) error {
	if dryRun {
		log.Printf("[DRY RUN] Would copy directory: %s -> %s", src, dst)
		return nil
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath, false); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	if err := os.Chmod(dst, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// RemoveContents removes all contents of a directory but keeps the directory itself
func RemoveContents(dir string, dryRun bool) error {
	if dryRun {
		log.Printf("[DRY RUN] Would remove contents of: %s", dir)
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return nil
}

// EjectVolume ejects a volume (macOS only)
func EjectVolume(volumeName string, dryRun bool) error {
	if dryRun {
		log.Printf("[DRY RUN] Would eject volume: %s", volumeName)
		return nil
	}
	if runtime.GOOS != "darwin" {
		log.Printf("Eject not supported on %s, skipping", runtime.GOOS)
		return nil
	}

	cmd := exec.Command("diskutil", "eject", volumeName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to eject volume %s: %w\nOutput: %s", volumeName, err, string(output))
	}

	log.Printf("Successfully ejected volume: %s", volumeName)
	return nil
}

// CheckDirectoryExists checks if a directory exists
func CheckDirectoryExists(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("failed to check directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	return nil
}

// runWorkflow executes the copy-rename-delete workflow
func Run(config *config.Config, dryRun bool) error {
	tmpDir := config.TmpDir
	sourceDCIM := filepath.Join(config.TargetPath, "DCIM")

	// Check if source directory exists
	if err := CheckDirectoryExists(config.DestinationPath); err != nil {
		return fmt.Errorf("destination check failed: %w", err)
	}

	if err := CheckDirectoryExists(sourceDCIM); err != nil {
		return fmt.Errorf("source DCIM check failed: %w", err)
	}

	// Create temporary directory
	if !dryRun {
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return fmt.Errorf("failed to create temporary directory: %w", err)
		}
	} else {
		log.Printf("[DRY RUN] Would create temporary directory: %s", tmpDir)
	}

	log.Printf("Copying photos from %s to %s", sourceDCIM, tmpDir)
	if err := CopyDir(sourceDCIM, tmpDir, dryRun); err != nil {
		return fmt.Errorf("failed to copy files to temp directory: %w", err)
	}

	log.Printf("Renaming directories in %s", tmpDir)
	if !dryRun {
		if err := rename.Directories(tmpDir); err != nil {
			return fmt.Errorf("failed to rename directories: %w", err)
		}
	} else {
		log.Printf("[DRY RUN] Would rename directories in: %s", tmpDir)
	}

	log.Printf("Copying renamed directories to %s", config.DestinationPath)
	if err := CopyDir(tmpDir, config.DestinationPath, dryRun); err != nil {
		return fmt.Errorf("failed to copy to destination: %w", err)
	}

	log.Printf("Deleting photos from source: %s", sourceDCIM)
	if err := RemoveContents(sourceDCIM, dryRun); err != nil {
		return fmt.Errorf("failed to delete source files: %w", err)
	}

	log.Printf("Cleaning up temporary directory: %s", tmpDir)
	if err := RemoveContents(tmpDir, dryRun); err != nil {
		return fmt.Errorf("failed to clean temporary directory: %w", err)
	}

	// Extract volume name from path (e.g., "/Volumes/1-1" -> "1-1")
	volumeName := filepath.Base(config.TargetPath)
	log.Printf("Ejecting volume: %s", volumeName)
	if err := EjectVolume(volumeName, dryRun); err != nil {
		log.Printf("Warning: %v", err)
	}

	return nil
}

// runBackupCleanup deletes all files from backup SD card and ejects it
func RunBackupCleanup(config *config.Config, dryRun bool) error {
	backupDCIM := filepath.Join(config.BackupPath, "DCIM")

	// Check if backup directory exists
	if err := CheckDirectoryExists(config.BackupPath); err != nil {
		return fmt.Errorf("backup path check failed: %w", err)
	}

	if err := CheckDirectoryExists(backupDCIM); err != nil {
		return fmt.Errorf("backup DCIM check failed: %w", err)
	}

	log.Printf("Deleting photos from backup: %s", backupDCIM)
	if err := RemoveContents(backupDCIM, dryRun); err != nil {
		return fmt.Errorf("failed to delete backup files: %w", err)
	}

	// Extract volume name from path
	volumeName := filepath.Base(config.BackupPath)
	log.Printf("Ejecting backup volume: %s", volumeName)
	if err := EjectVolume(volumeName, dryRun); err != nil {
		log.Printf("Warning: %v", err)
	}

	return nil
}
