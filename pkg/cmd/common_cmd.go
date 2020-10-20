package cmd

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/jmespath/go-jmespath"
)

const twoSpaces = "  "

type CommonCmd struct {
	Query string
}

func (c *CommonCmd) AddOptions(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&c.Query, "query", "q", "", "JMESPath query filter")
}

func (c *CommonCmd) Filter(data interface{}) (string, error) {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return "", errors.Wrap(err, "marshal failed")
	}

	var dataInterface interface{}
	err = json.Unmarshal(marshalledData, &dataInterface)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal failed")
	}

	filtered, err := jmespath.Search(c.Query, dataInterface)

	if err != nil {
		return "", errors.Wrap(err, "filter failed")
	}

	formattedOutput, err := json.MarshalIndent(filtered, "", twoSpaces)
	if err != nil {
		return "", errors.Wrap(err, "output failed")
	}

	return string(formattedOutput), nil
}