package apply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Apply(objects []runtime.Object, groupID string, generation int64) error {
	if len(objects) == 0 {
		return nil
	}

	ns, whitelist, content, err := constructApplyData(objects, groupID, generation)
	if err != nil {
		return err
	}

	return execApply(ns, whitelist, content, groupID)
}

func execApply(ns string, whitelist map[string]bool, content []byte, groupID string) error {
	output := &bytes.Buffer{}
	errOutput := &bytes.Buffer{}
	cmd := exec.Command("kubectl", "-n", ns, "apply", "--prune", "-l", "apply.cattle.io/groupID="+groupID, "-o", "json", "-f", "-")
	for group := range whitelist {
		cmd.Args = append(cmd.Args, "--prune-whitelist="+group)
	}
	cmd.Stdin = bytes.NewReader(content)
	cmd.Stdout = output
	cmd.Stderr = errOutput

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to apply: %s", errOutput.String())
	}

	fmt.Println("Applied", output.String())
	return nil
}

func constructApplyData(objects []runtime.Object, groupID string, generation int64) (string, map[string]bool, []byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	whitelist := map[string]bool{}
	ns := ""
	for _, obj := range objects {
		objType, err := meta.TypeAccessor(obj)
		if err != nil {
			return "", nil, nil, fmt.Errorf("resource type data can not be accessed")
		}

		metaObj, ok := obj.(v1.Object)
		if !ok {
			return "", nil, nil, fmt.Errorf("resource type is not a meta object")
		}
		labels := metaObj.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}
		labels["apply.cattle.io/groupID"] = groupID
		labels["apply.cattle.io/generationID"] = strconv.FormatInt(generation, 10)
		metaObj.SetLabels(labels)

		if len(ns) == 0 {
			ns = metaObj.GetNamespace()
		}

		whitelist[fmt.Sprintf("%s/%s", objType.GetAPIVersion(), objType.GetKind())] = true
		if err := encoder.Encode(obj); err != nil {
			return "", nil, nil, errors.Wrapf(err, "failed to encode %s/%s/%s/%s", objType.GetAPIVersion(), objType.GetKind(), metaObj.GetNamespace(), metaObj.GetName())
		}
	}

	return ns, whitelist, buffer.Bytes(), nil
}
