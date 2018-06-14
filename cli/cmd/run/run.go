package run

import (
	"github.com/rancher/rio/cli/cmd/create"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	"github.com/urfave/cli"
)

type Run struct {
	create.Create
	Scale int `desc:"scale" default:"1"`
}

func (r *Run) Run(app *cli.Context) error {
	return r.RunCallback(app, func(service *client.Service) *client.Service {
		service.Scale = int64(r.Scale)
		return service
	})
}
