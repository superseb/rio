package client

const (
	VolumeStatusType            = "volumeStatus"
	VolumeStatusFieldConditions = "conditions"
)

type VolumeStatus struct {
	Conditions []Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
}
