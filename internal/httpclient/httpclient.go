package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ondrovic/bambulab-authenticator/internal/consts"
	"github.com/ondrovic/bambulab-authenticator/internal/types"
)

// Client is a global variable holding the HTTP client used for making requests with authentication support.
// var Client *http.Client
// var Client interface {
//     Do(req *http.Request) (*http.Response, error)
// } = &http.Client{}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CustomResponse struct {
	Body    []byte
	Cookies []*http.Cookie
	Status  int
	Header  http.Header
}

var (
	Client         HTTPClient
	defaultHeaders = http.Header{
		"Content-Type": {"application/json"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0"},
		"Accept":       {"*/*"},
		"Connection":   {"keep-alive"},
		"Referer":      {string(consts.RefererURL)},
	}
)

type transportWithAuth struct {
	// authToken is the authentication token used for authorized requests.
	authToken string
	// rt is the underlying RoundTripper used for HTTP transport.
	rt http.RoundTripper
}

func addDefaultHeadersToRequest(req *http.Request) {
	for key, values := range defaultHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

func Request(method string, url string, payload []byte) (*types.LoginResponse, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	addDefaultHeadersToRequest(req)

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for empty body
	if resp.ContentLength == 0 || resp.Body == http.NoBody {
		return &types.LoginResponse{}, nil // return an empty LoginResponse
	}

	// Otherwise, read and unmarshal the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var loginResponse types.LoginResponse
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &loginResponse, nil
}

func CookieRequest(method string, url string, payload []byte) (*types.LoginResponse, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	addDefaultHeadersToRequest(req)

	resp, err := Client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %v", resp.Status)
	}

	return MapCookiesToResponse(resp.Cookies())
}

func MapCookiesToResponse(cookies []*http.Cookie) (*types.LoginResponse, error) {
	loginResponse := &types.LoginResponse{}

	for _, cookie := range cookies {
		switch cookie.Name {
		case "token":
			loginResponse.AccessToken = cookie.Value
		case "refreshToken":
			loginResponse.RefreshToken = cookie.Value
		case "expiresIn":
			expiresIn, err := strconv.Atoi(cookie.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid expiresIn value: %v", err)
			}
			loginResponse.ExpiresIn = expiresIn
		case "refreshExpiresIn":
			refreshExpiresIn, err := strconv.Atoi(cookie.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid refreshExpiresIn value: %v", err)
			}
			loginResponse.RefreshExpiresIn = refreshExpiresIn
		}
	}

	return loginResponse, nil
}

// InitClient initializes an HTTP client with authentication support using the provided authToken.
// It creates an HTTP client with a custom transport that includes the authentication token for requests.
// Parameters:
// - authToken: The authentication token to be used for authorized requests.
// Returns:
// - A pointer to the initialized http.Client.
// - An error if any issues occur during client initialization (returns nil in this implementation).
func InitClient(authToken string) error {
	Client = &http.Client{
		Transport: &transportWithAuth{
			authToken: authToken,
			rt:        http.DefaultTransport,
		},
	}

	return nil
}

// RoundTrip executes a single HTTP request using the transportWithAuth transport.
// If an authentication token is set, it adds an "Authorization" header to the request.
// Parameters:
// - req: The HTTP request to be sent.
// Returns:
// - A pointer to the http.Response received from the server.
// - An error if any issues occur during the request execution.
func (t *transportWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.authToken != consts.EMPTY_STRING {
		req.Header.Set("Authorization", "token "+t.authToken)
	}
	return t.rt.RoundTrip(req)
}
