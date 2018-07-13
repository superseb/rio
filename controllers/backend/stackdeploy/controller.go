package stackdeploy

import (
	"context"
	"reflect"
	"strings"

	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/pkg/deploy"
	"github.com/rancher/rio/pkg/namespace"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

const (
	stackByNS = "stackByNS"
)

func Register(ctx context.Context, rContext *types.Context) {
	s := &stackDeployController{
		stacks:          rContext.Rio.Stacks(""),
		stackController: rContext.Rio.Stacks("").Controller(),
		serviceLister:   rContext.Rio.Services("").Controller().Lister(),
		configLister:    rContext.Rio.Configs("").Controller().Lister(),
		volumeLister:    rContext.Rio.Volumes("").Controller().Lister(),
	}

	rContext.Rio.Stacks("").AddHandler("stack-deploy-controller", s.deploy)
	rContext.Rio.Configs("").AddHandler("stack-deploy-controller", func(key string, obj *v1beta1.Config) error {
		return s.enqueue(key)
	})
	rContext.Rio.Services("").AddHandler("stack-deploy-controller", func(key string, obj *v1beta1.Service) error {
		return s.enqueue(key)
	})
	rContext.Rio.Volumes("").AddHandler("stack-deploy-controller", func(key string, obj *v1beta1.Volume) error {
		return s.enqueue(key)
	})

	s.stackController.Informer().AddIndexers(cache.Indexers{
		stackByNS: func(obj interface{}) ([]string, error) {
			if obj == nil {
				return nil, nil
			}
			s, ok := obj.(*v1beta1.Stack)
			if !ok {
				return nil, nil
			}
			return []string{
				namespace.StackToNamespace(s),
			}, nil
		},
	})
}

type stackDeployController struct {
	stacks          v1beta1.StackInterface
	stackController v1beta1.StackController
	serviceLister   v1beta1.ServiceLister
	configLister    v1beta1.ConfigLister
	volumeLister    v1beta1.VolumeLister
}

func (s *stackDeployController) enqueue(key string) error {
	ns, name := kv.Split(key, "/")
	if ns == "" || name == "" {
		return nil
	}
	s.stackController.Enqueue("", "/"+ns)
	return nil
}

func (s *stackDeployController) deploy(key string, _ *v1beta1.Stack) error {
	if !strings.HasPrefix(key, "/") {
		return nil
	}

	objs, err := s.stackController.Informer().GetIndexer().ByIndex(stackByNS, key[1:])
	if err != nil {
		return err
	}

	if len(objs) != 1 {
		return nil
	}

	stack, ok := objs[0].(*v1beta1.Stack)
	if !ok {
		return nil
	}

	newStack := stack.DeepCopy()
	_, err = v1beta1.StackConditionDeployed.Do(newStack, func() (runtime.Object, error) {
		return newStack, s.deployNamespace(key[1:])
	})

	if !reflect.DeepEqual(stack, newStack) {
		s.stacks.Update(newStack)
	}

	return err
}

func (s *stackDeployController) getStack(namespace string) (*deploy.StackResources, error) {
	configs, err := s.configLister.List(namespace, labels.Everything())
	if err != nil {
		return nil, err
	}

	volumes, err := s.volumeLister.List(namespace, labels.Everything())
	if err != nil {
		return nil, err
	}

	services, err := s.serviceLister.List(namespace, labels.Everything())
	if err != nil {
		return nil, err
	}

	return &deploy.StackResources{
		Configs:  configs,
		Volumes:  volumes,
		Services: services,
	}, nil
}

func (s *stackDeployController) deployNamespace(namespace string) error {
	stack, err := s.getStack(namespace)
	if err != nil {
		return err
	}
	return deploy.Deploy(namespace, stack)
}
