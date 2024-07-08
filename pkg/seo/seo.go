package seo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"time"
)

func main() {	
  // results found on 
	type OnPageOptions struct {
		Target                  string `json:"target"`
		MaxCrawledPages         int    `json:"max_crawled_pages"`
		EnableContentParsing    bool   `json:"enable_content_parsing"`
		BrowserPreset           string `json:"browser_preset"`
		EnableJavascript        bool   `json:"enable_javascript"`
		LoadResources           bool   `json:"load_resources"`
		EnableBrowserRendering  bool   `json:"enable_browser_rendering"`
		ValidateMicromarkup     bool   `json:"validate_micromarkup"`
		CalculateKeywordDensity bool   `json:"calculate_keyword_density"`

	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	postBody, _ := json.Marshal(map[string]string{
		"name":  "Toby",
		"email": "Toby@example.com",
	})
	responseBody := bytes.NewBuffer(postBody)

	req, err := createRequest(http.MethodPost, "/on_page/task_post", postBody)
	if err != nil {
		log.Fatalf("an error occurred while creating the request %v", err)
	}
	resp, err := client.Do(req)
	resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}

func createRequest(method string, endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, path.Join("https://api.dataforseo.com/v3", endpoint), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("login", "password")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
