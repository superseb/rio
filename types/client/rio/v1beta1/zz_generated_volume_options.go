package client

const (
	VolumeOptionsType              = "volumeOptions"
	VolumeOptionsFieldDriverConfig = "driverConfig"
	VolumeOptionsFieldNoCopy       = "noCopy"
	VolumeOptionsFieldSubPath      = "subPath"
)

type VolumeOptions struct {
	DriverConfig *DriverConfig `json:"driverConfig,omitempty" yaml:"driverConfig,omitempty"`
	NoCopy       bool          `json:"noCopy,omitempty" yaml:"noCopy,omitempty"`
	SubPath      string        `json:"subPath,omitempty" yaml:"subPath,omitempty"`
}
