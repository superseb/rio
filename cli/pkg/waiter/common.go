package waiter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rancher/norman/types"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func byID(c *client.Client, id, schemaType string) (*types.Resource, error) {
	var resource types.Resource

	err := c.ByID(schemaType, id, &resource)
	return &resource, err
}

func byName(c *client.Client, name, schemaType string) (*types.Resource, error) {
	var collection types.ResourceCollection

	if err := c.List(schemaType, &types.ListOpts{
		Filters: map[string]interface{}{
			"name":         name,
			"removed_null": "1",
		},
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

func Lookup(c *client.Client, name string, typeNames ...string) (*types.Resource, error) {
	for _, schemaType := range typeNames {
		resourceByID, err := byID(c, name, schemaType)
		if err == nil {
			return resourceByID, nil
		}

		byName, err := byName(c, name, schemaType)
		if err != nil {
			return nil, err
		}

		if byName != nil {
			return byName, nil
		}
	}

	return nil, fmt.Errorf("Not found: %s", name)
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
