package mapper

import (
	"fmt"

	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/sirupsen/logrus"
)

type NewObject func() fmt.Stringer
type ToObjects func([]string) ([]interface{}, error)

type ObjectsToSlice struct {
	Field     string
	NewObject NewObject
	ToObjects ToObjects
}

func (p ObjectsToSlice) FromInternal(data map[string]interface{}) {
	if data == nil {
		return
	}

	objs, ok := data[p.Field]
	if !ok {
		return
	}

	var result []string
	for _, obj := range convert.ToMapSlice(objs) {
		target := p.NewObject()
		if err := convert.ToObj(obj, target); err != nil {
			logrus.Errorf("Failed to unmarshal slice to object: %v", err)
			continue
		}

		result = append(result, target.String())
	}

	if len(result) == 0 {
		delete(data, p.Field)
	} else {
		data[p.Field] = result
	}
}

func (p ObjectsToSlice) ToInternal(data map[string]interface{}) error {
	if data == nil {
		return nil
	}

	tmpfs, ok := data[p.Field]
	if !ok {
		return nil
	}

	strSlice := convert.ToStringSlice(tmpfs)
	if len(strSlice) == 0 {
		return nil
	}

	tmpfsObjects, err := p.ToObjects(strSlice)
	if err != nil {
		return err
	}

	var result []interface{}
	for _, tmpfsObject := range tmpfsObjects {
		obj, err := convert.EncodeToMap(tmpfsObject)
		if err != nil {
			return err
		}
		result = append(result, obj)
	}

	data[p.Field] = result
	return nil
}

func (ObjectsToSlice) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	return nil
}
