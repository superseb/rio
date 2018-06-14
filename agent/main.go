package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/coreos/flannel"
	"github.com/rancher/rio/pkg/clientaccess"
	"github.com/sirupsen/logrus"
	"k8s.io/kubernetes/cmd/agent"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
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
