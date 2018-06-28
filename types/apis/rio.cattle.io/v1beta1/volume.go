package v1beta1

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Volume struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeSpec   `json:"spec,omitempty"`
	Status VolumeStatus `json:"status,omitempty"`
}

type VolumeSpec struct {
	Description string `json:"description,omitempty"`
	Driver      string `json:"driver,omitempty"`
	Template    string `json:"template,omitempty,noupdate"`
	SizeInGB    int    `json:"sizeInGb,omitempty,required"`
	StackScoped
}

type VolumeStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}
