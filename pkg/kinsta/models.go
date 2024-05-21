package kinsta

type Company struct {
	id string
}

type Clients struct {
	client_id string
	name      string
}

type Domains struct {
	domain_id string
	name      string
	type_     string
	site_id   string //{ references: () => Sites.columns.site_id }),
}

type Environments struct {
	environment_id string
	name           string
	is_premium     bool
	is_blocked     bool
	site_id        string //{ references: () => Sites.columns.site_id }),
	primary_domain string //{ references: () => Domains.columns.domain_id }),
	// ssh port and ip are retrieved from /sites/{site_id}/environments
	// every other column can come from /sites/{site_id}
	// so may not be worth adding due to the extra API calls
	ssh_port string //{ optional: true }),
	ssh_ip   string //{ optional: true }),
}

type KinstaCompanies struct {
	company_id string
	name       string
}

// labels can be retrieved from /sites/{site_id}
type Labels struct {
	label_id string
	name     string
	site_id  string //{ references: () => Sites.columns.site_id }),
}

type Plugins struct {
	plugin_environment_id string //{ references: () => Environments.columns.environment_id }),
	name                  string
	title                 string
	status                string
	version               string
	update_available      bool
	update_version        string
	update_status         string
}

type Themes struct {
	theme_environment_id string //{ references: () => Environments.columns.environment_id }),
	name                 string
	title                string
	status               string
	version              string
	update_available     bool
	update_version       string
	update_status        string
}

type Sites struct {
	site_id      string
	name         string
	display_name string
	status       string
	company_id   string //{ references: () => KinstaCompanies.columns.company_id }),
	client_id    string //{ references: () => Clients.columns.client_id, optional: true }),
}
