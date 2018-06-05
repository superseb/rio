package stack

import (
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rio/api/stack"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1/schema"
	"k8s.io/apimachinery/pkg/runtime"
)

type template struct {
	name    string
	content []byte
}

func (s *stackController) parseServices(stack *v1beta1.Stack) (*v1beta1.InternalStack, error) {
	internalStack := &v1beta1.InternalStack{
		Services: map[string]v1beta1.Service{},
	}

	_, err := v1beta1.StackConditionParsed.Do(stack, func() (runtime.Object, error) {
		templates, err := getTemplates(stack)
		if err != nil {
			return nil, err
		}

		for _, template := range templates {
			parsed, err := templateToStack(template)
			if err != nil {
				return nil, err
			}

			for name, service := range parsed.Services {
				internalStack.Services[name] = service
			}
		}

		return nil, nil
	})

	return internalStack, err
}

func templateToStack(template template) (*v1beta1.InternalStack, error) {
	data, err := ParseYAML(template.content)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml for %s", template.name)
	}

	if err := stack.YAMLStackSchema.Mapper.ToInternal(data); err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml for %s", template.name)
	}

	if err := schema.APIStackSchema.Mapper.ToInternal(data); err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml for %s", template.name)
	}

	stack := &v1beta1.InternalStack{}
	if err := convert.ToObj(data, stack); err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml for %s", template.name)
	}

	return stack, nil
}

func getTemplates(stack *v1beta1.Stack) ([]template, error) {
	var templates []template
	for name, value := range stack.Spec.Templates {
		content, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template [%s]: %v", name, err)
		}
		templates = append(templates, template{
			name:    name,
			content: content,
		})
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].name < templates[j].name
	})

	return templates, nil
}
