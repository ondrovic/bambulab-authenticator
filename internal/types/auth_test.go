package types

import (
	"os"
	"testing"
)

func TestWriteAuthToFile(t *testing.T) {
	// Create a temp file for the write failure test
	tmpFile, err := os.CreateTemp("", "auth_test")
	if err != nil {
		t.Fatal(err)
	}
	tmpFilePath := tmpFile.Name()
	tmpFile.Close()
	if err := os.Chmod(tmpFilePath, 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFilePath)

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
			expectedErr: "failed to create or open file: open /invalid/path/to/file.json: The system cannot find the path specified.",
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

			// If no error was expected, verify the file was written
			if err == nil && !tt.expectErr {
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Errorf("expected file to be created, but it doesn't exist")
				}
			}
		})
	}
}
