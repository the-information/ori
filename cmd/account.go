// Package cmd supports the ori command-line utility.
package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	_ "strings"
)

func AddAccount(c *cli.Context) error {

	if c.NArg() != 2 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}

	requestBody := map[string]interface{}{
		"Email":    c.Args().Get(0),
		"Password": c.Args().Get(1),
	}

	if err := post(c, "accounts", requestBody, nil); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	return nil

}

func RemoveAccount(c *cli.Context) error {

	if c.NArg() != 1 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}

	key := base64.RawURLEncoding.EncodeToString([]byte(c.Args().Get(0)))

	if err := del(c, "accounts/"+key); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	return nil

}

func GetAccount(c *cli.Context) error {

	if c.NArg() != 1 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}

	account := json.RawMessage{}
	formattedAccount := bytes.NewBuffer(nil)

	key := base64.RawURLEncoding.EncodeToString([]byte(c.Args().Get(0)))

	if err := get(c, "accounts/"+key, &account); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	// format the JSON
	json.Indent(formattedAccount, account, "", "  ")

	fmt.Println(formattedAccount)

	return nil
}

func GetJwt(c *cli.Context) error {

	if c.NArg() != 1 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}
	var jwt string

	key := base64.RawURLEncoding.EncodeToString([]byte(c.Args().Get(0)))

	if err := get(c, "accounts/"+key+"/jwt", &jwt); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	fmt.Println(jwt)

	return nil

}

func ChangeAccountEmail(c *cli.Context) error {

	if c.NArg() != 2 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}

	key := base64.RawURLEncoding.EncodeToString([]byte(c.Args().Get(0)))
	body := map[string]string{
		"email": c.Args().Get(1),
	}

	if err := patch(c, "accounts/"+key, &body, nil); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	return nil

}

func ChangeAccountPassword(c *cli.Context) error {

	if c.NArg() != 2 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}

	key := base64.RawURLEncoding.EncodeToString([]byte(c.Args().Get(0)))
	body := c.Args().Get(1)

	if err := post(c, "accounts/"+key+"/password", &body, nil); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	return nil
}

type rolesRequest struct {
	Roles []string `json:"roles"`
}

func AddAccountRole(c *cli.Context) error {

	if c.NArg() < 2 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}

	key := base64.RawURLEncoding.EncodeToString([]byte(c.Args().Get(0)))
	newRoles := c.Args().Tail()

	account := rolesRequest{}

	if err := get(c, "accounts/"+key, &account); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	roleMap := map[string]bool{}

	for _, role := range account.Roles {
		roleMap[role] = true
	}

	for _, newRole := range newRoles {
		if _, ok := roleMap[newRole]; !ok {
			account.Roles = append(account.Roles, newRole)
		}
	}

	if err := patch(c, "accounts/"+key, &account, nil); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	return nil

}

func RemoveAccountRole(c *cli.Context) error {

	if c.NArg() < 2 {
		return cli.NewExitError("Too many or not enough arguments specified", 1)
	}

	key := base64.RawURLEncoding.EncodeToString([]byte(c.Args().Get(0)))
	removeRoles := c.Args().Tail()

	account := rolesRequest{}

	if err := get(c, "accounts/"+key, &account); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	roleMap := map[string]bool{}

	for _, role := range account.Roles {
		roleMap[role] = true
	}

	for _, removeRole := range removeRoles {
		delete(roleMap, removeRole)
	}

	account.Roles = []string{}

	for role, _ := range roleMap {
		account.Roles = append(account.Roles, role)
	}

	if err := patch(c, "accounts/"+key, &account, nil); err != nil {
		return cli.NewExitError("Server error: "+err.Error(), 1)
	}

	return nil
}
