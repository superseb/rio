package server

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	"github.com/urfave/cli"
)

type Context struct {
	Client           *client.Client
	DefaultStackName string
	SpaceID          string
}

func NewContext(app *cli.Context) (*Context, error) {
	c, err := client.NewClient(&clientbase.ClientOpts{
		Insecure: true,
		URL:      "http://localhost:8083/v1beta1-rio/spaces/default",
	})
	if err != nil {
		return nil, err
	}

	return &Context{
		Client:           c,
		DefaultStackName: "default",
		SpaceID:          "default",
	}, nil
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
