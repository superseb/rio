package stack

import (
	"strings"

	"github.com/rancher/norman/types/convert"
	"gopkg.in/yaml.v2"
)

func ParseYAML(content []byte) (map[string]interface{}, error) {
	data := map[interface{}]interface{}{}
	err := yaml.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return convertMap(data), nil
}

func convertSlice(data []interface{}) []interface{} {
	var result []interface{}
	for _, obj := range data {
		result = append(result, convertValue(obj))
	}
	return result
}

func convertMap(data map[interface{}]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range data {
		result[convertKey(k)] = convertValue(v)
	}
	return result
}

func convertValue(val interface{}) interface{} {
	switch v := val.(type) {
	case map[interface{}]interface{}:
		return convertMap(v)
	case []interface{}:
		return convertSlice(v)
	default:
		return val
	}
}

func convertKey(k interface{}) string {
	str := convert.ToString(k)
	parts := strings.Split(str, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	return strings.Join(parts, "")
}
