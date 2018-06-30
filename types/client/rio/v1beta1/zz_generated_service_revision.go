package client

const (
	ServiceRevisionType                        = "serviceRevision"
	ServiceRevisionFieldBatchSize              = "batchSize"
	ServiceRevisionFieldCPUs                   = "nanoCpus"
	ServiceRevisionFieldCapAdd                 = "capAdd"
	ServiceRevisionFieldCapDrop                = "capDrop"
	ServiceRevisionFieldCommand                = "command"
	ServiceRevisionFieldConditions             = "conditions"
	ServiceRevisionFieldConfigs                = "configs"
	ServiceRevisionFieldDNS                    = "dns"
	ServiceRevisionFieldDNSOptions             = "dnsOptions"
	ServiceRevisionFieldDNSSearch              = "dnsSearch"
	ServiceRevisionFieldDefaultVolumeDriver    = "defaultVolumeDriver"
	ServiceRevisionFieldDevices                = "devices"
	ServiceRevisionFieldEntrypoint             = "entrypoint"
	ServiceRevisionFieldEnvironment            = "environment"
	ServiceRevisionFieldExtraHosts             = "extraHosts"
	ServiceRevisionFieldHealthcheck            = "healthcheck"
	ServiceRevisionFieldHostname               = "hostname"
	ServiceRevisionFieldImage                  = "image"
	ServiceRevisionFieldInit                   = "init"
	ServiceRevisionFieldIpcMode                = "ipc"
	ServiceRevisionFieldLabels                 = "labels"
	ServiceRevisionFieldMemoryBytes            = "memoryBytes"
	ServiceRevisionFieldMemoryReservationBytes = "memoryReservationBytes"
	ServiceRevisionFieldNetworkMode            = "net"
	ServiceRevisionFieldOpenStdin              = "stdinOpen"
	ServiceRevisionFieldPidMode                = "pid"
	ServiceRevisionFieldPortBindings           = "ports"
	ServiceRevisionFieldPrivileged             = "privileged"
	ServiceRevisionFieldPromote                = "promote"
	ServiceRevisionFieldReadonlyRootfs         = "readOnly"
	ServiceRevisionFieldRestartPolicy          = "restart"
	ServiceRevisionFieldScale                  = "scale"
	ServiceRevisionFieldScaleStatus            = "scaleStatus"
	ServiceRevisionFieldSidecars               = "sidecars"
	ServiceRevisionFieldState                  = "state"
	ServiceRevisionFieldStopGracePeriodSeconds = "stopGracePeriod"
	ServiceRevisionFieldTmpfs                  = "tmpfs"
	ServiceRevisionFieldTransitioning          = "transitioning"
	ServiceRevisionFieldTransitioningMessage   = "transitioningMessage"
	ServiceRevisionFieldTty                    = "tty"
	ServiceRevisionFieldUpdateOrder            = "updateOrder"
	ServiceRevisionFieldUser                   = "user"
	ServiceRevisionFieldVolumes                = "volumes"
	ServiceRevisionFieldVolumesFrom            = "volumesFrom"
	ServiceRevisionFieldWeight                 = "weight"
	ServiceRevisionFieldWorkingDir             = "workingDir"
)

type ServiceRevision struct {
	BatchSize              int64                    `json:"batchSize,omitempty" yaml:"batchSize,omitempty"`
	CPUs                   string                   `json:"nanoCpus,omitempty" yaml:"nanoCpus,omitempty"`
	CapAdd                 []string                 `json:"capAdd,omitempty" yaml:"capAdd,omitempty"`
	CapDrop                []string                 `json:"capDrop,omitempty" yaml:"capDrop,omitempty"`
	Command                []string                 `json:"command,omitempty" yaml:"command,omitempty"`
	Conditions             []Condition              `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Configs                []ConfigMapping          `json:"configs,omitempty" yaml:"configs,omitempty"`
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
	Labels                 map[string]string        `json:"labels,omitempty" yaml:"labels,omitempty"`
	MemoryBytes            int64                    `json:"memoryBytes,omitempty" yaml:"memoryBytes,omitempty"`
	MemoryReservationBytes int64                    `json:"memoryReservationBytes,omitempty" yaml:"memoryReservationBytes,omitempty"`
	NetworkMode            string                   `json:"net,omitempty" yaml:"net,omitempty"`
	OpenStdin              bool                     `json:"stdinOpen,omitempty" yaml:"stdinOpen,omitempty"`
	PidMode                string                   `json:"pid,omitempty" yaml:"pid,omitempty"`
	PortBindings           []PortBinding            `json:"ports,omitempty" yaml:"ports,omitempty"`
	Privileged             bool                     `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	Promote                bool                     `json:"promote,omitempty" yaml:"promote,omitempty"`
	ReadonlyRootfs         bool                     `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	RestartPolicy          string                   `json:"restart,omitempty" yaml:"restart,omitempty"`
	Scale                  int64                    `json:"scale,omitempty" yaml:"scale,omitempty"`
	ScaleStatus            *ScaleStatus             `json:"scaleStatus,omitempty" yaml:"scaleStatus,omitempty"`
	Sidecars               map[string]SidecarConfig `json:"sidecars,omitempty" yaml:"sidecars,omitempty"`
	State                  string                   `json:"state,omitempty" yaml:"state,omitempty"`
	StopGracePeriodSeconds *int64                   `json:"stopGracePeriod,omitempty" yaml:"stopGracePeriod,omitempty"`
	Tmpfs                  []Tmpfs                  `json:"tmpfs,omitempty" yaml:"tmpfs,omitempty"`
	Transitioning          string                   `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage   string                   `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	Tty                    bool                     `json:"tty,omitempty" yaml:"tty,omitempty"`
	UpdateOrder            string                   `json:"updateOrder,omitempty" yaml:"updateOrder,omitempty"`
	User                   string                   `json:"user,omitempty" yaml:"user,omitempty"`
	Volumes                []Mount                  `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	VolumesFrom            []string                 `json:"volumesFrom,omitempty" yaml:"volumesFrom,omitempty"`
	Weight                 int64                    `json:"weight,omitempty" yaml:"weight,omitempty"`
	WorkingDir             string                   `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
}
