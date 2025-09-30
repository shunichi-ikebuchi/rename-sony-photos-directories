// Package main demonstrates basic directory renaming.
package main

import (
	"log"

	"github.com/shunichi-ikebuchi/rename-sony-photos-directories/internal/rename"
)

func main() {
	// Rename directories in current directory
	if err := rename.Directories("."); err != nil {
		log.Fatalf("Failed to rename directories: %v", err)
	}

	log.Println("Directory renaming completed!")
}
