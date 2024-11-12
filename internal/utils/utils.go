package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ondrovic/bambulab-authenticator/internal/types"
)

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return s == ""
}

// SaveLoginResponseToFile serializes the LoginResponse struct to JSON and saves it to the given file path
func SaveLoginResponseToFile(loginResponse types.LoginResponse, path string) error {

	fullPath := filepath.Join(path, "auth.json")
	// Marshal the struct to JSON (with indentation for readability)
	jsonData, err := json.MarshalIndent(loginResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(fullPath, jsonData, 0644) // 0644 is the file permission mode
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	fmt.Printf("Auth data saved to %s", fullPath)

	return nil
}
