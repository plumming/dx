package deprecation

// AlphaOnly defines replacements.
var AlphaOnly = map[string]AlphaOnlyInfo{}

// AlphaOnlyInfo keeps some deprecation details related to a command.
type AlphaOnlyInfo struct {
	Replacement string
}
