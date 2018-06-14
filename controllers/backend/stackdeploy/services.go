package stackdeploy

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (s *stackDeployController) services(objects []runtime.Object, namespace string) ([]runtime.Object, error) {
	services, err := s.serviceLister.List(namespace, labels.Everything())
	if err != nil {
		return objects, err
	}

	for _, service := range services {
		objects, err = s.service(objects, service.Name, namespace, service)
		if err != nil {
			return objects, errors.Wrapf(err, "failed to construct service for %s/%s", service.Namespace, service.Name)
		}
	}

	return objects, nil
}

func addScale(objects []runtime.Object, name, serviceName, namespace string, dep *v1beta2.Deployment, service *v1beta1.ServiceUnversionedSpec) []runtime.Object {
	scale := int32(service.Scale)
	batchSize := service.BatchSize

	if batchSize == 0 {
		batchSize = 1
	}

	if int32(batchSize) > scale {
		batchSize = int(scale)
	}

	surge := batchSize
	unavailable := 0

	if service.UpdateOrder == "start-first" {
		surge = batchSize
		unavailable = 0
	} else if service.UpdateOrder == "stop-first" {
		surge = 0
		unavailable = batchSize
	}

	maxSurge := intstr.FromInt(surge)
	maxUnavailable := intstr.FromInt(unavailable)

	dep.Spec.Replicas = &scale

	if scale > 0 && batchSize > 0 {
		dep.Spec.Strategy = v1beta2.DeploymentStrategy{
			RollingUpdate: &v1beta2.RollingUpdateDeployment{
				MaxSurge:       &maxSurge,
				MaxUnavailable: &maxUnavailable,
			},
			Type: v1beta2.RollingUpdateDeploymentStrategyType,
		}

		if batchSize < int(scale) {
			pdbSize := service.BatchSize
			if service.BatchSize > service.Scale {
				pdbSize = 1
			}
			pdbQuantity := intstr.FromInt(pdbSize)

			pdb := &v1beta12.PodDisruptionBudget{
				TypeMeta: metav1.TypeMeta{
					Kind:       "PodDisruptionBudget",
					APIVersion: "policy/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%d", name, pdbQuantity.IntVal),
					Namespace: namespace,
				},
				Spec: v1beta12.PodDisruptionBudgetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: dep.Spec.Selector.MatchLabels,
					},
					MaxUnavailable: &pdbQuantity,
				},
				Status: v1beta12.PodDisruptionBudgetStatus{
					DisruptedPods: map[string]metav1.Time{},
				},
			}

			objects = append(objects, pdb)
		}
	}

	return objects
}

func (s *stackDeployController) deployment(labels map[string]string, depName, serviceName, namespace string, service *v1beta1.ServiceUnversionedSpec) *v1beta2.Deployment {
	podSpec := podSpec(serviceName, service)

	return &v1beta2.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1beta2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        depName,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: map[string]string{},
		},
		Spec: v1beta2.DeploymentSpec{
			Paused: false,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: podSpec,
			},
		},
	}
}

func (s *stackDeployController) serviceSelector(name, namespace string, labels map[string]string) *v1.Service {
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

func MergeRevisionToService(service *v1beta1.Service, revision string) (*v1beta1.ServiceUnversionedSpec, error) {
	// TODO: do better merging
	newRevision := service.Spec.ServiceUnversionedSpec.DeepCopy()
	serviceRevision, ok := service.Spec.Revisions[revision]
	if !ok {
		return nil, fmt.Errorf("failed to find revision for %s", revision)
	}

	err := convert.ToObj(&serviceRevision.Spec, newRevision)
	return newRevision, err
}

func (s *stackDeployController) service(objects []runtime.Object, name, namespace string, service *v1beta1.Service) ([]runtime.Object, error) {
	objects, err := s.addService(objects, "latest", name, namespace, &service.Spec.ServiceUnversionedSpec)

	for revision := range service.Spec.Revisions {
		newRevision, err := MergeRevisionToService(service, revision)
		if err != nil {
			return nil, err
		}

		objects, err = s.addService(objects, revision, name, namespace, newRevision)
		if err != nil {
			return nil, err
		}
	}

	return objects, err
}

func (s *stackDeployController) nodePorts(objects []runtime.Object, name, namespace string, service *v1beta1.ServiceUnversionedSpec, labels map[string]string) []runtime.Object {
	nodePortService := &v1.Service{
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
			Type:     v1.ServiceTypeNodePort,
			Selector: labels,
		},
	}

	for i, portBinding := range service.PortBindings {
		if portBinding.Port > 0 {
			continue
		}

		servicePort := v1.ServicePort{
			Name:       fmt.Sprintf("port-%d", i),
			TargetPort: intstr.FromInt(int(portBinding.TargetPort)),
			Port:       int32(portBinding.TargetPort),
		}

		switch portBinding.Protocol {
		case "tcp":
			servicePort.Protocol = v1.ProtocolTCP
		case "udp":
			servicePort.Protocol = v1.ProtocolUDP
		default:
			continue
		}

		nodePortService.Spec.Ports = append(nodePortService.Spec.Ports, servicePort)
	}

	if len(nodePortService.Spec.Ports) > 0 {
		objects = append(objects, nodePortService)
	}

	return objects
}

