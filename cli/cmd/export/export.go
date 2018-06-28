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
	O_Output string `desc:"Output format (yaml/json)"`
}

func (e *Export) Run(app *cli.Context) error {
	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}

	format, err := output.Format(e.O_Output)
	if err != nil {
		return err
	}

	args := app.Args()
	if len(args) == 0 {
		args = []string{"default"}
	}

	for _, arg := range args {
		_, body, _, err := yamldownload.DownloadYAML(ctx, format, "export", arg, util.ExportEditTypes...)
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
