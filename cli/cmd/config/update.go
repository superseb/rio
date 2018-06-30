package config

import (
	"fmt"

	"encoding/base64"
	"unicode/utf8"

	"github.com/rancher/rio/cli/cmd/util"
	"github.com/rancher/rio/cli/pkg/lookup"
	"github.com/rancher/rio/cli/server"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	"github.com/urfave/cli"
)

type Update struct {
	L_Label map[string]string `desc:"Set meta data on a config"`
}

func (c *Update) Run(app *cli.Context) error {
	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}

	if len(app.Args()) != 2 {
		return fmt.Errorf("two arguments are required")
	}

	name := app.Args()[0]
	file := app.Args()[1]

	resource, err := lookup.Lookup(ctx.Client, name, client.ConfigType)
	if err != nil {
		return err
	}

	config, err := ctx.Client.Config.ByID(resource.ID)
	if err != nil {
		return err
	}

	content, err := util.ReadFile(file)
	if err != nil {
		return err
	}

	if len(c.L_Label) > 0 {
		config.Labels = c.L_Label
	}
	if utf8.Valid(content) {
		config.Content = string(content)
		config.Encoded = false
	} else {
		config.Content = base64.StdEncoding.EncodeToString(content)
		config.Encoded = true
	}

	config, err = ctx.Client.Config.Update(config, config)
	if err != nil {
		return err
	}

	fmt.Println(config.ID)
	return nil
}
