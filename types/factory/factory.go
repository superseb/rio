package factory

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/factory"
	m "github.com/rancher/norman/types/mapper"
	rm "github.com/rancher/rio/types/mapper"
	"github.com/rancher/types/mapper"
)

func Schemas(version *types.APIVersion) *types.Schemas {
	schemas := factory.Schemas(version)
	baseFunc := schemas.DefaultMappers
	schemas.DefaultMappers = func() []types.Mapper {
		mappers := append([]types.Mapper{
			&rm.StackScoped{},
			&mapper.Status{},
			&m.Embed{
				Field:    "status",
				Optional: true,
			},
			m.Drop{
				Field:            "conditions",
				IgnoreDefinition: true,
			},
			m.Drop{
				Field:            "ownerReferences",
				IgnoreDefinition: true,
			},
		}, baseFunc()...)
		return mappers
	}

	basePostFunc := schemas.DefaultPostMappers
	schemas.DefaultPostMappers = func() []types.Mapper {
		mappers := append(basePostFunc(),
			m.Drop{
				Field:            "annotations",
				IgnoreDefinition: true,
			})
		return mappers
	}
	return schemas
}
