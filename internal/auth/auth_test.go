package auth

import (
	"fmt"
	"testing"

	"github.com/ondrovic/bambulab-authenticator/internal/types"
)

func TestProcessLoginType(t *testing.T) {
	tests := []struct {
		name           string
		loginType      string
		expectError    bool
		err            error
		expectedResult *types.LoginResponse
	}{
		{
			name:           "Invalid type",
			loginType:      "invalid",
			expectError:    true,
			err:            fmt.Errorf("unknown login type: invalid"),
			expectedResult: &types.LoginResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock LoginResponse
			loginResponse := &types.LoginResponse{
				LoginType: tt.loginType,
				TfaKey:    "mock_tfa_key",
			}
			// Create a mock CliFlags
			opts := &types.CliFlags{
				UserAccount:  "test@example.com",
				UserPassword: "password123",
				UserRegion:   "us",
				OutputPath:   "output.json",
			}

			// Call the processLoginType function
			err := processLoginType(loginResponse, opts)

			// Verify the expected behavior
			if tt.expectError {
				if err == nil {
					t.Errorf("processLoginType() expected an error, but got nil")
				} else if err.Error() != tt.err.Error() {
					t.Errorf("processLoginType() returned unexpected error: got %v, want %v", err, tt.err)
				}
			} else {
				if err != nil {
					t.Errorf("processLoginType() returned an unexpected error: %v", err)
				}
			}
		})
	}
}
