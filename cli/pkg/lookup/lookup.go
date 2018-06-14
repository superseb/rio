package lookup

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func Lookup(c clientbase.APIBaseClientInterface, name string, typeNames ...string) (*types.Resource, error) {
	for _, schemaType := range typeNames {
		if strings.Contains(name, ":") && !strings.Contains(name, "/") {
			resourceByID, err := byID(c, name, schemaType)
			if err == nil {
				return resourceByID, nil
			}
		}

		byName, err := byName(c, name, schemaType)
		if err != nil {
			return nil, err
		}

		if byName != nil {
			return byName, nil
		}
	}

	return nil, fmt.Errorf("not found: %s", name)
}

func byID(c clientbase.APIBaseClientInterface, id, schemaType string) (*types.Resource, error) {
	var resource types.Resource

	err := c.ByID(schemaType, id, &resource)
	return &resource, err
}

func setupFilters(c clientbase.APIBaseClientInterface, name, schemaType string) (map[string]interface{}, error) {
	filters := map[string]interface{}{
		"name":         name,
		"removed_null": "1",
	}

	if schemaType != client.ServiceType {
		return filters, nil
	}

	parsedService := ParseServiceName(name)
	stackName, serviceName := parsedService.StackName, parsedService.ServiceName

	stack, err := byName(c, stackName, client.StackType)
	if err != nil {
		return nil, err
	}

	filters["stackId"] = stack.ID
	filters["name"] = serviceName

	return filters, nil
}

func byName(c clientbase.APIBaseClientInterface, name, schemaType string) (*types.Resource, error) {
	var collection types.ResourceCollection

	if schemaType == client.StackType && strings.Contains(name, "/") {
		// stacks can't be foo/bar
		return nil, nil
	}

	filters, err := setupFilters(c, name, schemaType)
	if err != nil {
		return nil, err
	}

	if err := c.List(schemaType, &types.ListOpts{
		Filters: filters,
	}, &collection); err != nil {
		return nil, err
	}

	if len(collection.Data) > 1 {
		var ids []string
		for _, data := range collection.Data {
			switch schemaType {
			default:
				ids = append(ids, fmt.Sprintf("%s (%s)", data.ID, name))
			}

		}
		index := selectFromList("Resources: ", ids)
		return &collection.Data[index], nil
	}

	if len(collection.Data) == 0 {
		return nil, nil
	}

	return &collection.Data[0], nil
}

func selectFromList(header string, choices []string) int {
	if header != "" {
		fmt.Println(header)
	}

	reader := bufio.NewReader(os.Stdin)
	selected := -1
	for selected <= 0 || selected > len(choices) {
		for i, choice := range choices {
			fmt.Printf("[%d] %s\n", i+1, choice)
		}
		fmt.Print("Select: ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		num, err := strconv.Atoi(text)
		if err == nil {
			selected = num
		}
	}
	return selected - 1
}
