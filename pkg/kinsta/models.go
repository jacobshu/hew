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
	SiteLabels   []Label       `json:"site_labels,omitempty"`
	Environments []Environment `json:"environments,omitempty"`
	CompanyID    string        `json:"-"` // references company.id
	ClientID     string        `json:"-"` // references client.id
}

type Clients struct {
	ID   string `json:"client_id"`
	Name string `json:"name"`
}

type Domain struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type_  string `json:"type"`
	SiteID string `json:"-"` // references site.id
}

type Environment struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	IsPremium     bool          `json:"is_premium"`
	IsBlocked     bool          `json:"is_blocked"`
	SiteID        string        `json:"-"`             // references site_id
	PrimaryDomain Domain        `json:"primaryDomain"` // references domain_id
	SSH           SSHConnection `json:"ssh_connection,omitempty"`
	// Domains       []Domain `json:"domains"`
	// ssh connection is  retrieved from /sites/{site_id}/environments
	// every other field can come from /sites/{site_id}
}

type SSHConnection struct {
	SSHPort string `json:"ssh_port,omitempty"`
	SSHIP   struct {
		ExternalIP string `json:"external_ip,omitempty"`
	} `json:"ssh_ip,omitempty"`
}

type KinstaCompanies struct {
	ID   string `json:"company_id"`
	Name string `json:"name"`
}

// labels can be retrieved from /sites/{site_id}
type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Plugin struct {
	EnvironmentID   string `json:"-"` // references environment_id
	Name            string `json:"name"`
	Title           string `json:"title"`
	Status          string `json:"status"`
	Version         string `json:"version"`
	UpdateAvailable bool   `json:"update_available"`
	UpdateVersion   string `json:"update_version"`
	UpdateStatus    string `json:"update_status"`
}

type Theme struct {
	EnvironmentID   string `json:"-"` // references environment_id
	Name            string `json:"name"`
	Title           string `json:"title"`
	Status          string `json:"status"`
	Version         string `json:"version"`
	UpdateAvailable bool   `json:"update_available"`
	UpdateVersion   string `json:"update_version"`
	UpdateStatus    string `json:"update_status"`
}
