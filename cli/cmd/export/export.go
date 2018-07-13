package export

import (
	"io"
	"os"

	"github.com/rancher/rio/cli/cmd/util"
	"github.com/rancher/rio/cli/pkg/output"
	"github.com/rancher/rio/cli/pkg/yamldownload"
	"github.com/rancher/rio/cli/server"
	"github.com/urfave/cli"
)

type Export struct {
	T_Type   string `desc:"Export specific type"`
	O_Output string `desc:"Output format (yaml/json)"`
}

func (e *Export) Run(app *cli.Context) error {
	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}
	defer ctx.Close()

	format, err := output.Format(e.O_Output)
	if err != nil {
		return err
	}

	args := app.Args()
	if len(args) == 0 {
		args = []string{ctx.DefaultStackName}
	}

	for _, arg := range args {
		types := util.ExportEditTypes
		if e.T_Type != "" {
			types = []string{e.T_Type}
		}
		_, body, _, err := yamldownload.DownloadYAML(ctx, format, "export", arg, types...)
		if err != nil {
			return err
		}
		defer body.Close()

		_, err = io.Copy(os.Stdout, body)
		if err != nil {
			return err
		}
	}

	return nil
}
