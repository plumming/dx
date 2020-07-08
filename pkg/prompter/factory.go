package prompter

import "gopkg.in/AlecAivazis/survey.v1"

type factory struct {
}

func (f *factory) SelectFromOptions(question string, options []string) (string, error) {
	var result string
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}

	err := survey.AskOne(prompt, &result, survey.Required)
	if err != nil {
		return result, err
	}

	return result, nil
}

func NewPrompter() Prompter {
	return &factory{}
}
