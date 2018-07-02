package stackdeploy

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	fieldRefs = map[string]string{
		"self/name":           "metadata.name",
		"self/namespace":      "metadata.namespace",
		"self/labels":         "metadata.labels",
		"self/annotations":    "metadata.annotations",
		"self/node":           "spec.nodeName",
		"self/serviceAccount": "spec.serviceAccountName",
		"self/hostIp":         "status.hostIP",
		"self/ip":             "status.podIP",
	}
	resourceRefs = map[string]string{
		"limits/cpu":                 "limits.cpu",
		"limits/memory":              "limits.memory",
		"limits/ephemeral-storage":   "limits.ephemeral-storage",
		"requests/cpu":               "requests.cpu",
		"requests/memory":            "requests.memory",
		"requests/ephemeral-storage": "requests.ephemeral-storage",
	}
)

func container(name string, container v1beta1.ContainerConfig, volumes map[string]v1.Volume, volumeDefs map[string]*v1beta1.Volume, usedTemplates map[string]*v1beta1.Volume) v1.Container {
	c := v1.Container{
		Name:            name,
		Image:           container.Image,
		Command:         container.Entrypoint,
		Args:            container.Command,
		WorkingDir:      container.WorkingDir,
		ImagePullPolicy: v1.PullPolicy(container.ImagePullPolicy),
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

	populateEnv(&c, container)

	c.LivenessProbe, c.ReadinessProbe = toProbes(container.Healthcheck)

	for _, volume := range container.Volumes {
		addVolumes(&c, volume, volumes, volumeDefs, usedTemplates)
	}

	return c
}

func populateEnv(c *v1.Container, container v1beta1.ContainerConfig) {
	for _, env := range container.Environment {
		name, value := kv.Split(env, "=")
		c.Env = append(c.Env, toEnvVar(c.Name, name, value))
	}
}

func toEnvVar(containerName, name, value string) v1.EnvVar {
	basic := v1.EnvVar{
		Name:  name,
		Value: value,
	}

	if !strings.HasPrefix(value, "$(") || !strings.HasSuffix(value, ")") {
		return basic
	}

	key := value[2 : len(value)-1]

	if fieldRefValue, ok := fieldRefs[key]; ok {
		return v1.EnvVar{
			Name: name,
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: fieldRefValue,
				},
			},
		}
	}

	if resourceRefValue, ok := resourceRefs[key]; ok {
		return v1.EnvVar{
			Name: name,
			ValueFrom: &v1.EnvVarSource{
				ResourceFieldRef: &v1.ResourceFieldSelector{
					ContainerName: containerName,
					Resource:      resourceRefValue,
				},
			},
		}
	}

	k, v := kv.Split(key, "/")
	optional := strings.HasSuffix(v, "?")
	if optional {
		v = v[:len(v)-1]
	}

	if v == "" {
		return basic
	}

	switch k {
	case "config":
		return v1.EnvVar{
			Name: name,
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: v,
					},
					Key:      "content",
					Optional: &optional,
				},
			},
		}
	case "secret":
		return v1.EnvVar{
			Name: name,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: v,
					},
					Key:      "content",
					Optional: &optional,
				},
			},
		}
	default:
		if resourceRefValue, ok := resourceRefs[v]; ok {
			return v1.EnvVar{
				Name: name,
				ValueFrom: &v1.EnvVarSource{
					ResourceFieldRef: &v1.ResourceFieldSelector{
						ContainerName: k,
						Resource:      resourceRefValue,
					},
				},
			}
		}
	}

	return basic
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
