package deprecation

// BetaOnly defines replacements.
var BetaOnly = map[string]BetaOnlyInfo{}

// BetaOnlyInfo keeps some deprecation details related to a command.
type BetaOnlyInfo struct {
	Replacement string
}
