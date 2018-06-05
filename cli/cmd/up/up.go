package up

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/rio/cli/pkg/waiter"
	"github.com/rancher/rio/cli/server"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	"github.com/urfave/cli"
)

type Up struct {
	F_File  string `desc:"The template to apply"`
	S_Stack string `desc:"Stack to use (id or name)"`
}

func (u *Up) Run(app *cli.Context) error {
	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}

	bytes, err := readFile(u.F_File)
	if err != nil {
		return errors.Wrapf(err, "reading %s", u.F_File)
	}

	content := base64.StdEncoding.EncodeToString(bytes)

	stackName, err := getStackName(u.F_File, u.S_Stack)
	if err != nil {
		return err
	}

	if len(stackName) > 0 && !strings.HasSuffix(stackName, "/") {
		stackName += "/"
	}

	_, stackID, stackName, err := ctx.ResolveSpaceStackName(stackName)
	if err != nil {
		return err
	}

	stack, err := ctx.Client.Stack.ByID(stackID)
	if err != nil {
		return err
	}

	_, err = ctx.Client.Stack.Update(stack, &client.Stack{
		Templates: map[string]string{
			u.F_File: content,
		},
	})

	return waiter.WaitFor(app, stack.ID)
}

func readFile(file string) ([]byte, error) {
	if file == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(file)
}

func getStackName(file, stack string) (string, error) {
	if stack != "" {
		return stack, nil
	}
	if strings.HasSuffix(file, "-stack.yml") || strings.HasSuffix(file, "-stack.yaml") {
		file = strings.TrimSuffix(file, "-stack.yml")
		file = strings.TrimSuffix(file, "-stack.yaml")
		return file, nil
	}
	stack, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "can not determine stack file from current directory")
	}

	return stack, nil
}
