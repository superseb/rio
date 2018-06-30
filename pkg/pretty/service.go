package pretty

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/mapper"
	pm "github.com/rancher/rio/pkg/pretty/mapper"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func serviceMappers() []types.Mapper {
	return []types.Mapper{
		// Sorted by field name
		pm.SingleSlice{Field: "capAdd"},
		pm.SingleSlice{Field: "capDrop"},
		pm.SingleSlice{Field: "dns"},
		pm.SingleSlice{Field: "dnsOption"},
		pm.SingleSlice{Field: "dnsSearch"},
		pm.Shlex{Field: "command"},
		pm.NewConfigMapping("configs"),
		pm.MapToSlice{Field: "configs", Sep: ":"},
		pm.NewDeviceMapping("devices"),
		pm.MapToSlice{Field: "devices", Sep: ":"},
		pm.AliasField{Field: "environment", Names: []string{"env"}},
		pm.MapToSlice{Field: "environment", Sep: "="},
		pm.MapToSlice{Field: "extraHosts", Sep: ":"},
		&mapper.Embed{Field: "healthcheck"},
		mapper.Move{From: "memoryBytes", To: "memory"},
		pm.Bytes{"memory"},
		mapper.Move{From: "memoryReservationBytes", To: "memoryReservation"},
		pm.Bytes{"memoryReservation"},
		mapper.Move{From: "nanoCpus", To: "cpus"},
		pm.AliasField{Field: "net", Names: []string{"network"}},
		pm.AliasValue{Field: "net", Alias: map[string][]string{
			"default": {"bridge"}},
		},
		pm.NewPortBinding("ports"),
		pm.AliasValue{Field: "restart", Alias: map[string][]string{
			"never":      {"no"},
			"on-failure": {"OnFailure"}},
		},
		mapper.Drop{Field: "spaceId", IgnoreDefinition: true},
		mapper.Drop{Field: "stackId", IgnoreDefinition: true},
		pm.AliasField{Field: "stdinOpen", Names: []string{"interactive"}},
		pm.Duration{Field: "stopGracePeriod"},
		pm.NewTmpfs("tmpfs"),
		pm.SingleSlice{Field: "tmpfs"},
		pm.SingleSlice{Field: "volumesFrom"},
		pm.NewMounts("volumes"),
		pm.SingleSlice{Field: "volumes"},
	}
}

func services(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, client.ServiceRevision{}, serviceMappers()...).
		AddMapperForType(&Version, client.Service{}, serviceMappers()...)
}
