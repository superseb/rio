package exec

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/pkg/reexec"
	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/cli/server"
	spaceclient "github.com/rancher/rio/types/client/space/v1beta1"
	"github.com/urfave/cli"
)

type Exec struct {
	I_Stdin bool `desc:"Pass stdin to the container"`
	T_Tty   bool `desc:"Stdin is a TTY"`
}

func (e *Exec) Run(app *cli.Context) error {
	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}

	c, err := ctx.SpaceClient()
	if err != nil {
		return err
	}

	args := app.Args()
	if len(args) < 2 {
		return fmt.Errorf("at least two arguments are required CONTAINER CMD")
	}

	pod, containerName, err := FindPodAndContainer(args[0], c)
	if err != nil {
		return err
	}

	podNS, podName := kv.Split(pod.ID, ":")

	execArgs := []string{"kubectl", "--kubeconfig", "/home/darren/.gpath/142651f86a7318730641c8f1ca6c6150cbba0718/src/github.com/rancher/rio/a", "-v=9", "-n", podNS, "exec"}
	if e.I_Stdin {
		execArgs = append(execArgs, "-i")
	}
	if e.T_Tty {
		execArgs = append(execArgs, "-t")
	}

	execArgs = append(execArgs, podName, "-c", containerName)
	execArgs = append(execArgs, args[1:]...)

	fmt.Println("RUNNING", execArgs)

	//cmd := exec.Command(execArgs[0], execArgs[1:]...)
	cmd := reexec.Command(execArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func FindPodAndContainer(name string, c *spaceclient.Client) (*spaceclient.Pod, string, error) {
	podID, containerName := kv.Split(name, "/")
	stackName, serviceName := kv.Split(name, "/")
	if serviceName == "" {
		serviceName = stackName
		stackName = "default"
	}

	pods, err := c.Pod.List(nil)
	if err != nil {
		return nil, "", err
	}

	for i, pod := range pods.Data {
		if pod.ID == podID {
			containerName, ok := findContainerInPod(&pod, containerName)
			if ok {
				return &pods.Data[i], containerName, nil
			}
			return nil, "", fmt.Errorf("not found: %s", name)
		}

		if pod.Labels["rio.cattle.io/service"] == serviceName {
			ns := pod.Labels["rio.cattle.io/namespace"]
			i := strings.LastIndex(ns, "-")
			if i <= 0 {
				continue
			}

			if stackName == ns[:i] {
				containerName, ok := findContainerInPod(&pod, containerName)
				if ok {
					return &pod, containerName, nil
				}
			}
		}
	}

	return nil, "", fmt.Errorf("not found: %s", name)
}

func findContainerInPod(pod *spaceclient.Pod, containerName string) (string, bool) {
	if containerName == "" {
		return pod.Containers[0].Name, true
	}
	for _, container := range pod.Containers {
		if container.Name == containerName {
			return containerName, true
		}
	}
	for _, container := range pod.InitContainers {
		if container.Name == containerName {
			return containerName, true
		}
	}

	return containerName, false
}