func (s *stackDeployController) addService(objects []runtime.Object, revision, serviceName, namespace string, service *v1beta1.ServiceUnversionedSpec) ([]runtime.Object, error) {
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

	dep := s.deployment(labels, name, serviceName, namespace, service)
	k8sService := s.serviceSelector(name, namespace, labels)
	objects = s.nodePorts(objects, name+"-ports", namespace, service, labels)
	objects = append(objects, dep, k8sService)
	objects = addScale(objects, name, serviceName, namespace, dep, service)

	return objects, nil
}

func ports(podSpec *v1.PodSpec, service *v1beta1.ServiceUnversionedSpec) {
	for _, portBinding := range service.PortBindings {
		if portBinding.Port <= 0 {
			continue
		}

		port := v1.ContainerPort{
			HostPort:      int32(portBinding.Port),
			ContainerPort: int32(portBinding.TargetPort),
			HostIP:        portBinding.IP,
		}

		switch portBinding.Protocol {
		case "tcp":
			port.Protocol = v1.ProtocolTCP
		case "udp":
			port.Protocol = v1.ProtocolUDP
		default:
			continue
		}

		podSpec.Containers[0].Ports = append(podSpec.Containers[0].Ports, port)
	}
}

func dns(podSpec *v1.PodSpec, service *v1beta1.ServiceUnversionedSpec) {
	dnsConfig := &v1.PodDNSConfig{
		Nameservers: service.DNS,
		Searches:    service.DNSSearch,
	}

	for _, dnsOpt := range service.DNSOptions {
		k, v := kv.Split(dnsOpt, "=")
		opt := v1.PodDNSConfigOption{
			Name: k,
		}
		if len(v) > 0 {
			opt.Value = &v
		}
		dnsConfig.Options = append(dnsConfig.Options, opt)
	}

	if len(dnsConfig.Options) > 0 || len(dnsConfig.Searches) > 0 || len(dnsConfig.Nameservers) > 0 {
		podSpec.DNSConfig = dnsConfig
	}

}

func podSpec(serviceName string, service *v1beta1.ServiceUnversionedSpec) v1.PodSpec {
	f := false

	podSpec := v1.PodSpec{
		HostNetwork:                  service.NetworkMode == "host",
		HostIPC:                      service.IpcMode == "host",
		HostPID:                      service.PidMode == "pid",
		Hostname:                     service.Hostname,
		AutomountServiceAccountToken: &f,
	}

	dns(&podSpec, service)

	volumes := map[string]v1.Volume{}

	podSpec.Containers = append(podSpec.Containers, container(serviceName, service.ContainerConfig, volumes))
	for name, sidekick := range service.Sidecars {
		c := container(name, sidekick.ContainerConfig, volumes)
		if sidekick.InitContainer {
			podSpec.InitContainers = append(podSpec.InitContainers, c)
		} else {
			podSpec.Containers = append(podSpec.Containers, c)
		}
	}

	switch service.RestartPolicy {
	case "never":
		podSpec.RestartPolicy = v1.RestartPolicyNever
	case "on-failure":
		podSpec.RestartPolicy = v1.RestartPolicyOnFailure
	case "always":
		podSpec.RestartPolicy = v1.RestartPolicyAlways
	}

	if service.StopGracePeriodSeconds != nil {
		v := int64(*service.StopGracePeriodSeconds)
		podSpec.TerminationGracePeriodSeconds = &v
	}

	for _, host := range service.ExtraHosts {
		ip, host := kv.Split(host, ":")
		podSpec.HostAliases = append(podSpec.HostAliases, v1.HostAlias{
			IP:        ip,
			Hostnames: []string{host},
		})
	}

	for _, volume := range volumes {
		podSpec.Volumes = append(podSpec.Volumes, volume)
	}

	// Must be done after the first container is added
	ports(&podSpec, service)

	return podSpec
}

