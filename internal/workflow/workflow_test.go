package workflow

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "copy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	dstFile := filepath.Join(tmpDir, "destination.txt")
	if err := copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify destination file exists and has correct content
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if !bytes.Equal(dstContent, content) {
		t.Errorf("Content mismatch: got %q, want %q", string(dstContent), string(content))
	}
}

func TestCopyDir(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "copy-dir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	// Create test files
	testFile1 := filepath.Join(srcDir, "file1.txt")
	if err := os.WriteFile(testFile1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	subDir := filepath.Join(srcDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	testFile2 := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(testFile2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Copy directory
	dstDir := filepath.Join(tmpDir, "destination")
	if err := CopyDir(srcDir, dstDir, false); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Verify files were copied
	dstFile1 := filepath.Join(dstDir, "file1.txt")
	if _, err := os.Stat(dstFile1); os.IsNotExist(err) {
		t.Error("file1.txt was not copied")
	}

	dstFile2 := filepath.Join(dstDir, "subdir", "file2.txt")
	if _, err := os.Stat(dstFile2); os.IsNotExist(err) {
		t.Error("subdir/file2.txt was not copied")
	}
}

func TestCopyDirDryRun(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "copy-dir-dryrun-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	dstDir := filepath.Join(tmpDir, "destination")

	// Run CopyDir with dry-run
	if err := CopyDir(srcDir, dstDir, true); err != nil {
		t.Fatalf("CopyDir with dry-run failed: %v", err)
	}

	// Verify destination was NOT created
	if _, err := os.Stat(dstDir); !os.IsNotExist(err) {
		t.Error("Destination directory should not exist in dry-run mode")
	}
}

func TestRemoveContents(t *testing.T) {
	// Create temporary directory with files
	tmpDir, err := os.MkdirTemp("", "remove-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testDir := filepath.Join(tmpDir, "testdir")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}

	// Create test files
	file1 := filepath.Join(testDir, "file1.txt")
	if err := os.WriteFile(file1, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	subDir := filepath.Join(testDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Remove contents
	if err := RemoveContents(testDir, false); err != nil {
		t.Fatalf("RemoveContents failed: %v", err)
	}

	// Verify directory is empty but still exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory should still exist")
	}

	entries, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Directory should be empty, but has %d entries", len(entries))
	}
}

func TestRemoveContentsDryRun(t *testing.T) {
	// Create temporary directory with files
	tmpDir, err := os.MkdirTemp("", "remove-dryrun-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run RemoveContents with dry-run
	if err := RemoveContents(tmpDir, true); err != nil {
		t.Fatalf("RemoveContents with dry-run failed: %v", err)
	}

	// Verify file still exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("File should still exist in dry-run mode")
	}
}

func TestCheckDirectoryExists(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "check-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test existing directory
	if err := CheckDirectoryExists(tmpDir); err != nil {
		t.Errorf("CheckDirectoryExists failed for existing directory: %v", err)
	}

	// Test non-existent directory
	nonExistent := filepath.Join(tmpDir, "nonexistent")
	if err := CheckDirectoryExists(nonExistent); err == nil {
		t.Error("CheckDirectoryExists should fail for non-existent directory")
	}

	// Test file (not directory)
	testFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := CheckDirectoryExists(testFile); err == nil {
		t.Error("CheckDirectoryExists should fail for file")
	}
}
