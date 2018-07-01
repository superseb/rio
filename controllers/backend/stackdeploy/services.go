package stackdeploy

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *stackDeployController) services(objects []runtime.Object, volumes []*v1beta1.Volume, namespace string) ([]runtime.Object, error) {
	services, err := s.serviceLister.List(namespace, labels.Everything())
	if err != nil {
		return objects, err
	}

	volumeDefs := volumeMap(volumes)

	for _, service := range services {
		objects, err = s.service(objects, service.Name, namespace, service, volumeDefs)
		if err != nil {
			return objects, errors.Wrapf(err, "failed to construct service for %s/%s", service.Namespace, service.Name)
		}
	}

	return objects, nil
}

func MergeRevisionToService(service *v1beta1.Service, revision string) (*v1beta1.ServiceUnversionedSpec, map[string]string, error) {
	// TODO: do better merging
	newRevision := service.Spec.ServiceUnversionedSpec.DeepCopy()
	serviceRevision, ok := service.Spec.Revisions[revision]
	if !ok {
		return nil, nil, fmt.Errorf("failed to find revision for %s", revision)
	}

	err := convert.ToObj(&serviceRevision.Spec, newRevision)
	return newRevision, mergeLabels(service.Labels, serviceRevision.Labels), err
}

func (s *stackDeployController) service(objects []runtime.Object, name, namespace string, service *v1beta1.Service, volumeDefs map[string]*v1beta1.Volume) ([]runtime.Object, error) {
	objects, err := s.addService(objects, "latest", name, namespace, &service.Spec.ServiceUnversionedSpec, service.Labels, volumeDefs)

	for revision := range service.Spec.Revisions {
		newRevision, serviceLabels, err := MergeRevisionToService(service, revision)
		if err != nil {
			return nil, err
		}

		objects, err = s.addService(objects, revision, name, namespace, newRevision, serviceLabels, volumeDefs)
		if err != nil {
			return nil, err
		}
	}

	return objects, err
}

func (s *stackDeployController) addService(objects []runtime.Object, revision, serviceName, namespace string, service *v1beta1.ServiceUnversionedSpec,
	serviceLabels map[string]string, volumeDefs map[string]*v1beta1.Volume) ([]runtime.Object, error) {
	var (
		err error
	)

	labels := map[string]string{
		"rio.cattle.io":           "true",
		"rio.cattle.io/service":   serviceName,
		"rio.cattle.io/namespace": namespace,
		"rio.cattle.io/revision":  revision,
	}

	name := fmt.Sprintf("%s-%s", serviceName, revision)
	if revision == "latest" {
		name = serviceName
	}

	if service.Global {
		objects, err = daemonset(objects, labels, serviceLabels, name, serviceName, namespace, service, volumeDefs)
	} else if isDeployment(serviceName, namespace, service, volumeDefs) {
		objects = deployment(objects, labels, serviceLabels, name, serviceName, namespace, service, volumeDefs)
	} else {
		objects, err = statefulset(objects, labels, serviceLabels, name, serviceName, namespace, service, volumeDefs)
	}
	if err != nil {
		return objects, err
	}
	objects = serviceSelector(objects, name, namespace, labels)
	objects = nodePorts(objects, name+"-ports", namespace, service, labels)
	objects = pdbs(objects, name, namespace, labels, service)

	return objects, nil
}
