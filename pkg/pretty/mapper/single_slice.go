package mapper

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
)

type SingleSlice struct {
	Field string
}

func (d SingleSlice) FromInternal(data map[string]interface{}) {
	v, ok := data[d.Field]
	if !ok {
		return
	}

	ss := convert.ToStringSlice(v)
	if len(ss) == 1 {
		data[d.Field] = ss[1]
	}
}

func (d SingleSlice) ToInternal(data map[string]interface{}) error {
	v, ok := data[d.Field]
	if !ok {
		return nil
	}

	if str, ok := v.(string); ok {
		data[d.Field] = []string{str}
	}

	return nil
}

func (SingleSlice) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
