package stack

import (
	"github.com/docker/go-units"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
)

type Bytes struct {
	Field string
}

func (d Bytes) FromInternal(data map[string]interface{}) {
	v, ok := data[d.Field]
	if !ok {
		return
	}

	n, err := convert.ToNumber(v)
	if err != nil {
		return
	}

	data[d.Field] = units.BytesSize(float64(n))
}

func (d Bytes) ToInternal(data map[string]interface{}) error {
	v, ok := data[d.Field]
	if !ok {
		return nil
	}

	if str, ok := v.(string); ok {
		sec, err := units.RAMInBytes(str)
		if err != nil {
			return err
		}
		data[d.Field] = sec
	}

	return nil
}

func (Bytes) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
