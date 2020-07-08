package deprecation

var DeprecatedCommands = map[string]Info{}

// Info keeps some deprecation details related to a command.
type Info struct {
	Replacement string
	Date        string
	Info        string
}
