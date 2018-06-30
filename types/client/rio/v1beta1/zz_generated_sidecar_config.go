package client

const (
	SidecarConfigType                        = "sidecarConfig"
	SidecarConfigFieldCPUs                   = "nanoCpus"
	SidecarConfigFieldCapAdd                 = "capAdd"
	SidecarConfigFieldCapDrop                = "capDrop"
	SidecarConfigFieldCommand                = "command"
	SidecarConfigFieldConfigs                = "configs"
	SidecarConfigFieldDefaultVolumeDriver    = "defaultVolumeDriver"
	SidecarConfigFieldDevices                = "devices"
	SidecarConfigFieldEntrypoint             = "entrypoint"
	SidecarConfigFieldEnvironment            = "environment"
	SidecarConfigFieldHealthcheck            = "healthcheck"
	SidecarConfigFieldImage                  = "image"
	SidecarConfigFieldInit                   = "init"
	SidecarConfigFieldInitContainer          = "initContainer"
	SidecarConfigFieldMemoryBytes            = "memoryBytes"
	SidecarConfigFieldMemoryReservationBytes = "memoryReservationBytes"
	SidecarConfigFieldOpenStdin              = "stdinOpen"
	SidecarConfigFieldPrivileged             = "privileged"
	SidecarConfigFieldReadonlyRootfs         = "readOnly"
	SidecarConfigFieldTmpfs                  = "tmpfs"
	SidecarConfigFieldTty                    = "tty"
	SidecarConfigFieldUser                   = "user"
	SidecarConfigFieldVolumes                = "volumes"
	SidecarConfigFieldVolumesFrom            = "volumesFrom"
	SidecarConfigFieldWorkingDir             = "workingDir"
)

type SidecarConfig struct {
	CPUs                   string          `json:"nanoCpus,omitempty" yaml:"nanoCpus,omitempty"`
	CapAdd                 []string        `json:"capAdd,omitempty" yaml:"capAdd,omitempty"`
	CapDrop                []string        `json:"capDrop,omitempty" yaml:"capDrop,omitempty"`
	Command                []string        `json:"command,omitempty" yaml:"command,omitempty"`
	Configs                []ConfigMapping `json:"configs,omitempty" yaml:"configs,omitempty"`
	DefaultVolumeDriver    string          `json:"defaultVolumeDriver,omitempty" yaml:"defaultVolumeDriver,omitempty"`
	Devices                []DeviceMapping `json:"devices,omitempty" yaml:"devices,omitempty"`
	Entrypoint             []string        `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Environment            []string        `json:"environment,omitempty" yaml:"environment,omitempty"`
	Healthcheck            *HealthConfig   `json:"healthcheck,omitempty" yaml:"healthcheck,omitempty"`
	Image                  string          `json:"image,omitempty" yaml:"image,omitempty"`
	Init                   bool            `json:"init,omitempty" yaml:"init,omitempty"`
	InitContainer          bool            `json:"initContainer,omitempty" yaml:"initContainer,omitempty"`
	MemoryBytes            int64           `json:"memoryBytes,omitempty" yaml:"memoryBytes,omitempty"`
	MemoryReservationBytes int64           `json:"memoryReservationBytes,omitempty" yaml:"memoryReservationBytes,omitempty"`
	OpenStdin              bool            `json:"stdinOpen,omitempty" yaml:"stdinOpen,omitempty"`
	Privileged             bool            `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	ReadonlyRootfs         bool            `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Tmpfs                  []Tmpfs         `json:"tmpfs,omitempty" yaml:"tmpfs,omitempty"`
	Tty                    bool            `json:"tty,omitempty" yaml:"tty,omitempty"`
	User                   string          `json:"user,omitempty" yaml:"user,omitempty"`
	Volumes                []Mount         `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	VolumesFrom            []string        `json:"volumesFrom,omitempty" yaml:"volumesFrom,omitempty"`
	WorkingDir             string          `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
}
