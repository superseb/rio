package ps

import (
	"fmt"
	"strconv"

	"sort"

	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/cli/pkg/table"
	"github.com/rancher/rio/cli/server"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	"github.com/urfave/cli"
)

type ServiceData struct {
	ID       string
	Service  *client.Service
	Stack    *client.Stack
	Endpoint string
}

func FormatScale(data, data2 interface{}) (string, error) {
	scale, ok := data.(int64)
	if !ok {
		return fmt.Sprint(data), nil
	}
	scaleStr := strconv.FormatInt(scale, 10)

	scaleStatus, ok := data2.(*client.ScaleStatus)
	if !ok || scaleStatus == nil {
		return scaleStr, nil
	}

	if scaleStatus.Available == 0 && scaleStatus.Unavailable == 0 && scaleStatus.Ready == scale {
		return scaleStr, nil
	}

	percentage := ""
	if scale > 0 && scaleStatus.Updated > 0 && scale != scaleStatus.Updated {
		percentage = fmt.Sprintf(" %d%%", (scaleStatus.Updated*100)/scale)
	}

	return fmt.Sprintf("(%d/%d/%d)/%d%s", scaleStatus.Unavailable, scaleStatus.Available, scaleStatus.Ready, scale, percentage), nil
}

func (p *Ps) services(app *cli.Context, ctx *server.Context) error {
	stacks, err := ctx.Client.Stack.List(nil)
	if err != nil {
		return err
	}

	services, err := ctx.Client.Service.List(nil)
	if err != nil {
		return err
	}

	writer := table.NewWriter([][]string{
		{"NAME", "{{serviceName .Stack.Name .Service.Name}}"},
		{"IMAGE", "Service.Image"},
		{"CREATED", "{{.Service.Created | ago}}"},
		{"SCALE", "{{scale .Service.Scale .Service.ScaleStatus}}"},
		{"STATE", "Service.State"},
		{"ENDPOINT", "Endpoint"},
		{"DETAIL", "Service.TransitioningMessage"},
	}, app)
	defer writer.Close()

	writer.AddFormatFunc("scale", FormatScale)

	stackByID := map[string]*client.Stack{}
	for i, stack := range stacks.Data {
		stackByID[stack.ID] = &stacks.Data[i]
	}

	sort.Slice(services.Data, func(i, j int) bool {
		return services.Data[i].ID < services.Data[j].ID
	})

	for i, service := range services.Data {
		writer.Write(&ServiceData{
			ID:       service.ID,
			Service:  &services.Data[i],
			Stack:    stackByID[service.StackID],
			Endpoint: endpoint(ctx, stackByID[service.StackID], service.PortBindings, &service),
		})

		for revName, revision := range service.Revisions {
			newService := &client.Service{}
			if err := convert.ToObj(&revision, newService); err != nil {
				return err
			}
			newService.Name += service.Name + ":" + revName
			newService.Created = service.Created
			if newService.Image == "" {
				newService.Image = service.Image
			}

			writer.Write(&ServiceData{
				ID:      service.ID,
				Service: newService,
				Stack:   stackByID[service.StackID],
				// use parent service ports
				Endpoint: endpoint(ctx, stackByID[service.StackID], service.PortBindings, newService),
			})
		}
	}

	return writer.Err()
}

func endpoint(ctx *server.Context, stack *client.Stack, ports []client.PortBinding, service *client.Service) string {
	if ctx.Domain == "" {
		return ""
	}

	for _, port := range ports {
		if port.Protocol == "http" {
			name, rev := kv.Split(service.Name, ":")
			domain := fmt.Sprintf("%s.%s.%s", name, stack.Name, ctx.Domain)
			if rev != "" && rev != "latest" {
				domain = rev + "." + domain
			}

			return "http://" + domain
		}
	}

	return ""
}
