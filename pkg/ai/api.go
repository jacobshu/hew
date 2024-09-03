package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiURL = "https://api.anthropic.com/v1/messages" // This is a placeholder URL
)

type Client interface {
	SendMessage(message string) (string, error)
}

type APIClient struct {
	apiKey string
	client *http.Client
}

func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type requestBody struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type responseBody struct {
	Content string `json:"content"`
}

func (c *APIClient) SendMessage(message string) (string, error) {
	reqBody := requestBody{
		Model: "claude-3-opus-20240229",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var respBody responseBody
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %w", err)
	}

	return respBody.Content, nil
}
