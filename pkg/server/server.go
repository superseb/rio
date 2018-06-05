package server

import (
	"context"
	"net/http"
	"os"

	"github.com/rancher/norman/api"
	"github.com/rancher/norman/leader"
	"github.com/rancher/norman/signal"
	"github.com/rancher/rancher/k8s"
	"github.com/rancher/rancher/pkg/dynamiclistener"
	"github.com/rancher/rio/controllers/backend"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/space.cattle.io/v1beta1"
	"github.com/rancher/rio/types/config"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

func StartServer(ctx context.Context, httpPort, httpsPort int) error {
	ctx = signal.SigTermCancelContext(ctx)
	_, ctx, restConfig, err := k8s.GetConfig(ctx, "auto", os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	rContext, err := types.NewContext(*restConfig)
	if err != nil {
		return err
	}

	if err := config.SetupTypes(ctx, rContext); err != nil {
		return err
	}

	apiServer := api.NewAPIServer()
	if err := apiServer.AddSchemas(rContext.Schemas); err != nil {
		return err
	}

	apiRContext, err := types.NewContext(*restConfig)
	if err != nil {
		return err
	}
	apiRContext.Schemas = rContext.Schemas

	go leader.RunOrDie(ctx, "rio-controllers", rContext.K8s, func(ctx context.Context) {
		if err := backend.Register(ctx, rContext); err != nil {
			panic(err)
		}

		if err := rContext.Start(ctx); err != nil {
			panic(err)
		}

		<-ctx.Done()
	})

	if err := startServer(ctx, apiRContext, httpPort, httpsPort, apiServer); err != nil {
		return err
	}

	if err := apiRContext.Start(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func startServer(ctx context.Context, rContext *types.Context, httpPort, httpsPort int, handler http.Handler) error {
	s := &storage{
		listenConfigs:      rContext.Global.ListenConfigs(""),
		listenConfigLister: rContext.Global.ListenConfigs("").Controller().Lister(),
	}
	dynamiclistener.NewServer(ctx, s, handler, httpPort, httpsPort)
	return nil
}

type storage struct {
	listenConfigs      v1beta1.ListenConfigInterface
	listenConfigLister v1beta1.ListenConfigLister
}

func (s *storage) Update(lc *v3.ListenConfig) (*v3.ListenConfig, error) {
	updateLC := &v1beta1.ListenConfig{
		ListenConfig: *lc,
	}

	updateLC, err := s.listenConfigs.Update(updateLC)
	if err != nil {
		return nil, err
	}
	return &updateLC.ListenConfig, nil
}

func (s *storage) Get(namespace, name string) (*v3.ListenConfig, error) {
	lc, err := s.listenConfigLister.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	return &lc.ListenConfig, nil
}
