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

//func runFlannel(config *AgentConfig) error {
//	flannel.Main([]string{
//		"--ip-masq",
//		"--kubeconfig-file", config.KubeConfig,
//	})
//
//	logrus.Fatalf("flannel exited")
//	return nil
//}

func run() error {
	return server.StartServer(context.Background(), 8080, 8443)
}
