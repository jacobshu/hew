package kinsta

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Kinsta(method string, url string, body io.Reader) (string, error) {
  baseUrl := "https://api.kinsta.com/v2"
  token := os.Getenv("KINSTA_TOKEN")
	if token == "" {
		fmt.Println("KINSTA_TOKEN not set in .env file")
		os.Exit(1)
	}

	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest(method, baseUrl+url, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(bodyBytes), nil
}

