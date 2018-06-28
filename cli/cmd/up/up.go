package up

import (
	"os"
	"strings"

	"fmt"

	"github.com/pkg/errors"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rio/cli/cmd/util"
	"github.com/rancher/rio/cli/pkg/up"
	"github.com/rancher/rio/cli/pkg/waiter"
	"github.com/rancher/rio/cli/server"
	"github.com/rancher/rio/pkg/yaml"
	"github.com/urfave/cli"
)

type Up struct {
	F_File    string `desc:"The template to apply"`
	S_Stack   string `desc:"Stack to use (id or name)"`
	A_Answers string `desc:"Answer file in with key/value pairs in yaml or json"`
	Prompt    bool   `desc:"Re-ask all questions if answer is not found in environment variables"`
}

func (u *Up) Run(app *cli.Context) error {
	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}

	content, err := util.ReadFile(u.F_File)
	if err != nil {
		return errors.Wrapf(err, "reading %s", u.F_File)
	}

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

	answers, err := readAnswers(u.A_Answers)
	if err != nil {
		return fmt.Errorf("failed to parse answer file [%s]: %v", u.A_Answers, err)
	}

	if err := up.Run(ctx, content, stackID, false, u.Prompt, answers); err != nil {
		return err
	}

	return waiter.WaitFor(app, stackID)
}

func readAnswers(answersFile string) (map[string]string, error) {
	content, err := util.ReadFile(answersFile)
	if err != nil {
		return nil, err
	}

	data, err := yaml.Parse(content)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for k, v := range data {
		result[k] = convert.ToString(v)
	}

	return result, nil
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
