package stackdeploy

import (
	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
)

func podSpec(serviceName string, serviceLabels map[string]string, service *v1beta1.ServiceUnversionedSpec, volumeDefs map[string]*v1beta1.Volume) (map[string]*v1beta1.Volume, v1.PodSpec) {
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
	usedTemplates := map[string]*v1beta1.Volume{}

	podSpec.Containers = append(podSpec.Containers, container(serviceName, service.ContainerConfig, volumes, volumeDefs, usedTemplates))
	for name, sidekick := range service.Sidecars {
		c := container(name, sidekick.ContainerConfig, volumes, volumeDefs, usedTemplates)
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

	scheduling(&podSpec, service, serviceLabels)

	// Must be done after the first container is added
	ports(&podSpec, service)

	return usedTemplates, podSpec
}

func scheduling(podSpec *v1.PodSpec, service *v1beta1.ServiceUnversionedSpec, serviceLabels map[string]string) {
	nodeAffinity, err := service.Scheduling.ToNodeAffinity()
	if err == nil {
		podSpec.Affinity = &v1.Affinity{
			NodeAffinity: nodeAffinity,
		}
	} else {
		logrus.Errorf("failed to parse scheduling for service: %v", err)
	}

	podSpec.SchedulerName = service.Scheduling.Scheduler

	// mergeLabels will strip out rio.cattle.io labels
	for k, v := range mergeLabels(nil, serviceLabels) {
		toleration := v1.Toleration{
			Key:      k,
			Operator: v1.TolerationOpExists,
			Value:    v,
		}

		if len(toleration.Value) > 0 {
			toleration.Operator = v1.TolerationOpEqual
		}

		toleration.Effect = v1.TaintEffectNoExecute
		podSpec.Tolerations = append(podSpec.Tolerations)
		toleration.Effect = v1.TaintEffectNoSchedule
		podSpec.Tolerations = append(podSpec.Tolerations)
		toleration.Effect = v1.TaintEffectPreferNoSchedule
		podSpec.Tolerations = append(podSpec.Tolerations)
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
