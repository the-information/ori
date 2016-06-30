package main

import (
	"github.com/the-information/ori/cmd"
	"github.com/urfave/cli"
	"os"
)

const (
	VERSION = "1.0.0"
)

var pwd string
var flags []cli.Flag
var commands []cli.Command

func init() {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Value: pwd + "/ori.yaml",
			Usage: "Path to ori.yaml",
		},
		cli.StringFlag{
			Name:   "app",
			Value:  "http://localhost:8080",
			Usage:  "URL to address App Engine app at",
			EnvVar: "ORI_APP_NAME",
		},
		cli.StringFlag{
			Name:  "mount",
			Value: "/_ori/",
			Usage: "Mount point of admin routes in app",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "Auth secret for application",
			EnvVar: "ORI_APP_SECRET",
		},
	}

	commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Set initial configuration from ori.yaml to an application",
			Action: func(c *cli.Context) error { return nil },
		},
		{
			Name:  "config",
			Usage: "Get or set a configuration variable on an application",
			Subcommands: []cli.Command{
				{
					Name:   "get",
					Usage:  "Get a configuration variable from the app",
					Action: cmd.GetConfig,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "int",
							Usage: "Treat the variable as an integer",
						},
						cli.BoolFlag{
							Name:  "float",
							Usage: "Treat the variable as a floating-point number",
						},
						cli.BoolFlag{
							Name:  "bool",
							Usage: "Treat the variable as a boolean value ('true' and 'false')",
						},
					},
				},
				{
					Name:   "set",
					Usage:  "Set an environment variable on the app. Use syntax NAME=VALUE",
					Action: cmd.SetConfig,
					Flags:  []cli.Flag{},
				},
			},
		},
		{
			Name:  "account",
			Usage: "Modify accounts associated with the application",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "Add account",
				},
				{
					Name:  "remove",
					Usage: "Remove account",
				},
				{
					Name:  "password",
					Usage: "Change password on account",
				},
				{
					Name:  "change-email",
					Usage: "Change email address on account",
				},
			},
		},
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "ori"
	app.Usage = "Develop REST/JSON APIs on Google App Engine"
	app.Version = VERSION
	app.EnableBashCompletion = true
	app.Commands = commands
	app.Flags = flags

	app.Run(os.Args)

}
