package client

const (
	InternalStackType          = "internalStack"
	InternalStackFieldServices = "services"
)

type InternalStack struct {
	Services map[string]Service `json:"services,omitempty" yaml:"services,omitempty"`
}
