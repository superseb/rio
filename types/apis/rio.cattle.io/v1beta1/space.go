package v1beta1

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Space struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpaceSpec   `json:"spec"`
	Status StackStatus `json:"status"`
}

type SpaceSpec struct {
}

type SpaceStatus struct {
	Conditions []Condition `json:"conditions"`
}
