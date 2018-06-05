package v1beta1

import (
	"strconv"

	"bytes"

	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	ServiceSpec `json:"spec"`
}

type ServiceSpec struct {
	Scale       int    `json:"scale,omitempty"`
	BatchSize   int    `json:"batchSize,omitempty"`
	UpdateOrder string `json:"updateOrder,omitempty" norman:"type=enum,options=start-first|stop-first"`

	StackScoped
	PodConfig
	PrivilegedConfig
	Sidecars map[string]SidecarConfig `json:"sidecars,omitempty"`

	ContainerConfig
}

type PodConfig struct {
	Hostname string `json:"hostname,omitempty"`
	//Metadata               map[string]interface{}   `json:"metadata,omitempty"`        //alias annotations
	//StopSignal             string   `json:"stopSignal,omitempty" norman:"default=SIGTERM"`                                       // Signal to stop a container
	StopGracePeriodSeconds *int   `json:"stopGracePeriod,omitempty"`                                                           // support friendly numbers
	RestartPolicy          string `json:"restart,omitempty" norman:"type=enum,options=never|on-failure|always,default=always"` //support no and OnFailure
	//DNS                    []string `json:"dns,omitempty"`                                                                       // support string
	//DNSOptions             []string `json:"dnsOptions,omitempty"`                                                                // support string
	//DNSSearch              []string `json:"dnsSearch,omitempty"`                                                                 // support string
	ExtraHosts []string `json:"extraHosts,omitempty"` // support map
}

type PrivilegedConfig struct {
	NetworkMode  string        `json:"net,json" norman:"type=enum,options=default|host,default=default"` // alias network, support bridge
	PortBindings []PortBinding `json:"ports,omitempty"`                                                  // support []string
	IpcMode      string        `json:"ipc,omitempty" norman:"type=enum,options=default|host,default=default"`
	PidMode      string        `json:"pid,omitempty" norman:"type=enum,options=default|host,default=default"`
}

type ContainerPrivilegedConfig struct {
	Privileged bool `json:"privileged,omitempty"`
}

// PortBinding represents a binding between a Host IP address and a Host Port
type PortBinding struct {
	Port       int64  `json:"port"`
	Protocol   string `json:"protocol"`
	IP         string `json:"ip"`
	TargetPort int64  `json:"targetPort"`
}

func (p PortBinding) String() string {
	b := bytes.Buffer{}
	if p.Port != 0 && p.TargetPort != 0 {
		if p.IP != "" {
			b.WriteString(p.IP)
			b.WriteString(":")
		}
		b.WriteString(strconv.FormatInt(p.Port, 10))
		b.WriteString(":")
		b.WriteString(strconv.FormatInt(p.TargetPort, 10))
	} else if p.TargetPort != 0 {
		b.WriteString(strconv.FormatInt(p.TargetPort, 10))
	}

	if b.Len() > 0 && p.Protocol != "" {
		b.WriteString("/")
		b.WriteString(p.Protocol)
	}

	return b.String()
}

// TODO: add pull policy
type ContainerConfig struct {
	ContainerPrivilegedConfig

	ReadonlyRootfs         bool          `json:"readOnly,omitempty"`
	CapAdd                 []string      `json:"capAdd,omitempty"`  // support string
	CapDrop                []string      `json:"capDrop,omitempty"` // support string
	User                   string        `json:"user,omitempty"`
	Tty                    bool          `json:"tty,omitempty"`
	OpenStdin              bool          `json:"stdinOpen,omitempty"`   // alias interactive
	Environment            []string      `json:"environment,omitempty"` // alias env, support map
	Entrypoint             []string      `json:"entrypoint,omitempty"`
	Command                []string      `json:"command,omitempty"` // support string
	Image                  string        `json:"image,omitempty"`
	Init                   bool          `json:"init,omitempty"`
	Healthcheck            *HealthConfig `json:"healthcheck,omitempty"`
	Tmpfs                  []Tmpfs       `json:"tmpfs,omitempty"` // support []string too
	DefaultVolumeDriver    string        `json:"defaultVolumeDriver,omitempty"`
	Volumes                []Mount       `json:"volumes,omitempty"`     // support []string too
	VolumesFrom            []string      `json:"volumesFrom,omitempty"` // support []string too
	WorkingDir             string        `json:"workingDir,omitempty"`
	MemoryBytes            int64         `json:"memoryBytes,omitempty"`
	MemoryReservationBytes int64         `json:"memoryReservationBytes,omitempty"`
	CPUs                   string        `json:"nanoCpus,omitempty"`

	Devices []DeviceMapping `json:"devices,omitempty"` // support []string and map[string]string
}

type SidecarConfig struct {
	InitContainer bool `json:",omitempty"`
	ContainerConfig
}

type HealthConfig struct {
	// Test is the test to perform to check that the container is healthy.
	// An empty slice means to inherit the default.
	// The options are:
	// {} : inherit healthcheck
	// {"NONE"} : disable healthcheck
	// {"CMD", args...} : exec arguments directly
	// {"CMD-SHELL", command} : run command with system's default shell
	Test []string `json:"test,omitempty"` //alias string, deal with CMD, CMD-SHELL, NONE

	IntervalSeconds     int `json:"intervalSeconds,omitempty"`     // support friendly numbers, alias periodSeconds, period
	TimeoutSeconds      int `json:"timeoutSeconds,omitempty"`      // support friendly numbers
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"` //alias start_period
	HealthyThreshold    int `json:"healthyThreshold,omitempty"`    //alias retries, successThreshold
	UnhealthyThreshold  int `json:"unhealthyThreshold,omitempty"`  //alias failureThreshold, set to retries if unset
}

// DeviceMapping represents the device mapping between the host and the container.
type DeviceMapping struct {
	OnHost      string `json:"onHost"`
	InContainer string `json:"inContainer"`
	Permissions string `json:"permissions"`
}

func (d DeviceMapping) String() string {
	result := d.OnHost
	if len(d.InContainer) > 0 {
		if len(result) > 0 {
			result += ":"
		}
		result += d.InContainer
	}
	if len(d.Permissions) > 0 {
		if len(result) > 0 {
			result += ":"
		}
		result += d.Permissions
	}

	return result
}
