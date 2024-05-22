package kinsta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type RequestOpts struct {
	endpoint    string
	method      string // TODO: enum refactor
	body        io.Reader
	queryParams map[string]string
}

func kinsta(opts RequestOpts) ([]byte, error) {
	baseUrl := "https://api.kinsta.com/v2"
	token := os.Getenv("KINSTA_TOKEN")
	if token == "" {
		fmt.Println("KINSTA_TOKEN not set in .env file")
		os.Exit(1)
	}

	url, err := url.JoinPath(baseUrl, opts.endpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating url")
	}

	client := &http.Client{}

	req, err := http.NewRequest(opts.method, url, opts.body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if len(opts.queryParams) > 0 {
		q := req.URL.Query()
		for param, val := range opts.queryParams {
			q.Add(param, val)
		}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	return bodyBytes, nil
}

func GetSite(siteId string) (Site, error) {
  type GetSiteResponse struct {
    site Site
  }

	siteBody, err := kinsta(RequestOpts{method: "GET", endpoint: "/sites/" + siteId})
	if err != nil {
		fmt.Printf("error getting site %v\n", err)
	}

  fmt.Printf("stringy response:\n%v\n", string(siteBody))

	site := GetSiteResponse{}
	err = json.Unmarshal([]byte(siteBody), &site)
	if err != nil {
    fmt.Printf("error unmarshalling: %v\n", err)
	}
	fmt.Printf("site:\n%#v\n", site)
	return site.site, nil
}

func GetSites(companyId string) ([]Site, error) {
	sites, err := kinsta(RequestOpts{method: "GET", endpoint: "/sites", queryParams: map[string]string{"companyId": companyId}})
	if err != nil {
		fmt.Printf("error getting sites %v", err)
	}

	fmt.Printf("sites: \n%#v", sites)

	return nil, nil
}
