package export

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rancher/rio/cli/pkg/output"
	"github.com/rancher/rio/cli/pkg/yamldownload"
	"github.com/rancher/rio/cli/server"
	"github.com/urfave/cli"
)

type Export struct {
	O_Output string `desc:"Output format (yaml/json)"`
	F_File   string `desc:"Optional file to write to instead of stdout"`
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

	_, body, _, err := yamldownload.DownloadYAML(ctx, format, "export", app.Args())
	if err != nil {
		return err
	}
	defer body.Close()

	out := io.Writer(os.Stdout)
	if e.F_File != "" {
		f, err := os.Open(e.F_File)
		if err != nil {
			return errors.Wrapf(err, "failed to open %s", e.F_File)
		}
		defer f.Close()
		out = f
	}

	_, err = io.Copy(out, body)
	return err
}
