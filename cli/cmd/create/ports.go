package create

import (
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func ParsePorts(specs []string) ([]client.PortBinding, error) {
	var (
		result []client.PortBinding
		err    error
	)

	_, portBindings, err := nat.ParsePortSpecs(specs)
	if err != nil {
		return nil, err
	}

	for port, bindings := range portBindings {
		proto, port := nat.SplitProtoPort(string(port))
		for _, binding := range bindings {
			pb := client.PortBinding{
				Protocol: proto,
				IP:       binding.HostIP,
			}

			if portNum, err := strconv.Atoi(port); err == nil {
				pb.TargetPort = int64(portNum)
			}

			if portNum, err := strconv.Atoi(binding.HostPort); err == nil {
				pb.Port = int64(portNum)
			}

			result = append(result, pb)
		}
	}

	return result, nil
}
