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
		return nil, fmt.Errorf("error creating request: %v\n", err)
	}

	if len(opts.queryParams) > 0 {
		q := req.URL.Query()
		for param, val := range opts.queryParams {
			q.Add(param, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received status code %v\n", string(b))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v\n", err)
	}
	defer resp.Body.Close()

	return bodyBytes, nil
}

func GetSite(siteID string) (Site, error) {
	siteBody, err := kinsta(RequestOpts{method: "GET", endpoint: "/sites/" + siteID})
	if err != nil {
		return Site{}, err
	}

	site := struct {
		Site Site `json:"site"`
	}{}

	err = json.Unmarshal([]byte(siteBody), &site)
	if err != nil {
		return Site{}, err
	}

	return site.Site, nil
}

func GetSites(companyID string) ([]Site, error) {
	sitesBody, err := kinsta(RequestOpts{method: "GET", endpoint: "/sites", queryParams: map[string]string{"company": companyID}})
	if err != nil {
		return []Site{}, err
	}

	sites := struct {
		Company struct {
			Sites []Site `json:"sites"`
		} `json:"company"`
	}{}

	err = json.Unmarshal([]byte(sitesBody), &sites)
	if err != nil {
		return []Site{}, err
	}

	return sites.Company.Sites, nil
}

func GetEnvironments(siteID string) ([]Environment, error) {
	url := "/sites/" + siteID + "/environments"
	envBody, err := kinsta(RequestOpts{method: "GET", endpoint: url})
	if err != nil {
		return []Environment{}, err
	}

	envs := struct {
		Site struct {
			Environments []Environment `json:"environments"`
		} `json:"site"`
	}{}

	err = json.Unmarshal([]byte(envBody), &envs)
	if err != nil {
		return []Environment{}, err
	}

	return envs.Site.Environments, nil
}
