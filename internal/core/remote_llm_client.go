package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// RemoteLLMClient implements LLMClient interface for connecting to remote LLM APIs
// It supports configuration through environment variables
type RemoteLLMClient struct {
	apiURL      string
	apiKey      string
	model       string
	temperature float64
	timeout     time.Duration
	httpClient  *http.Client
}

// NewRemoteLLMClient creates a new RemoteLLMClient using environment variables
// Environment variables:
//   LLM_API_URL - The API endpoint URL
//   LLM_API_KEY - API key for authentication
//   LLM_MODEL - Model name to use (default: "mistral-tiny")
//   LLM_TEMPERATURE - Temperature for response randomness (default: 0.7)
//   LLM_TIMEOUT - Request timeout in seconds (default: 30)
//
// Supported APIs:
//   - Mistral AI: https://api.mistral.ai/v1/completions
//   - OpenAI: https://api.openai.com/v1/chat/completions
//   - Any OpenAI-compatible API endpoint
func NewRemoteLLMClient() (*RemoteLLMClient, error) {
	apiURL := os.Getenv("LLM_API_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("LLM_API_URL environment variable is required")
	}

	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY environment variable is required")
	}

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "mistral-tiny" // Default model
	}

	temperature := 0.7 // Default temperature
	if tempStr := os.Getenv("LLM_TEMPERATURE"); tempStr != "" {
		if _, err := fmt.Sscanf(tempStr, "%f", &temperature); err != nil {
			return nil, fmt.Errorf("invalid LLM_TEMPERATURE: %v", err)
		}
	}

	timeoutSeconds := 30 // Default timeout
	if timeoutStr := os.Getenv("LLM_TIMEOUT"); timeoutStr != "" {
		if _, err := fmt.Sscanf(timeoutStr, "%d", &timeoutSeconds); err != nil {
			return nil, fmt.Errorf("invalid LLM_TIMEOUT: %v", err)
		}
	}

	return &RemoteLLMClient{
		apiURL:      apiURL,
		apiKey:      apiKey,
		model:       model,
		temperature: temperature,
		timeout:     time.Duration(timeoutSeconds) * time.Second,
		httpClient:  &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
	}, nil
}

// Complete implements the LLMClient interface
// It sends a request to the remote LLM API and returns the completion
func (c *RemoteLLMClient) Complete(prompt, text string) (string, error) {
	// Construct the full prompt by combining the rule prompt with the text
	fullPrompt := fmt.Sprintf("%s\n\nText: %s\n\nAnswer:", prompt, text)

	// Prepare the request payload
	// Note: This supports both Mistral-style completions and OpenAI-style chat completions
	requestPayload := map[string]interface{}{
		"model":       c.model,
		"messages":    []map[string]string{
			{"role": "user", "content": fullPrompt},
		},
		"temperature": c.temperature,
	}

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request payload: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract the completion text from the response
	// The structure depends on the API provider
	if choices, ok := responseMap["choices"].([]interface{}); ok && len(choices) > 0 {
		if firstChoice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := firstChoice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					// Clean up the response by trimming whitespace and removing any trailing punctuation
					answer := strings.TrimSpace(content)
					// Remove any trailing periods, exclamation marks, or question marks
					answer = strings.TrimRight(answer, ".!?")
					return answer, nil
				}
			}
			// Fallback for some API formats
			if content, ok := firstChoice["text"].(string); ok {
				answer := strings.TrimSpace(content)
				answer = strings.TrimRight(answer, ".!?")
				return answer, nil
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from LLM API")
}

// HealthCheck verifies the connection to the LLM API
func (c *RemoteLLMClient) HealthCheck() error {
	// Create a simple test request
	testPayload := map[string]interface{}{
		"model":    c.model,
		"messages": []map[string]string{
			{"role": "user", "content": "Test connection"},
		},
	}

	payloadBytes, err := json.Marshal(testPayload)
	if err != nil {
		return fmt.Errorf("failed to create test payload: %v", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create test request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}