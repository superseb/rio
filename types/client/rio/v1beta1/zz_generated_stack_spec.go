package client

const (
	StackSpecType                 = "stackSpec"
	StackSpecFieldAdditionalFiles = "additionalFiles"
	StackSpecFieldAnswers         = "answers"
	StackSpecFieldDescription     = "description"
	StackSpecFieldQuestions       = "questions"
	StackSpecFieldTemplate        = "template"
)

type StackSpec struct {
	AdditionalFiles map[string]string `json:"additionalFiles,omitempty" yaml:"additionalFiles,omitempty"`
	Answers         map[string]string `json:"answers,omitempty" yaml:"answers,omitempty"`
	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	Questions       []Question        `json:"questions,omitempty" yaml:"questions,omitempty"`
	Template        string            `json:"template,omitempty" yaml:"template,omitempty"`
}
