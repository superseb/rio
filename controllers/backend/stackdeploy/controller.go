package stackdeploy

import (
	"context"
	"strings"

	"github.com/rancher/rio/pkg/apply"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, rContext *types.Context) {
	s := &stackDeployController{
		stackController: rContext.Rio.Stacks("").Controller(),
		serviceLister:   rContext.Rio.Services("").Controller().Lister(),
	}

	rContext.Rio.Stacks("").AddHandler("stack-deploy-controller", s.deploy)
	rContext.Rio.Services("").AddHandler("stack-deploy-controller", func(key string, obj *v1beta1.Service) error {
		if obj != nil {
			return s.enqueue(obj)
		}
		return nil
	})
}

type stackDeployController struct {
	stackController v1beta1.StackController
	serviceLister   v1beta1.ServiceLister
}

func (s *stackDeployController) enqueue(obj runtime.Object) error {
	if obj == nil {
		return nil
	}
	if meta, ok := obj.(metav1.Object); ok {
		s.stackController.Enqueue("", "/"+meta.GetNamespace())
	}

	return nil
}

func (s *stackDeployController) deploy(key string, _ *v1beta1.Stack) error {
	if !strings.HasPrefix(key, "/") {
		return nil
	}

	return s.deployNamespace(key[1:])
}

func (s *stackDeployController) deployNamespace(namespace string) error {
	objects, err := s.services(nil, namespace)
	if err != nil {
		return err
	}

	return apply.Apply(objects, "stackdeploy-"+namespace, 0)
}
