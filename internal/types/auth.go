package types

import (
	"encoding/json"
	"fmt"
	"os"
)

type Auth struct {
	RefreshToken     string `json:"refreshToken,omitempty"`
	Token            string `json:"token,omitempty"`
	ExpiresIn        string `json:"expiresIn"`
	RefreshExpiresIn string `json:"refreshExpiresIn"`
}

// WriteAuthToFile serializes the Auth struct to JSON and writes it to a file
func WriteAuthToFile(path string, data Auth) error {
	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal auth data: %v", err)
	}

	// Create or open the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create or open file: %v", err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}
