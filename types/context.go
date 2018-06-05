package types

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/store/proxy"
	"github.com/rancher/norman/types"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1/schema"
	spacev1beta1 "github.com/rancher/rio/types/apis/space.cattle.io/v1beta1"
	spaceSchema "github.com/rancher/rio/types/apis/space.cattle.io/v1beta1/schema"
	"github.com/rancher/types/apis/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Context struct {
	LocalConfig  *rest.Config
	Schemas      *types.Schemas
	Global       spacev1beta1.Interface
	Rio        v1beta1.Interface
	Core         v1.Interface
	K8s          kubernetes.Interface
	ClientGetter proxy.ClientGetter
}

func NewContext(restConfig rest.Config) (*Context, error) {
	g, err := spacev1beta1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	r, err := v1beta1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	c, err := v1.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	k, err := kubernetes.NewForConfig(&restConfig)
	if err != nil {
		return nil, err
	}

	cg, err := proxy.NewClientGetterFromConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &Context{
		LocalConfig: &restConfig,
		Schemas: types.NewSchemas().
			AddSchemas(spaceSchema.Schemas).
			AddSchemas(schema.Schemas),
		Global:       g,
		Rio:        r,
		Core:         c,
		K8s:          k,
		ClientGetter: cg,
	}, nil
}

func (c *Context) Start(ctx context.Context) error {
	return controller.SyncThenStart(ctx, 5,
		c.Rio,
		c.Core,
		c.Global)
}
