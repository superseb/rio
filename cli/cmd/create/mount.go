package create

import (
	"fmt"
	"strings"

	"github.com/docker/cli/cli/compose/loader"
	"github.com/docker/cli/cli/compose/types"
	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func ParseMounts(spec []string) ([]client.Mount, error) {
	var mounts []client.Mount
	for _, volume := range spec {
		serviceMount, err := loader.ParseVolume(volume)
		if err != nil {
			return nil, err
		}

		mount := createMount(serviceMount)
		mount, err = parseAdditionalOptions(mount, volume)
		if err != nil {
			return nil, err
		}

		mounts = append(mounts, createMount(serviceMount))
	}

	return mounts, nil
}

func parseAdditionalOptions(mount client.Mount, spec string) (client.Mount, error) {
	if mount.Source == "" {
		return mount, fmt.Errorf("invalid volume spec, no target path found: %s", spec)
	}

	if !strings.Contains(mount.Source, "/") && !strings.Contains(mount.Source, "\\") && mount.VolumeOptions == nil {
		mount.VolumeOptions = &client.VolumeOptions{}
	}

	if len(strings.SplitN(spec, ":", 3)) < 3 {
		return mount, nil
	}

	setDriverOpts := false
	driverOpts := &client.DriverConfig{
		Options: map[string]string{},
	}

	for _, opt := range strings.Split(spec, ",") {
		key, value := kv.Split(opt, "=")
		key = strings.ToLower(key)
		switch key {
		case "driver":
			setDriverOpts = true
			driverOpts.Name = value
		case "driver-opt":
			setDriverOpts = true
			key, value := kv.Split(value, "=")
			driverOpts.Options[key] = value
		case "subpath":
			if mount.VolumeOptions == nil {
				return mount, fmt.Errorf("subpath can only be used with volumes, not host bind mounts")
			}
			mount.VolumeOptions.SubPath = value
		}
	}

	if setDriverOpts {
		if mount.VolumeOptions == nil {
			return mount, fmt.Errorf("driver and driver-opts can only be used with volumes, not host bind mounts")
		}
		mount.VolumeOptions.DriverConfig = driverOpts
	}

	return mount, nil
}

func createMount(serviceMount types.ServiceVolumeConfig) client.Mount {
	mount := client.Mount{
		Kind:     serviceMount.Type,
		ReadOnly: serviceMount.ReadOnly,
		Source:   serviceMount.Source,
		Target:   serviceMount.Target,
	}

	if serviceMount.Bind != nil {
		mount.BindOptions = &client.BindOptions{
			Propagation: serviceMount.Bind.Propagation,
		}
	}

	if serviceMount.Volume != nil {
		mount.VolumeOptions = &client.VolumeOptions{
			NoCopy: serviceMount.Volume.NoCopy,
		}
	}

	return mount
}
