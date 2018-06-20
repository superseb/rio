package client

const (
	ServiceUnversionedSpecType                        = "serviceUnversionedSpec"
	ServiceUnversionedSpecFieldBatchSize              = "batchSize"
	ServiceUnversionedSpecFieldCPUs                   = "nanoCpus"
	ServiceUnversionedSpecFieldCapAdd                 = "capAdd"
	ServiceUnversionedSpecFieldCapDrop                = "capDrop"
	ServiceUnversionedSpecFieldCommand                = "command"
	ServiceUnversionedSpecFieldDNS                    = "dns"
	ServiceUnversionedSpecFieldDNSOptions             = "dnsOptions"
	ServiceUnversionedSpecFieldDNSSearch              = "dnsSearch"
	ServiceUnversionedSpecFieldDefaultVolumeDriver    = "defaultVolumeDriver"
	ServiceUnversionedSpecFieldDevices                = "devices"
	ServiceUnversionedSpecFieldEntrypoint             = "entrypoint"
	ServiceUnversionedSpecFieldEnvironment            = "environment"
	ServiceUnversionedSpecFieldExtraHosts             = "extraHosts"
	ServiceUnversionedSpecFieldHealthcheck            = "healthcheck"
	ServiceUnversionedSpecFieldHostname               = "hostname"
	ServiceUnversionedSpecFieldImage                  = "image"
	ServiceUnversionedSpecFieldInit                   = "init"
	ServiceUnversionedSpecFieldIpcMode                = "ipc"
	ServiceUnversionedSpecFieldMemoryBytes            = "memoryBytes"
	ServiceUnversionedSpecFieldMemoryReservationBytes = "memoryReservationBytes"
	ServiceUnversionedSpecFieldNetworkMode            = "net"
	ServiceUnversionedSpecFieldOpenStdin              = "stdinOpen"
	ServiceUnversionedSpecFieldPidMode                = "pid"
	ServiceUnversionedSpecFieldPortBindings           = "ports"
	ServiceUnversionedSpecFieldPrivileged             = "privileged"
	ServiceUnversionedSpecFieldReadonlyRootfs         = "readOnly"
	ServiceUnversionedSpecFieldRestartPolicy          = "restart"
	ServiceUnversionedSpecFieldScale                  = "scale"
	ServiceUnversionedSpecFieldSidecars               = "sidecars"
	ServiceUnversionedSpecFieldSpaceID                = "spaceId"
	ServiceUnversionedSpecFieldStackID                = "stackId"
	ServiceUnversionedSpecFieldStopGracePeriodSeconds = "stopGracePeriod"
	ServiceUnversionedSpecFieldTmpfs                  = "tmpfs"
	ServiceUnversionedSpecFieldTty                    = "tty"
	ServiceUnversionedSpecFieldUpdateOrder            = "updateOrder"
	ServiceUnversionedSpecFieldUser                   = "user"
	ServiceUnversionedSpecFieldVolumes                = "volumes"
	ServiceUnversionedSpecFieldVolumesFrom            = "volumesFrom"
	ServiceUnversionedSpecFieldWorkingDir             = "workingDir"
)

type ServiceUnversionedSpec struct {
	BatchSize              int64                    `json:"batchSize,omitempty" yaml:"batchSize,omitempty"`
	CPUs                   string                   `json:"nanoCpus,omitempty" yaml:"nanoCpus,omitempty"`
	CapAdd                 []string                 `json:"capAdd,omitempty" yaml:"capAdd,omitempty"`
	CapDrop                []string                 `json:"capDrop,omitempty" yaml:"capDrop,omitempty"`
	Command                []string                 `json:"command,omitempty" yaml:"command,omitempty"`
	DNS                    []string                 `json:"dns,omitempty" yaml:"dns,omitempty"`
	DNSOptions             []string                 `json:"dnsOptions,omitempty" yaml:"dnsOptions,omitempty"`
	DNSSearch              []string                 `json:"dnsSearch,omitempty" yaml:"dnsSearch,omitempty"`
	DefaultVolumeDriver    string                   `json:"defaultVolumeDriver,omitempty" yaml:"defaultVolumeDriver,omitempty"`
	Devices                []DeviceMapping          `json:"devices,omitempty" yaml:"devices,omitempty"`
	Entrypoint             []string                 `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Environment            []string                 `json:"environment,omitempty" yaml:"environment,omitempty"`
	ExtraHosts             []string                 `json:"extraHosts,omitempty" yaml:"extraHosts,omitempty"`
	Healthcheck            *HealthConfig            `json:"healthcheck,omitempty" yaml:"healthcheck,omitempty"`
	Hostname               string                   `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Image                  string                   `json:"image,omitempty" yaml:"image,omitempty"`
	Init                   bool                     `json:"init,omitempty" yaml:"init,omitempty"`
	IpcMode                string                   `json:"ipc,omitempty" yaml:"ipc,omitempty"`
	MemoryBytes            int64                    `json:"memoryBytes,omitempty" yaml:"memoryBytes,omitempty"`
	MemoryReservationBytes int64                    `json:"memoryReservationBytes,omitempty" yaml:"memoryReservationBytes,omitempty"`
	NetworkMode            string                   `json:"net,omitempty" yaml:"net,omitempty"`
	OpenStdin              bool                     `json:"stdinOpen,omitempty" yaml:"stdinOpen,omitempty"`
	PidMode                string                   `json:"pid,omitempty" yaml:"pid,omitempty"`
	PortBindings           []PortBinding            `json:"ports,omitempty" yaml:"ports,omitempty"`
	Privileged             bool                     `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	ReadonlyRootfs         bool                     `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	RestartPolicy          string                   `json:"restart,omitempty" yaml:"restart,omitempty"`
	Scale                  int64                    `json:"scale,omitempty" yaml:"scale,omitempty"`
	Sidecars               map[string]SidecarConfig `json:"sidecars,omitempty" yaml:"sidecars,omitempty"`
	SpaceID                string                   `json:"spaceId,omitempty" yaml:"spaceId,omitempty"`
	StackID                string                   `json:"stackId,omitempty" yaml:"stackId,omitempty"`
	StopGracePeriodSeconds *int64                   `json:"stopGracePeriod,omitempty" yaml:"stopGracePeriod,omitempty"`
	Tmpfs                  []Tmpfs                  `json:"tmpfs,omitempty" yaml:"tmpfs,omitempty"`
	Tty                    bool                     `json:"tty,omitempty" yaml:"tty,omitempty"`
	UpdateOrder            string                   `json:"updateOrder,omitempty" yaml:"updateOrder,omitempty"`
	User                   string                   `json:"user,omitempty" yaml:"user,omitempty"`
	Volumes                []Mount                  `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	VolumesFrom            []string                 `json:"volumesFrom,omitempty" yaml:"volumesFrom,omitempty"`
	WorkingDir             string                   `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
}
