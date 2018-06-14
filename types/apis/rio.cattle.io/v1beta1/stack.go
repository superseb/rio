package v1beta1

import (
	"github.com/rancher/norman/condition"
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	StackConditionNamespaceCreated = condition.Cond("NamespaceCreated")
	StackConditionParsed           = condition.Cond("Parsed")
)

type Stack struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec   `json:"spec"`
	Status StackStatus `json:"status"`
}

type StackSpec struct {
	Description string            `json:"description"`
	Templates   map[string]string `json:"templates"`
}

type StackStatus struct {
	Conditions []Condition `json:"conditions"`
}

type StackScoped struct {
	StackName string `json:"stackName,omitempty" norman:"type=reference[stack],required"`
	SpaceName string `json:"spaceName,omitempty" norman:"type=reference[/v1beta1-rio/schemas/space]"`
}

type InternalStack struct {
	Services map[string]Service
}
