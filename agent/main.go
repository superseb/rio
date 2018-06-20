package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/cmd/containerd/command"
	"github.com/coreos/flannel"
	"github.com/rancher/rio/pkg/clientaccess"
	"github.com/sirupsen/logrus"
	"k8s.io/kubernetes/cmd/agent"

	// Containerd
	_ "github.com/containerd/containerd/diff/walking/plugin"
	_ "github.com/containerd/containerd/gc/scheduler"
	_ "github.com/containerd/containerd/services/containers"
	_ "github.com/containerd/containerd/services/content"
	_ "github.com/containerd/containerd/services/diff"
	_ "github.com/containerd/containerd/services/events"
	_ "github.com/containerd/containerd/services/healthcheck"
	_ "github.com/containerd/containerd/services/images"
	_ "github.com/containerd/containerd/services/introspection"
	_ "github.com/containerd/containerd/services/leases"
	_ "github.com/containerd/containerd/services/namespaces"
	_ "github.com/containerd/containerd/services/snapshots"
	_ "github.com/containerd/containerd/services/tasks"
	_ "github.com/containerd/containerd/services/version"

	_ "github.com/containerd/containerd/linux"
	_ "github.com/containerd/containerd/metrics/cgroups"
	_ "github.com/containerd/containerd/snapshots/native"
	_ "github.com/containerd/containerd/snapshots/overlay"

	_ "github.com/containerd/cri"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	runContainerd()

	agentConfig, err := getConfig()
	if err != nil {
		return err
	}

	if err := agent.Agent(agentConfig); err != nil {
		return err
	}

	return runFlannel(agentConfig)
}

func runFlannel(config *agent.AgentConfig) error {
	flannel.Main([]string{
		"--ip-masq",
		"--kubeconfig-file", config.KubeConfig,
	})

	logrus.Fatalf("flannel exited")
	return nil
}

func runContainerd() {
	args := []string{
		"containerd",
		"-a", "/run/rio/containerd.sock",
		"--state", "/run/rio/containerd",
	}
	app := command.App()
	go func() {
		if err := app.Run(args); err != nil {
			fmt.Fprintf(os.Stderr, "containerd: %s\n", err)
			os.Exit(1)
		}
	}()
}

func getConfig() (*agent.AgentConfig, error) {
	u := os.Getenv("RIO_URL")
	if u == "" {
		return nil, fmt.Errorf("RIO_URL env var is required")
	}

	t := os.Getenv("RIO_TOKEN")
	if t == "" {
		return nil, fmt.Errorf("RIO_TOKEN env var is required")
	}

	dataDir := os.Getenv("RIO_DATA_DIR")
	if dataDir == "" {
		return nil, fmt.Errorf("RIO_DATA_DIR is required")
	}

	kubeConfig := filepath.Join(dataDir, "kubeconfig.yaml")

	_, cidr, _ := net.ParseCIDR("10.42.0.0/16")

	agentConfig := &agent.AgentConfig{
		ClusterCIDR:   *cidr,
		KubeConfig:    kubeConfig,
		RuntimeSocket: "/run/rio/containerd.sock",
		CNIBinDir:     "/usr/share/cni",
	}

	return agentConfig, clientaccess.AccessInfoToKubeConfig(kubeConfig, u, t)
}
