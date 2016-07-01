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
	return do(c, "GET", path, nil, dst)
}

func post(c *cli.Context, path string, src, dst interface{}) error {
	return do(c, "POST", path, src, dst)
}

func patch(c *cli.Context, path string, src, dst interface{}) error {
	return do(c, "PATCH", path, src, dst)
}

func del(c *cli.Context, path string) error {
	return do(c, "DELETE", path, nil, nil)
}

func do(c *cli.Context, method, path string, src interface{}, dst interface{}) error {

	app := c.GlobalString("app")
	mount := c.GlobalString("mount")
	secret := c.GlobalString("secret")

	body := bytes.NewBuffer(nil)

	if src != nil {
		jsonData, err := json.Marshal(src)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonData)
	}

	r, err := http.NewRequest(method, app+mount+path, body)
	if err != nil {
		return err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", secret)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	} else if data, err := readResponse(resp); err != nil {
		return err
	} else if dst == nil {
		return nil
	} else {
		return json.Unmarshal(data, dst)
	}

}

func readResponse(resp *http.Response) ([]byte, error) {

	// slurp response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode > 399 {
		x := struct {
			Message string `json:"message"`
		}{}
		json.Unmarshal(data, &x)
		return nil, cli.NewExitError(fmt.Sprintf("HTTP status %d: %s", resp.StatusCode, x.Message), 1)
	}

	return data, nil

}
