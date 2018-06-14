package stack

import (
	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/norman/types/mapper"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1/schema"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

var (
	Version = types.APIVersion{
		Version: "v1beta1",
		Group:   "export.cattle.io",
		Path:    "/v1beta1-export",
	}

	Schemas = NewSchemas().
		AddMapperForType(&Version, client.VolumeOptions{},
			AliasField{Field: "noCopy", Names: []string{"nocopy"}},
		).
		AddMapperForType(&Version, client.Mount{},
			AliasField{Field: "kind", Names: []string{"type"}},
		).
		AddMapperForType(&Version, client.HealthConfig{},
			Shlex{"test"},
			mapper.Move{From: "intervalSeconds", To: "interval"},
			Duration{Field: "interval"},
			AliasField{Field: "interval", Names: []string{"period", "periodSeconds"}},
			mapper.Move{From: "timeoutSeconds", To: "timeout"},
			Duration{Field: "timeout"},
			mapper.Move{From: "initialDelaySeconds", To: "initialDelay"},
			Duration{Field: "initialDelay"},
			AliasField{Field: "initialDelay", Names: []string{"startPeriod"}},
			AliasField{Field: "healthyThreshold", Names: []string{"retries", "successThreshold"}},
			AliasField{Field: "unhealthyThreshold", Names: []string{"failureThreshold"}},
		).
		AddMapperForType(&Version, client.ServiceRevision{},
			SingleSlice{Field: "capAdd"},
			SingleSlice{Field: "capDrop"},
			SingleSlice{Field: "dns"},
			SingleSlice{Field: "dnsOption"},
			SingleSlice{Field: "dnsSearch"},
			Shlex{Field: "command"},
			NewDeviceMapping("devices"),
			MapToSlice{Field: "devices", Sep: ":"},
			AliasField{Field: "environment", Names: []string{"env"}},
			MapToSlice{Field: "environment", Sep: "="},
			MapToSlice{Field: "extraHosts", Sep: ":"},
			&mapper.Embed{Field: "healthcheck"},
			mapper.Move{From: "memoryBytes", To: "memory"},
			Bytes{"memory"},
			mapper.Move{From: "memoryReservationBytes", To: "memoryReservation"},
			Bytes{"memoryReservation"},
			mapper.Move{From: "nanoCpus", To: "cpus"},
			AliasField{Field: "net", Names: []string{"network"}},
			AliasValue{Field: "net", Alias: map[string][]string{
				"default": {"bridge"}},
			},
			NewPortBinding("ports"),
			AliasValue{Field: "restart", Alias: map[string][]string{
				"never":      {"no"},
				"on-failure": {"OnFailure"}},
			},
			AliasField{Field: "stdinOpen", Names: []string{"interactive"}},
			Duration{Field: "stopGracePeriod"},
			NewTmpfs("tmpfs"),
			SingleSlice{Field: "tmpfs"},
			SingleSlice{Field: "volumesFrom"},
			NewMounts("volumes"),
			SingleSlice{Field: "volumes"},
			mapper.Drop{Field: "spaceId"},
			mapper.Drop{Field: "stackId"},
		).
		AddMapperForType(&Version, client.Service{},
			SingleSlice{Field: "capAdd"},
			SingleSlice{Field: "capDrop"},
			SingleSlice{Field: "dns"},
			SingleSlice{Field: "dnsOption"},
			SingleSlice{Field: "dnsSearch"},
			Shlex{Field: "command"},
			NewDeviceMapping("devices"),
			MapToSlice{Field: "devices", Sep: ":"},
			AliasField{Field: "environment", Names: []string{"env"}},
			MapToSlice{Field: "environment", Sep: "="},
			MapToSlice{Field: "extraHosts", Sep: ":"},
			&mapper.Embed{Field: "healthcheck"},
			mapper.Move{From: "memoryBytes", To: "memory"},
			Bytes{"memory"},
			mapper.Move{From: "memoryReservationBytes", To: "memoryReservation"},
			Bytes{"memoryReservation"},
			mapper.Move{From: "nanoCpus", To: "cpus"},
			AliasField{Field: "net", Names: []string{"network"}},
			AliasValue{Field: "net", Alias: map[string][]string{
				"default": {"bridge"}},
			},
			NewPortBinding("ports"),
			AliasValue{Field: "restart", Alias: map[string][]string{
				"never":      {"no"},
				"on-failure": {"OnFailure"}},
			},
			AliasField{Field: "stdinOpen", Names: []string{"interactive"}},
			Duration{Field: "stopGracePeriod"},
			NewTmpfs("tmpfs"),
			SingleSlice{Field: "tmpfs"},
			SingleSlice{Field: "volumesFrom"},
			NewMounts("volumes"),
			SingleSlice{Field: "volumes"},
			mapper.Drop{Field: "spaceId"},
			mapper.Drop{Field: "stackId"},
		).
		MustImport(&Version, client.Service{}).
		MustImport(&Version, Stack{})

	YAMLStackSchema = Schemas.Schema(&Version, "stack")
)

func NewSchemas() *types.Schemas {
	schemas := types.NewSchemas()
	schemas.DefaultPostMappers = func() []types.Mapper {
		return []types.Mapper{
			mapper.Drop{
				Field:            "type",
				IgnoreDefinition: true,
			},
		}
	}
	return schemas
}

type Stack struct {
	types.Resource
	Services map[string]client.Service `json:"services"`
}

type ExportFormatter struct {
}

func (s *ExportFormatter) Format(request *types.APIContext, resource *types.RawResource) {
	if request.Option("export") != "true" {
		return
	}

	stack := map[string]interface{}{}
	s.addServices(request, resource, stack)

	Schemas.Schema(&Version, client.StackType).Mapper.FromInternal(stack)

	*resource = types.RawResource{
		Schema:       resource.Schema,
		DropReadOnly: true,
		Values:       stack,
	}
}

func (s *ExportFormatter) FormatService(request *types.APIContext, resource *types.RawResource) {
	if request.Option("export") != "true" {
		return
	}

	stack := map[string]interface{}{
		"services": map[string]interface{}{
			resource.Values["name"].(string): resource.Values,
		},
	}
	s.addServices(request, resource, stack)
	delete(resource.Values, "name")
	Schemas.Schema(&Version, client.StackType).Mapper.FromInternal(stack)

	*resource = types.RawResource{
		Schema:       resource.Schema,
		DropReadOnly: true,
		Values:       stack,
	}
}

func (s *ExportFormatter) addServices(request *types.APIContext, resource *types.RawResource, data map[string]interface{}) {
	var collection []map[string]interface{}
	err := access.List(request, &schema.Version, client.ServiceType, byStackID(resource.ID), &collection)
	if err != nil || len(collection) == 0 {
		return
	}

	services := map[string]interface{}{}
	for _, data := range collection {
		name := convert.ToString(data["name"])
		delete(data, "name")
		services[name] = data
	}

	data["services"] = services
}

func byStackID(id string) *types.QueryOptions {
	return &types.QueryOptions{
		Conditions: []*types.QueryCondition{
			types.EQ(client.ServiceFieldStackId, id),
		},
	}
}
