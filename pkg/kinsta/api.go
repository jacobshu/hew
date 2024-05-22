package kinsta

import (
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
	fmt.Printf("%v %v %v\n", opts.method, opts.body, opts.queryParams)
	token := os.Getenv("KINSTA_TOKEN")
	if token == "" {
		fmt.Println("KINSTA_TOKEN not set in .env file")
		os.Exit(1)
	}

	client := &http.Client{}

  url, err := buildUrl(opts)
	req, err := http.NewRequest(opts.method, url, opts.body)
	if err != nil {
		return byte[], fmt.Errorf("error creating request: %v", err)
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
  defer resp.Body.Close()
   
	return bodyBytes, nil
}

func buildUrl(req http.Request, opts RequestOpts) (string, error) {
  baseUrl := "https://api.kinsta.com/v2"

  var u string
  switch opts.endpoint {
    case "/sites": // all sites for company
      u = ""
    case "/sites/${siteId}": // single site
      u = ""
    case "/sites/${siteId}/environments": // envs for site
      u = ""
    case "/sites/environments/${envId}/plugins":    // all plugins for site env
      u = ""
    case "/sites/environments/${envId}/plugins/bulk-update":
      u = ""
    case "/sites/environments/${envId}/themes": // all themes for site env
      u = ""
    case "/sites/environments/${envId}/themes/bulk-update":
      u = ""
    case "/sites/environments/${envId}/manual-backups": // POST to create backup
      u = ""
    case "/sites/environments/${envId}/backups": // get backups
      u = ""
    case "/sites/environments/${targetEnvId}/backups/restore": // POST to restore
      u = ""
  }

	url, err := url.JoinPath(baseUrl, opts.endpoint)
  if err != nil {
    return "", fmt.Errorf("error creating url")
  }

  if len(opts.queryParams) > 0 {
		q := req.URL.Query()
		for param, val := range opts.queryParams {
			q.Add(param, val)
		}
	}

	return url, nil 
}

func GetSite(siteId string) (Site, error) {
  site, err := kinsta(RequestOpts{method: "GET", endpoint: "/sites/"+siteId})
  if err != nil {
    fmt.Printf("error getting site %v", err)
  }

   trade := []Trade{}
    err = json.Unmarshal([]byte(body), &trade)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(trade)
  return nil, nil
}

func GetSites(companyId string) ([]Site, error) {
  sites, err := kinsta(RequestOpts{method: "GET", endpoint: "/sites", queryParams: map[string]string{"companyId": companyId}})
	if err != nil {
		fmt.Printf("error getting sites %v", err)
	}

	fmt.Printf("sites: \n%#v", sites)

	return nil, nil
}
