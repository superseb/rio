package main

import (
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1/schema"
	spaceSchema "github.com/rancher/rio/types/apis/space.cattle.io/v1beta1/schema"
	"github.com/rancher/rio/types/codegen/generator"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
)

func main() {
	generator.Generate(schema.Schemas)
	generator.Generate(spaceSchema.Schemas)
	generator.GenerateNativeTypes(v1.SchemeGroupVersion, []interface{}{
		v1.Endpoints{},
		v1.Pod{},
		v1.Service{},
		v1.Secret{},
		v1.ConfigMap{},
		v1.ServiceAccount{},
		v1.ReplicationController{},
	}, []interface{}{
		v1.Node{},
		v1.ComponentStatus{},
		v1.Namespace{},
		v1.Event{},
	})
	generator.GenerateNativeTypes(v1beta2.SchemeGroupVersion, []interface{}{
		v1beta2.Deployment{},
		v1beta2.DaemonSet{},
		v1beta2.StatefulSet{},
		v1beta2.ReplicaSet{},
	}, nil)
}
