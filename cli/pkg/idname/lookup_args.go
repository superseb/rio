package idname

import (
	"fmt"

	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func LookupArgs(args []string) (string, []string, error) {
	types := []string{client.ServiceType, client.StackType}
	name := ""
	if len(args) == 1 {
		name = args[0]
	} else if len(args) == 2 {
		types = []string{args[0]}
		name = args[1]
	} else {
		return "", nil, fmt.Errorf("wrong number of arguments (%d)", len(args))
	}

	return name, types, nil
}
