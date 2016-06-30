// Package cmd supports the ori command-line utility.
package cmd

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
)

func SetConfig(c *cli.Context) error {
	return nil
}

func GetConfig(c *cli.Context) error {

	mount := c.GlobalString("mount")
	app := c.GlobalString("app")
	secret := c.GlobalString("secret")

	r, err := http.NewRequest("GET", app+mount+"config", nil)
	if err != nil {
		return eri
		r
	}

	r.Header.Set("Authorization", secret)
	resp, err := http.DefaultClient.Do(r)

	// slurp response
	data, err := string(ioutil.ReadAll(resp.Body)); if err != nil {
		return err
	} else {
		resp.Body.Close()
	}

	if err != nil {
		return err
	}

	switch resp.StatusCode {
		default:
			return cli.NewExitError("Unknown error", 1)
	case http.StatusForbidden:
		case http.StatusOK:

	}

	}

	return nil
}

func AddUser(c *cli.Context) error {
	return nil
}

func RemoveUser(c *cli.Context) error {
	return nil
}

func ChangeUserEmail(c *cli.Context) error {
	return nil
}

func ChangeUserPassword(c *cli.Context) error {
	return nil
}

func AddUserRole(c *cli.Context) error {
	return nil
}

func RemoveUserRole(c *cli.Context) error {
	return nil
}
