package mapper

import (
	"fmt"

	"github.com/rancher/rio/cli/cmd/create"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
)

func NewConfigMapping(field string) ObjectsToSlice {
	return ObjectsToSlice{
		Field: field,
		NewObject: func() fmt.Stringer {
			return &v1beta1.ConfigMapping{}
		},
		ToObjects: func(str []string) ([]interface{}, error) {
			var result []interface{}
			objs, err := create.ParseConfigs(str)
			if err != nil {
				return nil, err
			}
			for _, obj := range objs {
				result = append(result, obj)
			}
			return result, nil
		},
	}
}