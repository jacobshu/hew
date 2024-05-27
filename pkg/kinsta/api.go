package kinsta

import (
	"bytes"
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
	body        interface{}
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

	var body []byte
	if opts.body != nil {
		b, err := json.Marshal(opts.body)
		if err != nil {
			return nil, fmt.Errorf("error encoding body:\n%#v", opts.body)
		}
		body = b
	}

	req, err := http.NewRequest(opts.method, url, bytes.NewReader(body))
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
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
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

func GetPlugins(envID string) ([]Plugin, error) {
	url := "/sites/environments/" + envID + "/plugins"
	pluginBody, err := kinsta(RequestOpts{method: "GET", endpoint: url})
	if err != nil {
		return []Plugin{}, err
	}

	plugins := struct {
		Environment struct {
			Container struct {
				WPPlugins struct {
					Data []Plugin `json:"data"`
				} `json:"wp_plugins"`
			} `json:"container_info"`
		} `json:"environment"`
	}{}

	err = json.Unmarshal([]byte(pluginBody), &plugins)
	if err != nil {
		return []Plugin{}, err
	}

	return plugins.Environment.Container.WPPlugins.Data, nil
}

func GetThemes(envID string) ([]Theme, error) {
	url := "/sites/environments/" + envID + "/themes"
	themeBody, err := kinsta(RequestOpts{method: "GET", endpoint: url})
	if err != nil {
		return []Theme{}, err
	}

	themes := struct {
		Environment struct {
			Container struct {
				WPThemes struct {
					Data []Theme `json:"data"`
				} `json:"wp_themes"`
			} `json:"container_info"`
		} `json:"environment"`
	}{}

	err = json.Unmarshal([]byte(themeBody), &themes)
	if err != nil {
		return []Theme{}, err
	}

	return themes.Environment.Container.WPThemes.Data, nil
}

func GetBackups(envID string) ([]Backup, error) {
	url := "/sites/environments/" + envID + "/backups"
	backupBody, err := kinsta(RequestOpts{method: "GET", endpoint: url})
	if err != nil {
		return []Backup{}, err
	}

	backups := struct {
		Environment struct {
			Backups []Backup `json:"backups"`
		} `json:"environment"`
	}{}

	err = json.Unmarshal(backupBody, &backups)
	if err != nil {
		return []Backup{}, err
	}

	return backups.Environment.Backups, nil
}

func CreateManualBackup(envID string, note string) (string, error) {
	url := "sites/environments/" + envID + "/manual-backups"

	tag := struct {
		Tag string `json:"tag"`
	}{
		Tag: note,
	}

	responseBody, err := kinsta(RequestOpts{method: "POST", endpoint: url, body: tag})
	if err != nil {
		return "", err
	}

	operation := struct {
		OperationID string `json:"operation_id"`
	}{}
	err = json.Unmarshal(responseBody, &operation)
	if err != nil {
		return "", err
	}

	return operation.OperationID, nil
}

func IsOperationFinished(operationID string) (bool, error) {
	url := "/operations/" + operationID
	responseBody, err := kinsta(RequestOpts{method: "GET", endpoint: url})
	if err != nil {
		return false, err
	}

	status := struct {
		Code int `json:"status"`
	}{}
	err = json.Unmarshal(responseBody, &status)
	if err != nil {
		return false, nil
	}

	if status.Code == 200 {
		return true, nil
	}
	return false, nil
}
