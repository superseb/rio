package stackdeploy

import (
	"fmt"

	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *stackDeployController) volumes(objects []runtime.Object, namespace string) ([]runtime.Object, error) {
	volumes, err := s.volumeLister.List(namespace, labels.Everything())
	if err != nil {
		return objects, err
	}

	for i, volume := range volumes {
		if volume.Spec.Template {
			continue
		}

		cfg := newVolume(volume.Name, namespace, map[string]string{
			"rio.cattle.io/namespace": namespace,
			"rio.cattle.io/volume":    volume.Name,
		})
		if volume.Spec.Driver != "" && volume.Spec.Driver != "default" {
			cfg.Spec.StorageClassName = &volumes[i].Spec.Driver
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

		objects = append(objects, cfg)
	}

	return objects, nil
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
