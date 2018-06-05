package main

import (
	"os"
	"path/filepath"

	"github.com/rancher/rio/cli/cmd/create"
	"github.com/rancher/rio/cli/cmd/edit"
	"github.com/rancher/rio/cli/cmd/export"
	"github.com/rancher/rio/cli/cmd/server"
	"github.com/rancher/rio/cli/cmd/up"
	"github.com/rancher/rio/cli/pkg/builder"
	"github.com/rancher/rio/cli/pkg/waiter"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	appName = filepath.Base(os.Args[0])
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "wait,w",
			Usage: "Wait for resource to reach resting state",
		},
		cli.IntFlag{
			Name:  "wait-timeout",
			Usage: "Timeout in seconds to wait",
			Value: 600,
		},
		cli.StringFlag{
			Name:  "wait-state",
			Usage: "State to wait for (active, healthy, etc)",
		},
		cli.StringFlag{
			Name:   "url",
			Usage:  "Specify the Rancher API endpoint URL",
			EnvVar: "RODEO_URL",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "Specify Rancher API access key",
			EnvVar: "RODEO_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "Specify Rancher API secret key",
			EnvVar: "RODEO_SECRET_KEY",
		},
	}

	app.Commands = []cli.Command{
		builder.Command(&server.Server{},
			"Run management server",
			appName+" create [OPTIONS]",
			""),
		builder.Command(&create.Create{},
			"Create a new service",
			appName+" create [OPTIONS] IMAGE [COMMAND] [ARG...]",
			""),
		builder.Command(&edit.Edit{},
			"Edit a resource",
			appName+" edit [TYPE] ID_OR_NAME",
			""),
		builder.Command(&export.Export{},
			"Export a stack",
			appName+" export STACK_ID_OR_NAME",
			""),
		builder.Command(&up.Up{},
			"Bring up a stack",
			appName+" up [OPTIONS]",
			""),
	}
	app.Commands = append(app.Commands,
		waiter.WaitCommand())

	err := app.Run(reformatArgs(os.Args))
	if err != nil {
		logrus.Fatal(err)
	}
}

func reformatArgs(args []string) []string {
	var result []string
	words := -1
	for i, arg := range args {
		if arg == "--" {
			return append(result, args[i:]...)
		}

		if len(arg) > 0 && arg[0:1] != "-" {
			words++
			if words > 1 {
				return append(result, args[i:]...)
			}
			result = append(result, arg)
			continue
		}

		words = 0

		if len(arg) <= 2 || arg[1:2] == "-" {
			result = append(result, arg)
			continue
		}

		for _, chars := range arg[1:] {
			result = append(result, "-"+string(chars))
		}
	}

	return result
}
