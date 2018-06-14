//go:generate go run types/codegen/cleanup/main.go
//go:generate go run types/codegen/main.go

package main

import (
	"context"

	"github.com/rancher/rio/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	logrus.SetLevel(logrus.DebugLevel)
	return server.StartServer(context.Background(), "./data-dir", 5080, 5443, true)
}
