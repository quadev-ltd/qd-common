package util

import (
	"fmt"
	"os"
)

func ChangeCurrentWorkingDirectory(to string) (*string, error) {
	// Save current working directory and change it
	originalWD, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Failed to get current working directory: %s", err)
	}
	// Change working directory to the desired one
	if err := os.Chdir(to); err != nil {
		return nil, fmt.Errorf("Failed to change working directory: %s", err)
	}
	// Defer the reset of the working directory
	return &originalWD, nil
}
