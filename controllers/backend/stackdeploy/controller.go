package stackdeploy

import (
	"context"
	"strings"

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
}

type stackDeployController struct {
	stackController v1beta1.StackController
	serviceLister   v1beta1.ServiceLister
}

func (s *stackDeployController) enqueue(obj runtime.Object) {
	if obj == nil {
		return
	}
	if meta, ok := obj.(metav1.Object); ok {
		s.stackController.Enqueue("", "/"+meta.GetNamespace())
	}
}

func (s *stackDeployController) deploy(key string, _ *v1beta1.Stack) error {
	if !strings.HasPrefix(key, "/") {
		return nil
	}

	return s.deployNamespace(key[1:])
}

func (s *stackDeployController) deployNamespace(namespace string) error {
	var objects []runtime.Object

	object, err := s.services(objects, namespace)
}
