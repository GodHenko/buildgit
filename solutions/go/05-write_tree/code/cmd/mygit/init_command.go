package main

import (
	"fmt"
	"os"
)

func initCmd() error {
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %v: %v", dir, err)
		}
	}

	headFileContents := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		return fmt.Errorf("write .git/HEAD: %v", err)
	}

	fmt.Println("Initialized git repository")

	return nil
}
