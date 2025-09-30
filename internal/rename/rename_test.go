package rename

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidDateDir(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid 8 digits", "12345678", true},
		{"valid date format", "02501231", true},
		{"too short", "1234567", false},
		{"too long", "123456789", false},
		{"contains letters", "1234567a", false},
		{"contains special chars", "1234-678", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidDateDir(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidDateDir(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertDirName(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		currentCentury string
		expected       string
		shouldError    bool
	}{
		{
			name:           "valid date 2025-12-31",
			input:          "02512310",
			currentCentury: "20",
			expected:       "2025-12-31",
			shouldError:    false,
		},
		{
			name:           "valid date 2024-06-15",
			input:          "02406150",
			currentCentury: "20",
			expected:       "2024-06-15",
			shouldError:    false,
		},
		{
			name:           "valid date 2020-10-20",
			input:          "02010200",
			currentCentury: "20",
			expected:       "2020-10-20",
			shouldError:    false,
		},
		{
			name:           "invalid length",
			input:          "1234567",
			currentCentury: "20",
			expected:       "",
			shouldError:    true,
		},
		{
			name:           "too long",
			input:          "123456789",
			currentCentury: "20",
			expected:       "",
			shouldError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertDirName(tt.input, tt.currentCentury)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ConvertDirName(%q, %q) expected error, got nil", tt.input, tt.currentCentury)
				}
			} else {
				if err != nil {
					t.Errorf("ConvertDirName(%q, %q) unexpected error: %v", tt.input, tt.currentCentury, err)
				}
				if result != tt.expected {
					t.Errorf("ConvertDirName(%q, %q) = %q, want %q", tt.input, tt.currentCentury, result, tt.expected)
				}
			}
		})
	}
}

func TestRenameDirectories(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "rename-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directories
	testDirs := []string{
		"02512310",   // Valid directory -> 2025-12-31
		"02406150",   // Valid directory -> 2024-06-15
		"invaliddir", // Invalid directory (should be skipped)
		"1234567",    // Invalid length (should be skipped)
	}

	for _, dir := range testDirs {
		dirPath := filepath.Join(tmpDir, dir)
		if err := os.Mkdir(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create a test file (should be ignored)
	testFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run Directories
	if err := Directories(tmpDir); err != nil {
		t.Fatalf("Directories failed: %v", err)
	}

	// Check that valid directories were renamed
	expectedDirs := map[string]bool{
		"2025-12-31": true,
		"2024-06-15": true,
		"invaliddir": true, // Should still exist (not renamed)
		"1234567":    true, // Should still exist (not renamed)
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	foundDirs := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			foundDirs[entry.Name()] = true
		}
	}

	for expectedDir := range expectedDirs {
		if !foundDirs[expectedDir] {
			t.Errorf("Expected directory %q not found", expectedDir)
		}
	}

	// Verify original valid directories no longer exist
	if foundDirs["02512310"] {
		t.Error("Original directory 02512310 should have been renamed")
	}
	if foundDirs["02406150"] {
		t.Error("Original directory 02406150 should have been renamed")
	}

	// Verify test file still exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Test file should not have been removed")
	}
}

func TestRenameDirectoriesNonExistentPath(t *testing.T) {
	err := Directories("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
}
