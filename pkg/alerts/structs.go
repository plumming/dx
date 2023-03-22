package alerts

import "time"

type Repository struct {
	NameWithOwner       string              `json:"nameWithOwner"`
	VulnerabilityAlerts VulnerabilityAlerts `json:"vulnerabilityAlerts"`
}

type VulnerabilityAlerts struct {
	Nodes []SecurityAdvisoryNode `json:"nodes"`
}

type SecurityAdvisoryNode struct {
	SecurityAdvisory SecurityAdvisory `json:"securityAdvisory"`
	State            string           `json:"state"`
	CreatedAt        *time.Time       `json:"createdAt"`
	DependabotUpdate DependabotUpdate `json:"dependabotUpdate"`
}

type SecurityAdvisory struct {
	GhsaID      string       `json:"ghsaId"`
	Severity    string       `json:"severity"`
	Summary     string       `json:"summary"`
	CVSS        CVSS         `json:"cvss"`
	Identifiers []Identifier `json:"identifiers"`
}

type DependabotUpdate struct {
	PullRequest *PullRequest `json:"pullRequest"`
}

type PullRequest struct {
	Number int `json:"number"`
}

type CVSS struct {
	VectorString string  `json:"vectorString"`
	Score        float32 `json:"score"`
}

type Identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
