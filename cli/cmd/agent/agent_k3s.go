// +build k3s

package agent

import (
	"fmt"
	"os"

	"path/filepath"

	"time"

	"github.com/rancher/rio/cli/pkg/resolvehome"
	"github.com/rancher/rio/pkg/clientaccess"
	"github.com/rancher/rio/pkg/enterchroot"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func (a *Agent) Run(app *cli.Context) error {
	if os.Getuid() != 0 {
		return fmt.Errorf("agent must be ran as root")
	}

	if len(a.T_Token) == 0 {
		return fmt.Errorf("--token is required")
	}

	if len(a.S_Server) == 0 {
		return fmt.Errorf("--server is required")
	}

	dataDir, err := resolvehome.Resolve(a.D_DataDir)
	if err != nil {
		return err
	}

	for {
		tmpFile, err := clientaccess.AccessInfoToTempKubeConfig("", a.S_Server, a.T_Token)
		if err != nil {
			logrus.Error(err)
			time.Sleep(2 * time.Second)
			continue
		}
		os.Remove(tmpFile)
		break
	}

	os.Setenv("RIO_URL", a.S_Server)
	os.Setenv("RIO_TOKEN", a.T_Token)
	os.Setenv("RIO_DATA_DIR", filepath.Join(a.D_DataDir, "root"))

	os.MkdirAll(dataDir, 0700)

	return enterchroot.Mount(dataDir)
}
