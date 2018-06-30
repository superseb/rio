//    Copyright 2013-2017 Docker, Inc.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       https://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package create

import (
	"strconv"

	"github.com/rancher/rio/cli/pkg/kv"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func ParseConfigs(configs []string) ([]client.ConfigMapping, error) {
	var result []client.ConfigMapping
	for _, config := range configs {
		mapping, err := parseConfig(config)
		if err != nil {
			return nil, err
		}
		result = append(result, mapping)
	}

	return result, nil
}

func parseConfig(device string) (client.ConfigMapping, error) {
	result := client.ConfigMapping{}

	mapping, optStr := kv.Split(device, ",")
	result.Source, result.Target = kv.Split(mapping, ":")
	opts := kv.SplitMap(optStr, ",")

	if i, err := strconv.Atoi(opts["uid"]); err == nil {
		result.UID = int64(i)
	}

	if i, err := strconv.Atoi(opts["gid"]); err == nil {
		result.GID = int64(i)
	}
	result.Mode = opts["mode"]

	return result, nil
}
