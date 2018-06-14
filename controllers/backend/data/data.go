package data

import (
	"github.com/rancher/rio/types"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/api/errors"
)

func AddData(rContext *types.Context) error {
	if err := addNameSpace(rContext); err != nil {
		return err
	}
	if rContext.Embedded {
		if err := addCoreDNS(); err != nil {
			return err
		}
	}
	if err := addConduit(); err != nil {
		return err
	}
	return nil
}

func addNameSpace(rContext *types.Context) error {
	_, err := rContext.Core.Namespaces("").Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rio-system",
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}
