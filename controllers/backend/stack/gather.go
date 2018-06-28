package stack

import (
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/rio/pkg/namespace"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *stackController) gatherObjects(stack *v1beta1.Stack, internalStack *v1beta1.InternalStack) []runtime.Object {
	var resources []runtime.Object

	ns := namespace.StackToNamespace(stack)
	resources = services(resources, stack, ns, internalStack)

	return resources
}

func services(resources []runtime.Object, stack *v1beta1.Stack, ns string, internalStack *v1beta1.InternalStack) []runtime.Object {
	for name, service := range internalStack.Services {
		newResource := service.DeepCopy()
		newResource.Kind = "Service"
		newResource.APIVersion = v1beta1.SchemeGroupVersion.String()
		newResource.Name = name
		newResource.Namespace = ns
		newResource.Spec.SpaceName = stack.Namespace
		newResource.Spec.StackName = ref.FromStrings(stack.Namespace, stack.Name)

		resources = append(resources, newResource)
	}

	return resources
}
