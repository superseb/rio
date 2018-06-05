package stack

import (
	"github.com/rancher/rio/pkg/namespace"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *stackController) gatherObjects(stack *v1beta1.Stack, internalStack *v1beta1.InternalStack) []runtime.Object {
	var resources []runtime.Object

	ns := namespace.StackToNamespace(stack)
	resources = stacks(resources, ns, internalStack)

	return resources
}

func stacks(resources []runtime.Object, ns string, internalStack *v1beta1.InternalStack) []runtime.Object {
	for name, service := range internalStack.Services {
		newResource := service.DeepCopy()
		newResource.Kind = "Service"
		newResource.APIVersion = v1beta1.SchemeGroupVersion.String()
		newResource.Name = name
		newResource.Namespace = ns

		resources = append(resources, newResource)
	}

	return resources
}
