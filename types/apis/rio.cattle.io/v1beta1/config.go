package v1beta1

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConfigSpec `json:"spec,omitempty"`
}

type ConfigSpec struct {
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	Encoded     bool   `json:"encoded,omitempty"`
	StackScoped
}
