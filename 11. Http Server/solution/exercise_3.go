package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ApiResponse represents a generic API response
type ApiResponse struct {
	Source  string
	Data    map[string]interface{}
	Error   error
	Latency time.Duration
}

// FetchAPI makes an HTTP request to the given URL and returns the response
func FetchAPI(url string, source string) ApiResponse {
	startTime := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(startTime)

	if err != nil {
		return ApiResponse{
			Source:  source,
			Error:   err,
			Latency: latency,
		}
	}
	defer resp.Body.Close()

	var data map[string]interface{}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ApiResponse{
				Source:  source,
				Error:   fmt.Errorf("failed to read response body: %w", err),
				Latency: latency,
			}
		}

		if err := json.Unmarshal(body, &data); err != nil {
			return ApiResponse{
				Source:  source,
				Error:   fmt.Errorf("failed to parse JSON: %w", err),
				Latency: latency,
			}
		}
	} else {
		return ApiResponse{
			Source:  source,
			Error:   fmt.Errorf("API returned status code %d", resp.StatusCode),
			Latency: latency,
		}
	}

	return ApiResponse{
		Source:  source,
		Data:    data,
		Latency: latency,
	}
}

func main() {
	// List of APIs to fetch (using httpbin for demonstration)
	apis := []struct {
		URL    string
		Source string
	}{
		{"https://httpbin.org/get", "HTTPBin Get"},
		{"https://httpbin.org/ip", "IP Info"},
		{"https://httpbin.org/user-agent", "User Agent"},
		{"https://httpbin.org/headers", "Headers"},
		{"https://httpbin.org/delay/2", "Delayed Response"}, // This one will take longer
	}

	// Channel to collect responses
	responses := make(chan ApiResponse, len(apis))

	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup
	wg.Add(len(apis))

	// Set a custom timeout for the HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Override the default HTTP client
	http.DefaultClient = client

	fmt.Println("Making concurrent API requests...")

	// Make requests concurrently
	for _, api := range apis {
		go func(url, source string) {
			defer wg.Done()
			responses <- FetchAPI(url, source)
		}(api.URL, api.Source)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(responses)
	}()

	// Process the results
	for resp := range responses {
		if resp.Error != nil {
			fmt.Printf("[%s] Error: %v (took %v)\n", resp.Source, resp.Error, resp.Latency)
		} else {
			fmt.Printf("[%s] Success (took %v)\n", resp.Source, resp.Latency)
			// Print a sample of the data
			for k, v := range resp.Data {
				fmt.Printf("  - %s: %v\n", k, v)
				break // Just show one item for brevity
			}
		}
	}

	fmt.Println("All requests completed!")
}
