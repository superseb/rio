// +build k3s

package server

import (
	"context"
	"flag"

	"github.com/emicklei/go-restful-swagger12"
	"github.com/rancher/rio/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func setupLogging(app *cli.Context) {
	if !app.GlobalBool("debug") {
		swagger.LogInfo = func(format string, v ...interface{}) {
			logrus.Debugf(format, v...)
		}
		flag.Set("stderrthreshold", "3")
		flag.Set("alsologtostderr", "false")
		flag.Set("logtostderr", "false")
	}
}

func (s *Server) Run(app *cli.Context) error {
	setupLogging(app)
	logrus.Info("Starting Rio ", app.App.Version)
	return server.StartServer(context.Background(), s.D_DataDir, s.L_HttpListenPort, s.P_HttpsListenPort, !s.DisableControllers)
}
