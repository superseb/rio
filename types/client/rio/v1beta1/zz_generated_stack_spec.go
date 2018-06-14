package client

const (
	StackSpecType             = "stackSpec"
	StackSpecFieldDescription = "description"
	StackSpecFieldTemplates   = "templates"
)

type StackSpec struct {
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Templates   map[string]string `json:"templates,omitempty" yaml:"templates,omitempty"`
}
