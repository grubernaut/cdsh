package main

import "gopkg.in/urfave/cli.v2"

func setFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Value:   false,
			Aliases: []string{"v"},
			Usage:   "Verbose output",
		},
		&cli.StringFlag{
			Name:    "server",
			Value:   "",
			Aliases: []string{"s"},
			Usage:   "Consul Server Address",
		},
		&cli.StringFlag{
			Name:    "service",
			Value:   "",
			Aliases: []string{"S"},
			Usage:   "Target Consul Service",
		},
		&cli.StringFlag{
			Name:    "user",
			Value:   "",
			Aliases: []string{"u"},
			Usage:   "Remote User to connect as",
		},
	}
}
