package main

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/reexec"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/rio/cli/cmd/agent"
	"github.com/rancher/rio/cli/cmd/config"
	"github.com/rancher/rio/cli/cmd/create"
	"github.com/rancher/rio/cli/cmd/edit"
	"github.com/rancher/rio/cli/cmd/exec"
	"github.com/rancher/rio/cli/cmd/export"
	"github.com/rancher/rio/cli/cmd/inspect"
	"github.com/rancher/rio/cli/cmd/node"
	"github.com/rancher/rio/cli/cmd/promote"
	"github.com/rancher/rio/cli/cmd/ps"
	"github.com/rancher/rio/cli/cmd/rm"
	"github.com/rancher/rio/cli/cmd/run"
	"github.com/rancher/rio/cli/cmd/scale"
	"github.com/rancher/rio/cli/cmd/server"
	"github.com/rancher/rio/cli/cmd/stack"
	"github.com/rancher/rio/cli/cmd/stage"
	"github.com/rancher/rio/cli/cmd/up"
	"github.com/rancher/rio/cli/cmd/volume"
	"github.com/rancher/rio/cli/cmd/weight"
	"github.com/rancher/rio/cli/pkg/builder"
	"github.com/rancher/rio/cli/pkg/waiter"
	_ "github.com/rancher/rio/pkg/conduit"
	_ "github.com/rancher/rio/pkg/kubectl"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	desc = `For node scheduling arguments (--node-require, --node-require-any,
   --node-preferred) the expression is evaluated against the node labels using the
   following syntax.


      foo=bar              Node with label foo must have the value bar
      foo!=bar             If node has a label with key foo it must not equal to bar
      foo                  Node must have a label with key foo and any value
      foo in (bar, baz)    Node must have a label with key foo and a value of bar or baz
      foo notin (bar, baz) If node has a label with key foo it must not equal to bar or baz
      foo > 3              Node must have a label with key foo and a value greater than 3
      foo < 3              Node must have a label with key foo and a value less than 3
      !foo                 Node must not have a label with key foo with any value
      expr && expr         Any above expression can be combined so that all expression must be true
`
)

var (
	appName = filepath.Base(os.Args[0])
	Version = "dev"
)

func main() {
	if reexec.Init() {
		return
	}

	app := cli.NewApp()
	app.Name = appName
	app.Usage = "Containers as simple as they should be"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Turn on debug logs",
		},
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
			Usage:  "Specify the Rio API endpoint URL",
			EnvVar: "RIO_URL",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Specify Rio API token",
			EnvVar: "RIO_TOKEN",
		},
	}

	app.Commands = []cli.Command{
		stack.Stack(),
		config.Config(app),
		volume.Volume(app),
		builder.Command(&ps.Ps{},
			"List services and containers",
			appName+" ps [OPTIONS] [STACK...]",
			""),
		builder.Command(&exec.Exec{},
			"Run a command in a running container",
			appName+" exec [OPTIONS] CONTAINER COMMAND [ARG...]",
			""),
		builder.Command(&scale.Scale{},
			"Scale a service",
			appName+" scale [SERVICE=NUMBER...]",
			""),
		builder.Command(&server.Server{},
			"Run management server",
			appName+" server [OPTIONS]",
			""),
		builder.Command(&agent.Agent{},
			"Run node agent",
			appName+" agent [OPTIONS]",
			""),
		builder.Command(&run.Run{},
			"Run a new service",
			appName+" run [OPTIONS] IMAGE [COMMAND] [ARG...]",
			desc),
		builder.Command(&create.Create{},
			"Create a new service",
			appName+" create [OPTIONS] IMAGE [COMMAND] [ARG...]",
			desc),
		builder.Command(&rm.Rm{},
			"Delete a service or stack",
			appName+" rm ID_OR_NAME",
			""),
		builder.Command(&edit.Edit{},
			"Edit a service or stack",
			appName+" edit ID_OR_NAME",
			""),
		builder.Command(&export.Export{},
			"Export a stack",
			appName+" export STACK_ID_OR_NAME",
			""),
		builder.Command(&up.Up{},
			"Bring up a stack",
			appName+" up [OPTIONS]",
			""),
		builder.Command(&stage.Stage{},
			"Stage a new revision of a service",
			appName+" stage [OPTIONS] [SERVICE_ID_NAME]",
			""),
		builder.Command(&promote.Promote{},
			"Promote a staged version to latest",
			appName+" promote [SERVICE_ID_NAME]",
			""),
		builder.Command(&weight.Weight{},
			"Weight a percentage of traffic to a staged service",
			appName+" weight [OPTIONS] [SERVICE_REVISION=PERCENTAGE...]",
			""),
		builder.Command(&inspect.Inspect{},
			"Print the raw API output of a resource",
			appName+" inspect [ID_OR_NAME...]",
			""),
		node.Node(),
	}
	app.Commands = append(app.Commands,
		waiter.WaitCommand())
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			clientbase.Debug = true
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

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
