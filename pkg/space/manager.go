package space

import (
	ntypes "github.com/rancher/norman/types"
	"github.com/rancher/rio/types"
)

const (
	SpaceLabel = "space.rio.cattle.io"
)

type Manager struct {
}

func NewManager(context *types.Context) {
	context.Core.Namespaces("")

}

func (m *Manager) NamespaceForSpace(apiContext *ntypes.APIContext, space string) (string, error) {
	return "", nil
}
