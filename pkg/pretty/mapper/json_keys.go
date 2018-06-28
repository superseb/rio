package mapper

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
)

type JSONKeys struct {
	Field string
	Names []string
}

func (d JSONKeys) FromInternal(data map[string]interface{}) {
}

func (d JSONKeys) ToInternal(data map[string]interface{}) error {
	for key, value := range data {
		newKey := convert.ToJSONKey(key)
		if newKey != key {
			data[newKey] = value
			delete(data, key)
		}
	}
	return nil
}

func (JSONKeys) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
