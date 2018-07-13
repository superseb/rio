package deploy

import (
	"fmt"
	"strings"

	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func serviceSelector(objects []runtime.Object, name, namespace string, service *v1beta1.ServiceUnversionedSpec, labels map[string]string) []runtime.Object {
	svc := newServiceSelector(name, namespace, labels)

	eps := service.ExposedPorts
	for _, container := range service.Sidekicks {
		eps = append(eps, container.ExposedPorts...)
	}

	if len(eps) > 0 {
		svc.Spec.ClusterIP = ""
		svc.Spec.Ports = nil
	}

	names := map[string]bool{}
	for i, port := range eps {
		name := port.Name
		if name == "" {
			name = fmt.Sprintf("port-%d", i)
		}
		if names[name] {
			name = fmt.Sprintf("%s-%d", name, i)
		}
		names[name] = true
		servicePort := v1.ServicePort{
			Name:       name,
			TargetPort: intstr.FromInt(int(port.TargetPort)),
			Port:       int32(port.Port),
			Protocol:   v1.ProtocolTCP,
		}

		if servicePort.Port == 0 {
			servicePort.Port = servicePort.TargetPort.IntVal
		}

		if strings.EqualFold(port.Protocol, "udp") {
			servicePort.Protocol = v1.ProtocolUDP
		}

		if port.IP != "" {
			svc.Spec.ClusterIP = port.IP
		}

		svc.Spec.Ports = append(svc.Spec.Ports, servicePort)
	}

	objects = append(objects, svc)
	return objects
}

func newServiceSelector(name, namespace string, labels map[string]string) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: v1.ClusterIPNone,
			Selector:  labels,
			Ports: []v1.ServicePort{
				{
					Name:       "default",
					Protocol:   v1.ProtocolTCP,
					TargetPort: intstr.FromInt(80),
					Port:       80,
				},
			},
		},
	}
}
