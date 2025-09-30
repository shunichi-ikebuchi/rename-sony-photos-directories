// Package rename provides directory renaming functionality for Sony camera date formats.
package rename

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	// ExpectedDirNameLength is the expected length of Sony camera directory names
	ExpectedDirNameLength = 8 // Expected format: YYYMMDD (8 digits)
)

// IsValidDateDir checks if the directory name matches the expected Sony camera date format.
// Sony camera format: 8 digits (YYYMMDD where YYY is year from 2000)
func IsValidDateDir(name string) bool {
	matched, _ := regexp.MatchString(`^\d{8}$`, name)
	return matched
}

// ConvertDirName converts Sony camera date format (0YYMMDD0) to yyyy-mm-dd format.
// Sony format: 0YYMMDD0 where YY is the last 2 digits of the year
// Example: 02512310 -> 2025-12-31 (first and last digits are padding)
func ConvertDirName(name, currentCentury string) (string, error) {
	if len(name) != ExpectedDirNameLength {
		return "", fmt.Errorf("invalid directory name length: %s", name)
	}

	// Convert 0YYMMDD0 to YYYY-MM-DD
	// Format: [0][YY][MM][DD][0]
	// Positions: 0  1-2  3-4  5-6  7
	yearSuffix := name[1:3] // Extract YY (positions 1-2)
	year := currentCentury + yearSuffix
	month := name[3:5] // Positions 3-4
	day := name[5:7]   // Positions 5-6

	return fmt.Sprintf("%s-%s-%s", year, month, day), nil
}

// Directories renames directories in the specified path from Sony camera format to yyyy-mm-dd format.
func Directories(targetPath string) error {
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", targetPath, err)
	}

	currentYear := fmt.Sprintf("%d", time.Now().Year())
	currentCentury := currentYear[:2] // First 2 digits (e.g., "20")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()
		if !IsValidDateDir(dirName) {
			log.Printf("Skipping directory (invalid format): %s", dirName)
			continue
		}

		newName, err := ConvertDirName(dirName, currentCentury)
		if err != nil {
			log.Printf("Error converting directory name %s: %v", dirName, err)
			continue
		}

		oldPath := filepath.Join(targetPath, dirName)
		newPath := filepath.Join(targetPath, newName)

		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("Error renaming %s to %s: %v", oldPath, newPath, err)
			continue
		}

		log.Printf("Renamed: %s -> %s", dirName, newName)
	}

	return nil
}
