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
			Usage:  "Generate a new auth secret and load data from seed.json",
			Action: cmd.Initialize,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "seed",
					Value: pwd + "/seed.json",
					Usage: "Path to seed.json",
				},
			},
		},
		{
			Name:   "load",
			Usage:  "Import JSON from stdin to an application",
			Action: cmd.LoadEntities,
		},
		{
			Name:  "config",
			Usage: "Get or set a configuration variable on an application",
			Subcommands: []cli.Command{
				{
					Name:      "get",
					Usage:     "Get a configuration variable from the app (or the whole config if no one variable is specified)",
					ArgsUsage: "[key]",
					Action:    cmd.GetConfig,
				},
				{
					Name:      "set",
					Usage:     "Set an environment variable on the app.",
					ArgsUsage: "key value [key] [value] ...",
					Action:    cmd.SetConfig,
					Flags:     []cli.Flag{},
				},
			},
		},
		{
			Name:  "account",
			Usage: "Modify accounts associated with the application",
			Subcommands: []cli.Command{
				{
					Name:      "add",
					Usage:     "Add account",
					ArgsUsage: "email password",
					Action:    cmd.AddAccount,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "roles",
							Usage: "Attach comma-separated `ROLES_LIST` to this new account",
						},
					},
				},
				{
					Name:      "remove",
					Usage:     "Remove account",
					ArgsUsage: "email",
					Action:    cmd.RemoveAccount,
				},
				{
					Name:      "get",
					Usage:     "Show account",
					ArgsUsage: "email",
					Action:    cmd.GetAccount,
				},
				{
					Name:  "roles",
					Usage: "Actions related to user roles",
					Subcommands: []cli.Command{
						{
							Name:      "add",
							Usage:     "Add roles to account",
							ArgsUsage: "role [role...]",
							Action:    cmd.AddAccountRole,
						},
						{
							Name:      "remove",
							Usage:     "Remove roles from account",
							ArgsUsage: "role [role...]",
							Action:    cmd.RemoveAccountRole,
						},
					},
				},
				{
					Name:      "password",
					Usage:     "Change password on account",
					ArgsUsage: "new_password",
					Action:    cmd.ChangeAccountPassword,
				},
				{
					Name:      "change-email",
					Usage:     "Change email address on account",
					ArgsUsage: "new_address",
					Action:    cmd.ChangeAccountEmail,
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
