package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "print the version",
		Aliases: []string{"V"},
		Value:   false,
	}

	app := &cli.App{
		Name:  "cdsh",
		Usage: "Consul DSH",
		Flags: setFlags(),
		Action: func(c *cli.Context) error {
			return run(c)
		},
	}
	app.Run(os.Args)
}
