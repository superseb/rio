package v1beta1

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/go-units"
)

type Mount struct {
	Kind string `json:"kind,omitempty" norman:"type=enum,options=bind|volume"`
	// Source specifies the name of the mount. Depending on mount type, this
	// may be a volume name or a host path, or even ignored.
	Source   string `json:"source,omitempty"`
	Target   string `json:"target,omitempty"`
	ReadOnly bool   `json:"readonly,omitempty"`

	BindOptions   *BindOptions   `json:"bind,omitempty"`
	VolumeOptions *VolumeOptions `json:"volume,omitempty"`
}

func (m Mount) String() string {
	result := ""
	if m.Source != "" {
		result = m.Source + ":"
	}
	result += m.Target

	opts := ""
	if m.ReadOnly {
		addOpt(opts, "ro")
	}
	if m.BindOptions != nil {
		addOpt(opts, string(m.BindOptions.Propagation))
	}
	if m.VolumeOptions != nil {
		if m.VolumeOptions.NoCopy {
			addOpt(opts, "nocopy")
		}
		if m.VolumeOptions.SubPath != "" {
			addOpt(opts, "subPath="+m.VolumeOptions.SubPath)
		}
		if m.VolumeOptions.DriverConfig != nil && m.VolumeOptions.DriverConfig.Name != "" {
			addOpt(opts, m.VolumeOptions.DriverConfig.String())
		}
	}

	if len(opts) == 0 {
		return result
	}

	return result + ":" + opts
}

func addOpt(opt, val string) string {
	if val == "" {
		return opt
	}

	if len(opt) == 0 {
		opt = ":"
	} else {
		opt += ","
	}
	return opt + val
}

// Propagation represents the propagation of a mount.
type Propagation string

const (
	// PropagationRPrivate RPRIVATE
	PropagationRPrivate Propagation = "rprivate"
	// PropagationPrivate PRIVATE
	PropagationPrivate Propagation = "private"
	// PropagationRShared RSHARED
	PropagationRShared Propagation = "rshared"
	// PropagationShared SHARED
	PropagationShared Propagation = "shared"
	// PropagationRSlave RSLAVE
	PropagationRSlave Propagation = "rslave"
	// PropagationSlave SLAVE
	PropagationSlave Propagation = "slave"
)

// Propagations is the list of all valid mount propagations
var Propagations = []Propagation{
	PropagationRPrivate,
	PropagationPrivate,
	PropagationRShared,
	PropagationShared,
	PropagationRSlave,
	PropagationSlave,
}

// BindOptions defines options specific to mounts of type "bind".
type BindOptions struct {
	Propagation Propagation `json:"propagation,omitempty"`
}

// VolumeOptions represents the options for a mount of type volume.
type VolumeOptions struct {
	NoCopy       bool          `json:"noCopy,omitempty"`
	DriverConfig *DriverConfig `json:"driverConfig,omitempty"`
	SubPath      string        `json:"subPath,omitempty"`
}

// Driver represents a volume driver.
type DriverConfig struct {
	Name    string            `json:"name,omitempty"`
	Options map[string]string `json:"options,omitempty"`
}

func (d DriverConfig) String() string {
	if d.Name == "" {
		return ""
	}
	str := ""
	if len(d.Name) > 0 {
		str = "driver=" + d.Name
	}

	if len(d.Options) > 0 {
		var opts []string
		for k, v := range d.Options {
			opts = append(opts, fmt.Sprintf("%s=%s", k, v))
		}
		sort.Strings(opts)
		str += "," + strings.Join(opts, ",")
	}

	return str
}

// TmpfsOptions defines options specific to mounts of type "tmpfs".
type Tmpfs struct {
	SizeBytes int64  `json:"sizeBytes,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	Path      string `json:"path,omitempty" norman:"required"`
}

func (t Tmpfs) String() string {
	opts := ""

	if t.SizeBytes == 0 {
		opts = addOpt(opts, "size="+units.BytesSize(float64(t.SizeBytes)))
	}

	if t.ReadOnly {
		opts = addOpt(opts, "ro")
	}

	return t.Path + opts
}
