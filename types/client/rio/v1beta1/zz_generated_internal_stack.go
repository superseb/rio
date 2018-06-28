package client

const (
	InternalStackType          = "internalStack"
	InternalStackFieldConfigs  = "configs"
	InternalStackFieldServices = "services"
	InternalStackFieldVolumes  = "volumes"
)

type InternalStack struct {
	Configs  map[string]Config  `json:"configs,omitempty" yaml:"configs,omitempty"`
	Services map[string]Service `json:"services,omitempty" yaml:"services,omitempty"`
	Volumes  map[string]Volume  `json:"volumes,omitempty" yaml:"volumes,omitempty"`
}
