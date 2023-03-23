package securityconfig

type Repository struct {
	NameWithOwner                 string `json:"nameWithOwner"`
	URL                           string `json:"url"`
	HasVulnerabilityAlertsEnabled bool   `json:"hasVulnerabilityAlertsEnabled"`
}
