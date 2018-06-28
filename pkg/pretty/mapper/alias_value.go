package mapper

import (
	"github.com/rancher/norman/types"
)

type AliasValue struct {
	Field string
	Alias map[string][]string
}

func (d AliasValue) FromInternal(data map[string]interface{}) {
}

func (d AliasValue) ToInternal(data map[string]interface{}) error {
	v, ok := data[d.Field]
	if !ok {
		return nil
	}

	for name, values := range d.Alias {
		for _, value := range values {
			if value == v {
				data[d.Field] = name
			}
		}
	}

	return nil
}

func (AliasValue) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
