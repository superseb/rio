package data

import (
	"github.com/rancher/rio/pkg/space"
	"github.com/rancher/rio/types"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/api/errors"
)

func AddData(rContext *types.Context) error {
	if err := addNameSpace(rContext); err != nil {
		return err
	}
	//if rContext.Embedded {
	//	if err := addCoreDNS(); err != nil {
	//		return err
	//	}
	//}
	//if err := addConduit(); err != nil {
	//	return err
	//}
	return nil
}

func addNameSpace(rContext *types.Context) error {
	_, err := rContext.Core.Namespaces("").Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rio-system",
			Labels: map[string]string{
				space.SpaceLabel: "true",
			},
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	def, err := rContext.Core.Namespaces("").Get("default", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	if def.Labels == nil {
		def.Labels = map[string]string{}
	}

	if def.Labels[space.SpaceLabel] != "true" {
		def.Labels[space.SpaceLabel] = "true"
		_, err := rContext.Core.Namespaces("").Update(def)
		return err
	}

	return nil
}
