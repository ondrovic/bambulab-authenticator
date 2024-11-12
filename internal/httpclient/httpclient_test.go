package httpclient

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ondrovic/bambulab-authenticator/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define a mock client that implements the HTTPClient interface.
type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// Helper to create a *http.Response with a given status code, body, and optional headers.
func createMockResponse(statusCode int, body string, headers map[string]string) *http.Response {
	resp := httptest.NewRecorder()
	resp.WriteHeader(statusCode)
	resp.Body = bytes.NewBufferString(body)
	for k, v := range headers {
		resp.Header().Set(k, v)
	}

	return resp.Result()
}

// containsSubstring is a helper function for partial error message matching.
func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Helper function to create a test cookie
func createCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
	}
}

func TestInitClient(t *testing.T) {
	// Arrange
	authToken := "test-auth-token"
	expectedTransport := &transportWithAuth{
		authToken: authToken,
		rt:        http.DefaultTransport,
	}

	// Act
	err := InitClient(authToken)

	// Assert
	require.NoError(t, err, "InitClient should not return an error")
	require.NotNil(t, Client, "Client should be initialized")

	// Type assertion to check if Client is *http.Client
	httpClient, ok := Client.(*http.Client)
	require.True(t, ok, "Client should be of type *http.Client")

	assert.IsType(t, &transportWithAuth{}, httpClient.Transport, "Client transport should be of type *transportWithAuth")

	// Type assertion to check transport properties
	actualTransport, ok := httpClient.Transport.(*transportWithAuth)
	require.True(t, ok, "Transport should be of type *transportWithAuth")
	assert.Equal(t, expectedTransport.authToken, actualTransport.authToken, "Auth token should match")
	assert.Equal(t, expectedTransport.rt, actualTransport.rt, "RoundTripper should match")
}

func TestRoundTrip_WithAuthToken(t *testing.T) {
	// Arrange
	authToken := "test-auth-token"
	InitClient(authToken)

	// Create a test server that echoes back the request
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token "+authToken, r.Header.Get("Authorization"), "Authorization header should be set correctly")
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Create a request to the test server
	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	require.NoError(t, err, "Failed to create request")

	// Act
	resp, err := Client.Do(req)

	// Assert
	require.NoError(t, err, "Client.Do should not return an error")
	require.NotNil(t, resp, "Client.Do should return a non-nil response")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status code should be OK")
}

func TestRoundTrip_WithoutAuthToken(t *testing.T) {
	// Arrange
	InitClient("") // Initialize with an empty auth token

	// Create a test server that echoes back the request
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"), "Authorization header should not be set")
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Create a request to the test server
	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	require.NoError(t, err, "Failed to create request")

	// Act
	resp, err := Client.Do(req)

	// Assert
	require.NoError(t, err, "Client.Do should not return an error")
	require.NotNil(t, resp, "Client.Do should return a non-nil response")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status code should be OK")
}

func TestAddDefaultHeadersToRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Call func to add default headers
	addDefaultHeadersToRequest(req)

	// define expected headers
	expectedHeaders := defaultHeaders

	// Compare the headers
	if !reflect.DeepEqual(req.Header, expectedHeaders) {
		t.Errorf("Headers do not match. Expected: %v, got: %v", expectedHeaders, req.Header)
	}
}

func TestMapCookiesToResponse(t *testing.T) {
	tests := []struct {
		name        string
		cookies     []*http.Cookie
		expected    *types.LoginResponse
		expectError bool
	}{
		{
			name: "All cookies present",
			cookies: []*http.Cookie{
				createCookie("token", "access-token-value"),
				createCookie("refreshToken", "refresh-token-value"),
				createCookie("expiresIn", "3600"),
				createCookie("refreshExpiresIn", "3600"),
			},
			expected: &types.LoginResponse{
				AccessToken:      "access-token-value",
				RefreshToken:     "refresh-token-value",
				ExpiresIn:        3600,
				RefreshExpiresIn: 3600,
			},
			expectError: false,
		},
		{
			name: "Missing expiresIn",
			cookies: []*http.Cookie{
				createCookie("token", "access-token-value"),
				createCookie("refreshToken", "refresh-token-value"),
			},
			expected: &types.LoginResponse{
				AccessToken:  "access-token-value",
				RefreshToken: "refresh-token-value",
			},
			expectError: false,
		},
		{
			name: "Invalid expiresIn",
			cookies: []*http.Cookie{
				createCookie("token", "access-token-value"),
				createCookie("refreshToken", "refresh-token-value"),
				createCookie("expiresIn", "invalid"),
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid refreshExpiresIn",
			cookies: []*http.Cookie{
				createCookie("token", "access-token-value"),
				createCookie("refreshToken", "refresh-token-value"),
				createCookie("expiresIn", "3600"),
				createCookie("refreshExpiresIn", "invalid"),
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "No cookies provided",
			cookies:     []*http.Cookie{},
			expected:    &types.LoginResponse{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MapCookiesToResponse(tt.cookies)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.AccessToken, result.AccessToken)
				assert.Equal(t, tt.expected.RefreshToken, result.RefreshToken)
				assert.Equal(t, tt.expected.ExpiresIn, result.ExpiresIn)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		payload        []byte
		mockResponse   *http.Response
		mockError      error
		expectedResult *types.LoginResponse
		expectedError  string
	}{
		{
			name:    "Successful request",
			method:  http.MethodPost,
			url:     "http://example.com/login",
			payload: []byte(`{"username":"user","password":"pass"}`),
			mockResponse: createMockResponse(http.StatusOK, `{"accessToken":"mockToken"}`, map[string]string{
				"Content-Type": "application/json",
			}),
			expectedResult: &types.LoginResponse{AccessToken: "mockToken"},
		},
		{
			name:          "HTTP client error",
			method:        http.MethodGet,
			url:           "http://example.com/fail",
			payload:       nil,
			mockError:     errors.New("network error"),
			expectedError: "request failed: network error",
		},
		{
			name:    "Invalid JSON response",
			method:  http.MethodGet,
			url:     "http://example.com/invalid-json",
			payload: nil,
			mockResponse: createMockResponse(http.StatusOK, `invalid-json`, map[string]string{
				"Content-Type": "application/json",
			}),
			expectedError: "failed to unmarshal response body: invalid character 'i' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock client with the desired behavior.
			mockClient := &mockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockResponse, nil
				},
			}

			// Set the mock client in place of the real client
			Client = mockClient

			// Run the Request function
			result, err := Request(tt.method, tt.url, tt.payload)

			// Check the expected error
			if tt.expectedError != "" {
				if err == nil || !containsSubstring(err.Error(), tt.expectedError) {
					t.Fatalf("expected error %q, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check the expected result
			if tt.expectedResult != nil && (result == nil || result.AccessToken != tt.expectedResult.AccessToken) {
				t.Fatalf("expected result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
