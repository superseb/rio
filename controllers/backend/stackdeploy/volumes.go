package stackdeploy

import (
	"fmt"
	"strings"

	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func volumeToPVC(namespace string, labels map[string]string, volume v1beta1.Volume) (*v1.PersistentVolumeClaim, error) {
	cfg := newVolume(volume.Name, namespace, labels)
	if volume.Spec.Driver != "" && volume.Spec.Driver != "default" {
		cfg.Spec.StorageClassName = &volume.Spec.Driver
	}

	q, err := resource.ParseQuantity(fmt.Sprintf("%dGi", volume.Spec.SizeInGB))
	if err != nil {
		return nil, fmt.Errorf("failed to parse size [%d] on volume %s", volume.Spec.SizeInGB, volume.Name)
	}

	switch strings.ToLower(volume.Spec.AccessMode) {
	case "readwritemany":
		cfg.Spec.AccessModes = []v1.PersistentVolumeAccessMode{
			v1.ReadWriteMany,
		}
	case "readonlymany":
		cfg.Spec.AccessModes = []v1.PersistentVolumeAccessMode{
			v1.ReadOnlyMany,
		}
	default:
		cfg.Spec.AccessModes = []v1.PersistentVolumeAccessMode{
			v1.ReadWriteOnce,
		}
	}

	cfg.Spec.Resources.Requests = v1.ResourceList{
		v1.ResourceStorage: q,
	}

	return cfg, nil
}

func (s *stackDeployController) volumes(objects []runtime.Object, namespace string) ([]*v1beta1.Volume, []runtime.Object, error) {
	volumes, err := s.volumeLister.List(namespace, labels.Everything())
	if err != nil {
		return volumes, objects, err
	}

	for _, volume := range volumes {
		if volume.Spec.Template {
			continue
		}

		labels := map[string]string{
			"rio.cattle.io/namespace": namespace,
			"rio.cattle.io/volume":    volume.Name,
		}

		cfg, err := volumeToPVC(namespace, labels, *volume)
		if err != nil {
			return nil, nil, err
		}

		objects = append(objects, cfg)
	}

	return volumes, objects, nil
}

func newVolume(name, namespace string, labels map[string]string) *v1.PersistentVolumeClaim {
	return &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: map[string]string{},
		},
	}
}
