package stackdeploy

import (
	"encoding/base64"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *stackDeployController) configs(objects []runtime.Object, namespace string) ([]runtime.Object, error) {
	configs, err := s.configLister.List(namespace, labels.Everything())
	if err != nil {
		return objects, err
	}

	for _, config := range configs {
		cfg := newConfig(config.Name, namespace, map[string]string{
			"rio.cattle.io/namespace": namespace,
		})
		if config.Spec.Encoded {
			bytes, err := base64.StdEncoding.DecodeString(config.Spec.Content)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to decode data for %s", config.Name)
			}
			cfg.BinaryData = map[string][]byte{
				"content": bytes,
			}
		} else {
			cfg.Data = map[string]string{
				"content": config.Spec.Content,
			}
		}
		objects = append(objects, cfg)
	}

	return objects, nil
}

func newConfig(name, namespace string, labels map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
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
