package schema

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/mapper"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/rancher/rio/types/factory"
	rm "github.com/rancher/rio/types/mapper"
	mapper2 "github.com/rancher/types/mapper"
)

var (
	Version = types.APIVersion{
		Version:          "v1beta1",
		Group:            "rio.cattle.io",
		Path:             "/v1beta1-rio/spaces",
		SubContext:       true,
		SubContextSchema: "/v1beta1-rio/schemas/space",
	}

	Schemas = factory.Schemas(&Version).
		Init(stackTypes).
		Init(serviceTypes)

	APIStackSchema = Schemas.Schema(&Version, "internalStack")
)

func serviceTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v1beta1.ServiceRevision{},
			mapper2.Status{},
			&mapper.Embed{Field: "spec"},
			&mapper.ReadOnly{
				Field:     "status",
				SubFields: true,
			},
			&mapper.Embed{Field: "status"},
		).
		AddMapperForType(&Version, v1beta1.ServiceStatus{},
			&rm.DeploymentStatus{},
			mapper.Drop{Field: "deploymentStatus"},
		).
		AddMapperForType(&Version, v1beta1.Service{},
			mapper.Drop{Field: "namespace"},
		).
		MustImport(&Version, v1beta1.Service{}).
		MustImport(&Version, v1beta1.InternalStack{})
}

func stackTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v1beta1.Stack{},
			mapper.Move{From: "namespace", To: "spaceId", CodeName: "SpaceID"}).
		MustImportAndCustomize(&Version, v1beta1.Stack{}, func(schema *types.Schema) {
			schema.MustCustomizeField("spaceId", func(f types.Field) types.Field {
				f.Type = "reference[space]"
				return f
			})
		})
}
