package main

import (
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1/schema"
	spaceSchema "github.com/rancher/rio/types/apis/space.cattle.io/v1beta1/schema"
	"github.com/rancher/rio/types/codegen/generator"
)

func main() {
	generator.Generate(schema.Schemas)
	generator.Generate(spaceSchema.Schemas)
}
