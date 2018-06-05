package server

import (
	"context"

	"github.com/rancher/rio/pkg/server"
	"github.com/urfave/cli"
)

type Server struct {
	HttpsListenPort int `desc:"HTTPS listen port" default:"443"`
	HttpListenPort  int `desc:"HTTP listen port" default:"80"`
}

func (s *Server) Run(app *cli.Context) error {
	return server.StartServer(context.Background(), s.HttpListenPort, s.HttpsListenPort)
}
