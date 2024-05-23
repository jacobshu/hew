package kinsta

type Company struct {
	ID    string `json:"id"`
	Sites []Site
}

type Site struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	DisplayName  string        `json:"display_name"`
	Status       string        `json:"status"`
	SiteLabels   []Label       `json:"site_labels"`
	Environments []Environment `json:"environments"`
  CompanyID    string        `json:"company_id"` // references company.id
	ClientID     string        // references client.id
}

type Clients struct {
	ID   string `json:"client_id"`
	Name string `json:"name"`
}

type Domain struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type_  string `json:"type"`
	SiteID string `json:"site_id"` // references site.id
}

type Environment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsPremium bool   `json:"is_premium"`
	IsBlocked bool   `json:"is_blocked"`
	SiteID    string `json:"site_id"` // references site_id
	// Domains       []Domain `json:"domains"`
	PrimaryDomain Domain `json:"primary_domain"` // references domain_id
	// ssh port and ip are retrieved from /sites/{site_id}/environments
	// every other column can come from /sites/{site_id}
	SSHPort string `json:"ssh_port"`
	SSHIP   string `json:"ssh_ip"`
}

type KinstaCompanies struct {
	ID   string `json:"company_id"`
	Name string `json:"name"`
}

// labels can be retrieved from /sites/{site_id}
type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Plugins struct {
	EnvironmentID   string `json:"environment_id"` // references environment_id
	Name            string `json:"name"`
	Title           string `json:"title"`
	Status          string `json:"status"`
	Version         string `json:"version"`
	UpdateAvailable bool   `json:"update_available"`
	UpdateVersion   string `json:"update_version"`
	UpdateStatus    string `json:"update_status"`
}

type Themes struct {
	EnvironmentID   string `json:"environment_id"` // references environment_id
	Name            string `json:"name"`
	Title           string `json:"title"`
	Status          string `json:"status"`
	Version         string `json:"version"`
	UpdateAvailable bool   `json:"update_available"`
	UpdateVersion   string `json:"update_version"`
	UpdateStatus    string `json:"update_status"`
}
