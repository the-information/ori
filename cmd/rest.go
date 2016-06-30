package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
)

func get(c *cli.Context, path string, dst interface{}) error {

	mount := c.GlobalString("mount")
	app := c.GlobalString("app")
	secret := c.GlobalString("secret")

	r, err := http.NewRequest("GET", app+mount+path, nil)
	if err != nil {
		return err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", secret)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	data, err := readResponse(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dst)

}

func patch(c *cli.Context, path string, src, dst interface{}) error {

	mount := c.GlobalString("mount")
	app := c.GlobalString("app")
	secret := c.GlobalString("secret")

	body, err := json.Marshal(src)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("PATCH", app+mount+path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", secret)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	data, err := readResponse(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dst)

}

func readResponse(resp *http.Response) ([]byte, error) {

	// slurp response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode > 299 {
		x := struct {
			Message string `json:"message"`
		}{}
		json.Unmarshal(data, &x)
		return nil, cli.NewExitError(fmt.Sprintf("HTTP status %d: %s", resp.StatusCode, x.Message), 1)
	}

	return data, nil

}
