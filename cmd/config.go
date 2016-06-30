package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
)

func SetConfig(c *cli.Context) error {

	if c.NArg() < 2 {
		return cli.NewExitError("Must supply a key and a value to set it to", 1)
	} else if c.NArg()%2 != 0 {
		return cli.NewExitError("Must supply a list of key-value pairs", 1)
	}

	patchData := map[string]interface{}{}
	conf := json.RawMessage{}

	for i := 0; i < c.NArg(); i += 2 {

		key := c.Args().Get(i)
		value := []byte(c.Args().Get(i + 1))

		var unmarshaledValue interface{}

		if err := json.Unmarshal(value, &unmarshaledValue); err != nil {
			unmarshaledValue = c.Args().Get(i + 1)
		}
		patchData[key] = unmarshaledValue

	}

	if err := patch(c, "config", &patchData, &conf); err != nil {
		return cli.NewExitError("Error from server: "+err.Error(), 1)
	}

	indentedJSON := bytes.NewBuffer(nil)
	if err := json.Indent(indentedJSON, []byte(conf), "", "  "); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(indentedJSON)

	return nil

}

func GetConfig(c *cli.Context) error {

	if c.Args().First() == "" {
		// get the full configuration.
		return getFullConfig(c)
	} else {
		return getSingleConfig(c)
	}

}

func getFullConfig(c *cli.Context) error {

	conf := json.RawMessage{}
	if err := get(c, "config", &conf); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	result := bytes.NewBuffer(nil)
	if err := json.Indent(result, []byte(conf), "", "  "); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(result)

	return nil

}

func getSingleConfig(c *cli.Context) error {

	requestedVar := c.Args().First()

	conf := map[string]json.RawMessage{}
	if err := get(c, "config", &conf); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	val, ok := conf[requestedVar]
	if !ok {
		return cli.NewExitError("No such configuration variable: "+requestedVar, 1)
	}

	fmt.Println(string(val))
	return nil

}
