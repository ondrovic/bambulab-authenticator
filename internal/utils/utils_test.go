package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ondrovic/bambulab-authenticator/internal/types"
)

// helper containsErrorMessage
func containsErrorMessage(err error, substring string) bool {
	return fmt.Sprintf("%v", err) != "" && strings.Contains(fmt.Sprintf("%v", err), substring)
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "Empty string",
			s:        "",
			expected: true,
		},
		{
			name:     "Non-Empty string",
			s:        "hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEmpty(tt.s)
			if got != tt.expected {
				t.Errorf("IsEmpty() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestSaveLoginResponseToFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "login-response-test")
	if err != nil {
		t.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	testCases := []struct {
		name     string
		response types.LoginResponse
		wantErr  bool
		errMsg   string
	}{
		{
			name: "successful_save",
			response: types.LoginResponse{
				AccessToken:      "abc123",
				RefreshToken:     "def456",
				ExpiresIn:        3600,
				RefreshExpiresIn: 7200,
				TfaKey:           "ghi789",
				AccessMethod:     "password",
				LoginType:        "normal",
			},
			wantErr: false,
		},
		// {
		// 	name: "invalid_response_missing_fields",
		// 	response: types.LoginResponse{
		// 		AccessToken:      "",
		// 		RefreshToken:     "",
		// 		ExpiresIn:        0,
		// 		RefreshExpiresIn: 0,
		// 		TfaKey:           "",
		// 		AccessMethod:     "",
		// 		LoginType:        "",
		// 	},
		// 	wantErr: true,
		// 	errMsg:  "failed to marshal data:",
		// },
		{
			name: "file_write_error",
			response: types.LoginResponse{
				AccessToken:      "abc123",
				RefreshToken:     "def456",
				ExpiresIn:        3600,
				RefreshExpiresIn: 7200,
				TfaKey:           "ghi789",
				AccessMethod:     "password",
				LoginType:        "normal",
			},
			wantErr: true,
			errMsg:  "failed to write to file:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// For the file write error test, create a directory that the user cannot write to
			var path string
			if tc.name == "file_write_error" {
				path = "/this/path/does/not/exist"
			} else {
				path = tempDir
			}

			err := SaveLoginResponseToFile(tc.response, path)
			if (err != nil) != tc.wantErr {
				t.Errorf("SaveLoginResponseToFile() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr && err != nil && !containsErrorMessage(err, tc.errMsg) {
				t.Errorf("SaveLoginResponseToFile() error = %v, expected error message containing %q", err, tc.errMsg)
			}

			if !tc.wantErr {
				// Verify that the file was created
				filePath := filepath.Join(path, "auth.json")
				_, err := os.Stat(filePath)
				if os.IsNotExist(err) {
					t.Errorf("file was not created: %s", filePath)
				}
			}
		})
	}
}
