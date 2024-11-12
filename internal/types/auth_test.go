package types

import (
	"fmt"
	"os"
	"testing"
)

func TestWriteAuthToFile(t *testing.T) {
	// Create a temp file for the write failure test
	tmpFile, err := os.CreateTemp("", "auth_test")
	if err != nil {
		t.Fatal(err)
	}

	// Close it and make it read-only
	tmpFile.Close()
	if err := os.Chmod(tmpFile.Name(), 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name        string
		authData    Auth
		filePath    string
		expectErr   bool
		expectedErr string
	}{
		{
			name: "Invalid File Path",
			authData: Auth{
				RefreshToken:     "sampleRefreshToken",
				Token:            "sampleToken",
				ExpiresIn:        "3600",
				RefreshExpiresIn: "3600",
			},
			filePath:    "/invalid/path/to/file.json",
			expectErr:   true,
			expectedErr: "failed to create or open file",
		},
		{
			name: "Failed to Create File",
			authData: Auth{
				RefreshToken:     "sampleRefreshToken",
				Token:            "sampleToken",
				ExpiresIn:        "3600",
				RefreshExpiresIn: "3600",
			},
			filePath:    tmpFile.Name(), // Use our read-only temp file
			expectErr:   true,
			expectedErr: fmt.Sprintf("failed to create or open file: open %s: Access is denied.", tmpFile.Name()),
		},
		{
			name: "Successful Write",
			authData: Auth{
				RefreshToken:     "sampleRefreshToken",
				Token:            "sampleToken",
				ExpiresIn:        "3600",
				RefreshExpiresIn: "3600",
			},
			filePath:    "/tmp/auth.json",
			expectErr:   false,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call WriteAuthToFile
			err := WriteAuthToFile(tt.filePath, tt.authData)

			// Check if an error was expected
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			// If an error was expected, ensure it's the correct error
			if err != nil && tt.expectErr && err.Error()[:len(tt.expectedErr)] != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err.Error())
			}

			// If no error was expected, verify the file was written
			if err == nil && !tt.expectErr {
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Errorf("expected file to be created, but it doesn't exist")
				}
			}
		})
	}
}