func container(name string, container v1beta1.ContainerConfig, volumes map[string]v1.Volume) v1.Container {
	c := v1.Container{
		Name:            name,
		Image:           container.Image,
		Command:         container.Entrypoint,
		Args:            container.Command,
		WorkingDir:      container.WorkingDir,
		ImagePullPolicy: v1.PullIfNotPresent,
		SecurityContext: &v1.SecurityContext{
			ReadOnlyRootFilesystem: &container.ReadonlyRootfs,
			Capabilities: &v1.Capabilities{
				Add:  toCaps(container.CapAdd),
				Drop: toCaps(container.CapDrop),
			},
			Privileged: &container.Privileged,
		},
		TTY:       container.Tty,
		StdinOnce: container.OpenStdin,
		Resources: v1.ResourceRequirements{
			Limits:   v1.ResourceList{},
			Requests: v1.ResourceList{},
		},
	}

	populateResources(&c, container)

	if n, err := strconv.ParseInt(container.User, 10, 0); err == nil {
		c.SecurityContext.RunAsUser = &n
	}

	for _, env := range container.Environment {
		name, value := kv.Split(env, "=")
		c.Env = append(c.Env, v1.EnvVar{
			Name:  name,
			Value: value,
		})
	}

	c.LivenessProbe, c.ReadinessProbe = toProbes(container.Healthcheck)

	for i, volume := range container.Volumes {
		name := ""
		if volume.Kind == "bind" {
			name = "host-" + strings.Replace(volume.Source, "/", "-", -1)
			volumes[name] = v1.Volume{
				VolumeSource: v1.VolumeSource{
					HostPath: &v1.HostPathVolumeSource{
						Path: volume.Source,
					},
				},
				Name: name,
			}
		} else if volume.Kind == "volume" && volume.Source != "" {
			name = volume.Source
			volumes[name] = v1.Volume{
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: name,
						ReadOnly:  volume.ReadOnly,
					},
				},
				Name: name,
			}
		}

		if name == "" {
			continue
		}

		mount := v1.VolumeMount{
			Name:      fmt.Sprintf("mount%d", i),
			ReadOnly:  volume.ReadOnly,
			MountPath: volume.Target,
		}

		if volume.BindOptions != nil {
			if strings.Contains(string(volume.BindOptions.Propagation), "shared") {
				prop := v1.MountPropagationBidirectional
				mount.MountPropagation = &prop
			} else if strings.Contains(string(volume.BindOptions.Propagation), "private") ||
				strings.Contains(string(volume.BindOptions.Propagation), "slave") {
				prop := v1.MountPropagationHostToContainer
				mount.MountPropagation = &prop
			}
		}

		if volume.VolumeOptions != nil {
			mount.SubPath = volume.VolumeOptions.SubPath
		}

		c.VolumeMounts = append(c.VolumeMounts, mount)
	}

	return c
}

func populateResources(c *v1.Container, container v1beta1.ContainerConfig) {
	if container.MemoryBytes > 0 {
		c.Resources.Limits[v1.ResourceMemory] = *resource.NewQuantity(container.MemoryBytes, resource.DecimalSI)
	}

	if container.MemoryReservationBytes > 0 {
		c.Resources.Requests[v1.ResourceMemory] = *resource.NewQuantity(container.MemoryBytes, resource.DecimalSI)
	}

	if container.CPUs != "" {
		q, err := resource.ParseQuantity(container.CPUs)
		if err == nil {
			c.Resources.Requests[v1.ResourceCPU] = q
		}
		logrus.Errorf("Failed to parse CPU request: %v", err)
	}
}

func toProbes(healthcheck *v1beta1.HealthConfig) (*v1.Probe, *v1.Probe) {
	if healthcheck == nil {
		return nil, nil
	}

	probe := v1.Probe{
		InitialDelaySeconds: int32(healthcheck.InitialDelaySeconds),
		TimeoutSeconds:      int32(healthcheck.TimeoutSeconds),
		PeriodSeconds:       int32(healthcheck.IntervalSeconds),
		SuccessThreshold:    int32(healthcheck.HealthyThreshold),
		FailureThreshold:    int32(healthcheck.UnhealthyThreshold),
	}

	test := healthcheck.Test[0]
	if strings.HasPrefix(test, "http://") || strings.HasPrefix(test, "https://") {
		u, err := url.Parse(test)
		if err == nil {
			probe.HTTPGet = &v1.HTTPGetAction{
				Path: u.Path,
				Port: intstr.Parse(u.Port()),
			}
			if strings.HasPrefix(test, "http://") {
				probe.HTTPGet.Scheme = v1.URISchemeHTTP
			} else if strings.HasPrefix(test, "https://") {
				probe.HTTPGet.Scheme = v1.URISchemeHTTPS
			}

			for i := 1; i < len(healthcheck.Test); i++ {
				name, value := kv.Split(healthcheck.Test[i], "=")
				probe.HTTPGet.HTTPHeaders = append(probe.HTTPGet.HTTPHeaders, v1.HTTPHeader{
					Name:  name,
					Value: value,
				})
			}
		}
	} else if strings.HasPrefix(test, "tcp://") {
		u, err := url.Parse(test)
		if err == nil {
			probe.TCPSocket = &v1.TCPSocketAction{
				Port: intstr.Parse(u.Port()),
			}
		}
	} else if test == "CMD" {
		probe.Exec = &v1.ExecAction{
			Command: healthcheck.Test[1:],
		}
	} else if test == "CMD-SHELL" {
		if len(healthcheck.Test) == 2 {
			probe.Exec = &v1.ExecAction{
				Command: []string{"sh", "-c", healthcheck.Test[1]},
			}
		}
	} else if test == "NONE" {
		return nil, nil
	} else {
		probe.Exec = &v1.ExecAction{
			Command: healthcheck.Test,
		}
	}

	return &probe, &probe
}

func toCaps(args []string) []v1.Capability {
	var caps []v1.Capability
	for _, arg := range args {
		caps = append(caps, v1.Capability(arg))
	}
	return caps
}
