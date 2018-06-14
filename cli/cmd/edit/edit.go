package edit

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/rio/cli/pkg/output"
	"github.com/rancher/rio/cli/pkg/waiter"
	"github.com/rancher/rio/cli/pkg/yamldownload"
	"github.com/rancher/rio/cli/server"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	"github.com/urfave/cli"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/editor"
)

type Edit struct {
	O_Output string `desc:"Output format to edit (yaml|json)" default:"yaml"`
}

func (edit *Edit) Run(app *cli.Context) error {
	ctx, err := server.NewContext(app)
	if err != nil {
		return err
	}

	format, err := output.Format(edit.O_Output)
	if err != nil {
		return err
	}

	waiter, err := waiter.NewWaiter(app)
	if err != nil {
		return err
	}

	for _, arg := range app.Args() {
		obj, body, url, err := yamldownload.DownloadYAML(ctx, format, "edit", arg, client.ServiceType)
		if err != nil {
			return err
		}

		e := editor.NewDefaultEditor(os.Environ())
		content, _, err := e.LaunchTempFile("rio-", "-edit", body)
		if err != nil {
			return err
		}

		if err := update(ctx, url, content); err != nil {
			return err
		}

		waiter.Add(obj.ID)
	}

	return waiter.Wait()
}

func update(ctx *server.Context, self string, content []byte) error {
	parsed, err := url.Parse(self)
	if err != nil {
		return err
	}

	q := parsed.Query()
	q.Set("_replace", "true")
	parsed.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPut, parsed.String(), bytes.NewReader(content))
	if err != nil {
		return err
	}

	ctx.Client.Ops.SetupRequest(req)
	req.Header.Set("Content-Type", "application/yaml")

	resp, err := ctx.Client.Ops.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 300 {
		return clientbase.NewAPIError(resp, parsed.String())
	}

	return nil
}
