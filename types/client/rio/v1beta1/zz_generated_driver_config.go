package client

const (
	DriverConfigType         = "driverConfig"
	DriverConfigFieldName    = "name"
	DriverConfigFieldOptions = "options"
)

type DriverConfig struct {
	Name    string            `json:"name,omitempty" yaml:"name,omitempty"`
	Options map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
}
