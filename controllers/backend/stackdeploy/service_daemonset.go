package stackdeploy

import (
	"fmt"

	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func daemonset(objects []runtime.Object, labels map[string]string, serviceLabels map[string]string, depName, serviceName, namespace string, service *v1beta1.ServiceUnversionedSpec, volumes map[string]*v1beta1.Volume) ([]runtime.Object, error) {
	usedTemplates, podSpec := podSpec(serviceName, serviceLabels, service, volumes)
	scaleParams := parseScaleParams(service)

	daemonSet := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        depName,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: map[string]string{},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: mergeLabels(labels, serviceLabels),
				},
				Spec: podSpec,
			},
		},
	}

	if service.UpdateStrategy == "on-delete" {
		daemonSet.Spec.UpdateStrategy.Type = appsv1.OnDeleteStatefulSetStrategyType
	} else {
		daemonSet.Spec.UpdateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType
		if scaleParams.Scale > 0 && scaleParams.BatchSize > 0 {
			daemonSet.Spec.UpdateStrategy.RollingUpdate = &appsv1.RollingUpdateDaemonSet{
				MaxUnavailable: scaleParams.MaxUnavailable,
			}
		}
	}

	if len(usedTemplates) > 0 {
		return nil, fmt.Errorf("globally scheduling services can not use volume templates")
	}

	return append(objects, daemonSet), nil
}
