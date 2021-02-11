package mocks

import (
	"github.com/plumming/dx/pkg/util"
)

// MockCommandRunner is the mock command.
type MockCommandRunner struct {
	RunWithoutRetryFunc func(c *util.Command) (string, error)
	Commands            []string
}

// RunWithoutRetry is the mock command's `RunWithoutRetry` func.
func (m *MockCommandRunner) RunWithoutRetry(c *util.Command) (string, error) {
	m.Commands = append(m.Commands, c.String())
	return m.RunWithoutRetryFunc(c)
}
