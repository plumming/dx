package securityconfig

type Repository struct {
	NameWithOwner                 string           `json:"nameWithOwner"`
	URL                           string           `json:"url"`
	HasVulnerabilityAlertsEnabled bool             `json:"hasVulnerabilityAlertsEnabled"`
	IsSecurityPolicyEnabled       bool             `json:"isSecurityPolicyEnabled"`
	DefaultBranchRef              DefaultBranchRef `json:"defaultBranchRef"`
}

type DefaultBranchRef struct {
	Name                 string                `json:"name"`
	BranchProtectionRule *BranchProtectionRule `json:"branchProtectionRule"`
}

type BranchProtectionRule struct {
	Pattern              string `json:"pattern"`
	RequiresStatusChecks bool   `json:"requiresStatusChecks"`
	RestrictsPushes      bool   `json:"restrictsPushes"`
}

// sort functions

type ByNameAndOwner []Repository

func (p ByNameAndOwner) Len() int {
	return len(p)
}

func (p ByNameAndOwner) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByNameAndOwner) Less(i, j int) bool {
	return p[i].NameWithOwner < p[j].NameWithOwner
}
