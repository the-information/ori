package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

func Initialize(c *cli.Context) error {

	// read seed.json and send the data to the db
	reader, err := os.Open(c.String("seed"))
	if err != nil {
		return cli.NewExitError("Error opening seed.json: "+err.Error(), 1)
	}

	seedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return cli.NewExitError("Error reading seed.json: "+err.Error(), 1)
	}

	var src = json.RawMessage(seedData)

	// generate a random string for the secret
	randomBytes := make([]byte, 48)

	_, err = rand.Read(randomBytes)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	newAppSecret := base64.RawURLEncoding.EncodeToString(randomBytes)
	conf := struct{ AuthSecret string }{newAppSecret}

	if err := patch(c, "config", &conf, nil); err != nil {
		return cli.NewExitError("Error from server: "+err.Error(), 1)
	}

	fmt.Println("Add this secret to your environment -- it won't be available again.")
	fmt.Printf("export ORI_APP_SECRET=%s\n", newAppSecret)
	c.GlobalSet("secret", newAppSecret)

	// send to admin
	if err = post(c, "load", &src, nil); err != nil {
		return err
	}
	return nil

}
