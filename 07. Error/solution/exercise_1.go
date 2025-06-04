package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Define custom error types
type APIError struct {
	StatusCode int
	URL        string
	Message    string
	Err        error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d) on %s: %s", e.StatusCode, e.URL, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// Define sentinel errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrTimeout      = errors.New("request timed out")
)

// APIClient for making HTTP requests
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// NewAPIClient creates a new client with default settings
func NewAPIClient(baseURL, token string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		AuthToken: token,
	}
}

// GetUser fetches a user from the API
func (c *APIClient) GetUser(userID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization if available
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// Handle timeout specifically
		if errors.Is(err, http.ErrHandlerTimeout) {
			return nil, ErrTimeout
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle different status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Success - parse the JSON
		var user map[string]interface{}
		if err := json.Unmarshal(body, &user); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return user, nil

	case http.StatusNotFound:
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			URL:        url,
			Message:    "User not found",
			Err:        ErrNotFound,
		}

	case http.StatusUnauthorized:
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			URL:        url,
			Message:    "Invalid or expired token",
			Err:        ErrUnauthorized,
		}

	default:
		// Generic error for other status codes
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			URL:        url,
			Message:    fmt.Sprintf("API returned status %d", resp.StatusCode),
			Err:        errors.New("unexpected API response"),
		}
	}
}

func main() {
	// Create a client
	client := NewAPIClient("https://api.example.com", "valid-token")

	// Make a request
	user, err := client.GetUser("123")

	// Handle errors with appropriate type checks
	if err != nil {
		var apiErr *APIError

		switch {
		case errors.Is(err, ErrNotFound):
			fmt.Println("User not found")

		case errors.Is(err, ErrUnauthorized):
			fmt.Println("Please log in again")

		case errors.Is(err, ErrTimeout):
			fmt.Println("Request timed out, please try again")

		case errors.As(err, &apiErr):
			fmt.Printf("API error (%d): %s\n", apiErr.StatusCode, apiErr.Message)

		default:
			fmt.Printf("Unexpected error: %v\n", err)
		}

		return
	}

	// Process user data
	fmt.Printf("User: %v\n", user)
}
