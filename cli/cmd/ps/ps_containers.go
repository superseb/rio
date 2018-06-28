package ps

import (
	"strings"

	"github.com/rancher/rio/cli/cmd/util"
	"github.com/rancher/rio/cli/pkg/table"
	"github.com/rancher/rio/cli/server"
	spaceclient "github.com/rancher/rio/types/client/space/v1beta1"
	"github.com/urfave/cli"
)

var (
	ignoreImages = []string{
		"gcr.io/runconduit/proxy:",
		"gcr.io/runconduit/proxy-init:",
	}
)

type ContainerData struct {
	ID        string
	Pod       *spaceclient.Pod
	Container *spaceclient.Container
}

func (p *Ps) containers(app *cli.Context, ctx *server.Context) error {
	c, err := ctx.SpaceClient()
	if err != nil {
		return err
	}

	pods, err := c.Pod.List(util.DefaultListOpts())
	if err != nil {
		return err
	}

	writer := table.NewWriter([][]string{
		{"NAME", "{{.Container.Name}}"},
		{"IMAGE", "Container.Image"},
		{"CREATED", "{{.Pod.Created | ago}}"},
		{"NODE", "Pod.NodeName"},
		{"IP", "Pod.PodIP"},
		{"STATE", "Container.State"},
		{"DETAIL", "Container.TransitioningMessage"},
	}, app)
	defer writer.Close()

	for i, pod := range pods.Data {
		containers := append(pod.Containers, pod.InitContainers...)
		for j, container := range containers {
			serviceName := pod.Labels["rio.cattle.io/service"]
			if !p.A_All && (serviceName == "" || ignoreImage(container.Image)) {
				continue
			}
			id := pod.ID
			if container.Name != pod.Labels["rio.cattle.io/service"] {
				id = id + "/" + container.Name
			}
			writer.Write(&ContainerData{
				ID:        id,
				Pod:       &pods.Data[i],
				Container: &containers[j],
			})
		}
	}

	return writer.Err()
}

func ignoreImage(image string) bool {
	for _, bad := range ignoreImages {
		if strings.HasPrefix(image, bad) {
			return true
		}
	}
	return false
}
