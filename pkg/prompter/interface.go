package prompter

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Prompter

type Prompter interface {
	SelectFromOptions(question string, options []string) (string, error)
	SelectMultipleFromOptions(question string, options []string) ([]string, error)
	SelectFromOptionsWithDefault(question, defaultValue string, options []string) (string, error)
}
