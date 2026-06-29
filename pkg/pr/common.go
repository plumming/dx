package pr

import "time"

type Comments struct {
	TotalCount int `json:"totalCount"`
}

type Repository struct {
	NameWithOwner string `json:"nameWithOwner"`
}

type Author struct {
	Login string `json:"login"`
}

type Labels struct {
	Nodes []Label `json:"nodes"`
}

type Label struct {
	Name string `json:"name"`
}

type Commits struct {
	Nodes []CommitEntry `json:"nodes"`
}

type CommitEntry struct {
	Commit Commit `json:"commit"`
}

type Commit struct {
	StatusCheckRollup StatusCheckRollup `json:"statusCheckRollup"`
}

type StatusCheckRollup struct {
	State    string        `json:"state"`
	Contexts StatusContext `json:"contexts"`
}

type StatusContext struct {
	Nodes []Context `json:"nodes"`
}

type Context struct {
	State       string    `json:"state"`
	Description string    `json:"description"`
	Context     string    `json:"context"`
	Conclusion  string    `json:"conclusion"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	StartedAt   time.Time `json:"startedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

// timestamp returns the time the context was last run, preferring a CheckRun's
// startedAt and falling back to a StatusContext's createdAt.
func (c Context) timestamp() time.Time {
	if !c.StartedAt.IsZero() {
		return c.StartedAt
	}
	return c.CreatedAt
}
