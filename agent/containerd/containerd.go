package containerd

import (
	"fmt"
	"os"

	// containerd builtin
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

	// containerd cri
	_ "github.com/containerd/cri"

	// containerd linux
	"os/exec"

	"time"

	_ "github.com/containerd/containerd/linux"
	_ "github.com/containerd/containerd/metrics/cgroups"
	_ "github.com/containerd/containerd/snapshots/native"
	_ "github.com/containerd/containerd/snapshots/overlay"
)

func Run() {
	args := []string{
		"containerd",
		"-a", "/run/rio/containerd.sock",
		"--state", "/run/rio/containerd",
	}
	//app := command.App()
	go func() {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		fmt.Fprintf(os.Stderr, "containerd: %s\n", err)
		os.Exit(1)
		//if err := app.Run(args); err != nil {
		//	fmt.Fprintf(os.Stderr, "containerd: %s\n", err)
		//	os.Exit(1)
		//}
	}()

	time.Sleep(1 * time.Second)
}
