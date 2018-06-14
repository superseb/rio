package server

import (
	"strings"

	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	spaceclient "github.com/rancher/rio/types/client/space/v1beta1"
	"github.com/urfave/cli"
)

type Context struct {
	Domain           string
	Client           *client.Client
	DefaultStackName string
	SpaceID          string
}

func NewContext(app *cli.Context) (*Context, error) {
	c, err := client.NewClient(&clientbase.ClientOpts{
		Insecure: true,
		URL:      "https://localhost:7443/v1beta1-rio/spaces/default",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, "https://localhost:7443/domain", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Ops.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	domain, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Context{
		Domain:           string(domain),
		Client:           c,
		DefaultStackName: "default",
		SpaceID:          "default",
	}, nil
}

func (c *Context) SpaceClient() (*spaceclient.Client, error) {
	return spaceclient.NewClient(&clientbase.ClientOpts{
		Insecure: true,
		URL:      "https://localhost:7443/v1beta1-rio",
	})
}

func (c *Context) ResolveSpaceStackName(in string) (string, string, string, error) {
	stackName, name := kv.Split(in, "/")
	if stackName != "" && name == "" {
		if !strings.HasSuffix(in, "/") {
			name = stackName
			stackName = ""
		}
	}

	if stackName == "" {
		stackName = c.DefaultStackName
	}

	stacks, err := c.Client.Stack.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": stackName,
		},
	})
	if err != nil {
		return "", "", "", errors.Wrapf(err, "failed to determine stack")
	}

	var s *client.Stack
	if len(stacks.Data) == 0 {
		s, err = c.Client.Stack.Create(&client.Stack{
			Name:    stackName,
			SpaceId: c.SpaceID,
		})
		if err != nil {
			return "", "", "", errors.Wrapf(err, "failed to create stack %s", stackName)
		}
	} else {
		s = &stacks.Data[0]
	}

	return s.SpaceId, s.ID, name, nil
}
