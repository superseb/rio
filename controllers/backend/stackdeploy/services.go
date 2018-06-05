package stackdeploy

import (
	"net/url"
	"strconv"
	"strings"

	"fmt"

	"github.com/pkg/errors"
	"github.com/rancher/rio/cli2/pkg/kv"
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
		objects, err = s.service(objects, service)
		if err != nil {
			return objects, errors.Wrapf(err, "failed to construct service for %s/%s", service.Namespace, service.Name)
		}

	}

	return objects, nil
}

func (s *stackDeployController) service(objects []runtime.Object, service *v1beta1.Service) ([]runtime.Object, error) {
	labels := map[string]string{
		"service.cattle.io": service.Name,
	}

	scale := int32(service.Scale)
	surge := service.BatchSize
	unavailable := 0

	if service.UpdateOrder == "start-first" {
		surge = service.BatchSize
		unavailable = 0
	} else if service.UpdateOrder == "stop-first" {
		surge = 0
		unavailable = service.BatchSize
	}

	maxSurge := intstr.FromInt(surge)
	maxUnavailable := intstr.FromInt(unavailable)

	dep := &v1beta2.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1beta2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
		Spec: v1beta2.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: podSpec(service),
			},
			Replicas: &scale,
			Strategy: v1beta2.DeploymentStrategy{
				RollingUpdate: &v1beta2.RollingUpdateDeployment{
					MaxSurge:       &maxSurge,
					MaxUnavailable: &maxUnavailable,
				},
				Type: v1beta2.RollingUpdateDeploymentStrategyType,
			},
		},
	}

	k8sService := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: v1.ClusterIPNone,
			Selector:  labels,
		},
	}

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
			Name:      service.Name,
			Namespace: service.Namespace,
		},
		Spec: v1beta12.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			MaxUnavailable: &pdbQuantity,
		},
	}

	return append(objects, dep, k8sService, pdb), nil
}

func podSpec(service *v1beta1.Service) v1.PodSpec {
	podSpec := v1.PodSpec{
		HostNetwork: service.NetworkMode == "host",
		HostIPC:     service.IpcMode == "host",
		HostPID:     service.PidMode == "pid",
		Hostname:    service.Hostname,
	}

	volumes := map[string]v1.Volume{}

	podSpec.Containers = append(podSpec.Containers, container(service.Name, service.ContainerConfig, volumes))
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

	return podSpec
}

func container(name string, container v1beta1.ContainerConfig, volumes map[string]v1.Volume) v1.Container {
	c := v1.Container{
		Name:       name,
		Image:      container.Image,
		Command:    container.Entrypoint,
		Args:       container.Command,
		WorkingDir: container.WorkingDir,
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
