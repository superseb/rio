package conduit

import (
	"os"

	"github.com/docker/docker/pkg/reexec"
	"github.com/runconduit/conduit/cli/cmd"
)

func init() {
	reexec.Register("conduit", Main)
}
func Main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
