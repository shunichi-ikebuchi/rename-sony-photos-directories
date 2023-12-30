package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// This program renames directories in the current directory
func main() {
	dir, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	currentYear := fmt.Sprintf("%d", time.Now().Year())
	first_three_digits := (currentYear[:3])
	for _, entry := range dir {
		if entry.IsDir() {
			tmp := first_three_digits + entry.Name()[3:]
			tmp = tmp[:4] + "-" + tmp[4:6] + "-" + tmp[6:]
			os.Rename(entry.Name(), tmp)
		}
	}
}
