package kinsta

type Company struct {
	Id    string `json:"id"`
	Sites []Site
}

type Site struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	CompanyId    string `json:"company_id"` // references company.id
	Status       string `json:"status"`
	SiteLabels   []Label
	Environments []Environments `json:"environments"`
	ClientId     string         // references client.id
}

type Clients struct {
	Id   string `json:"client_id"`
	Name string `json:"name"`
}

type Domains struct {
	Domain_id string `json:"domain_id"`
	Name      string `json:"name"`
	Type_     string `json:"type_"`
	Site_id   string `json:"site_id"` // references site.id
}

type Environments struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IsPremium      bool   `json:"is_premium"`
	IsBlocked      bool   `json:"is_blocked"`
	Site_id        string `json:"site_id"`        //{ references: () => Sites.columns.site_id }),
	Primary_domain string `json:"primary_domain"` //{ references: () => Domains.columns.domain_id }),
	// ssh port and ip are retrieved from /sites/{site_id}/environments
	// every other column can come from /sites/{site_id}
	// so may not be worth adding due to the extra API calls
	Ssh_port string `json:"ssh_port"` //{ optional: true }),
	Ssh_ip   string `json:"ssh_ip"`   //{ optional: true }),
}

type KinstaCompanies struct {
	Company_id string `json:"company_id"`
	Name       string `json:"name"`
}

// labels can be retrieved from /sites/{site_id}
type Label struct {
	Label_id string `json:"label_id"`
	Name     string `json:"name"`
}

type Plugins struct {
	Plugin_environment_id string `json:"plugin_environment_id"` //{ references: () => Environments.columns.environment_id }),
	Name                  string `json:"name"`
	Title                 string `json:"title"`
	Status                string `json:"status"`
	Version               string `json:"version"`
	Update_available      bool   `json:"update_available"`
	Update_version        string `json:"update_version"`
	Update_status         string `json:"update_status"`
}

type Themes struct {
	Theme_environment_id string `json:"theme_environment_id"` //{ references: () => Environments.columns.environment_id }),
	Name                 string `json:"name"`
	Title                string `json:"title"`
	Status               string `json:"status"`
	Version              string `json:"version"`
	Update_available     bool   `json:"update_available"`
	Update_version       string `json:"update_version"`
	Update_status        string `json:"update_status"`
}
