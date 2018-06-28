package mapper

import (
	"time"

	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rio/cli/cmd/create"
)

type Duration struct {
	Field string
}

func (d Duration) FromInternal(data map[string]interface{}) {
	v, ok := data[d.Field]
	if !ok {
		return
	}

	n, err := convert.ToNumber(v)
	if err != nil {
		return
	}

	data[d.Field] = (time.Duration(n) * time.Second).String()
}

func (d Duration) ToInternal(data map[string]interface{}) error {
	v, ok := data[d.Field]
	if !ok {
		return nil
	}

	if str, ok := v.(string); ok {
		sec, err := create.ParseSeconds(str, d.Field)
		if err != nil {
			return err
		}
		data[d.Field] = sec
	}

	return nil
}

func (Duration) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
