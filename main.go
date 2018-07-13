//go:generate go run types/codegen/cleanup/main.go
//go:generate go run types/codegen/main.go

package main

import (
	"context"
	"os"

	"github.com/docker/docker/pkg/reexec"
	_ "github.com/rancher/rio/pkg/kubectl"
	"github.com/rancher/rio/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if reexec.Init() {
		return
	}

	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	if os.Getenv("RIO_DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return server.StartServer(context.Background(), "./data-dir", 5080, 5443, true)
}
