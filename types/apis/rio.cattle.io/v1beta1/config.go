package v1beta1

import (
	"fmt"

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

type ConfigMapping struct {
	Source string `json:"source,omitempty" norman:"required"`
	Target string `json:"target,omitempty"`
	UID    int    `json:"uid,omitempty"`
	GID    int    `json:"gid,omitempty"`
	Mode   string `json:"mode,omitempty"`
}

func (c ConfigMapping) String() string {
	if c.Target == "/"+c.Source {
		c.Target = ""
	}

	msg := c.Source
	if c.Target != "" {
		msg += ":" + c.Target
	}

	if c.UID > 0 {
		msg = fmt.Sprintf("%s,uid=%d", msg, c.UID)
	}

	if c.GID > 0 {
		msg = fmt.Sprintf("%s,gid=%d", msg, c.GID)
	}

	if c.Mode != "" {
		msg = fmt.Sprintf("%s,mode=%s", msg, c.Mode)
	}

	return msg
}
