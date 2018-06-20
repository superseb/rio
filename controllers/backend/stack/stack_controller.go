package stack

import (
	"context"

	"github.com/rancher/rio/pkg/apply"
	"github.com/rancher/rio/pkg/namespace"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/rancher/types/apis/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	projectID = "field.cattle.io/projectId"
)

func Register(ctx context.Context, rContext *types.Context) {
	s := &stackController{
		namespaceLister: rContext.Core.Namespaces("").Controller().Lister(),
		namespaces:      rContext.Core.Namespaces(""),
	}
	rContext.Rio.Stacks("").AddLifecycle("stack-controller", s)
}

type stackController struct {
	namespaceLister v1.NamespaceLister
	namespaces      v1.NamespaceInterface
}

func (s *stackController) Create(obj *v1beta1.Stack) (*v1beta1.Stack, error) {
	return obj, nil
}

func (s *stackController) Remove(obj *v1beta1.Stack) (*v1beta1.Stack, error) {
	err := s.namespaces.Delete(namespace.StackToNamespace(obj), nil)
	if errors.IsNotFound(err) {
		return obj, nil
	}
	return obj, err
}

func (s *stackController) Updated(stack *v1beta1.Stack) (*v1beta1.Stack, error) {
	stack, err := s.createBackingNamespace(stack)
	if err != nil {
		return stack, err
	}

	internalStack, err := s.parseServices(stack)
	if err != nil {
		// if parsing fails we don't return err because it's a user error
		return stack, nil
	}

	objects := s.gatherObjects(stack, internalStack)

	err = apply.Apply(objects, "stack-"+stack.Name, stack.Generation)
	return stack, err
}
