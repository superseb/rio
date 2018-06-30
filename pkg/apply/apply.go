package apply

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/runconduit/conduit/cli/cmd"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func Content(content []byte) error {
	errOutput := &bytes.Buffer{}
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(content)
	cmd.Stdout = nil
	cmd.Stderr = errOutput

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to apply: %s", errOutput.String())
	}

	return nil
}

func Apply(objects []runtime.Object, groupID string, generation int64) error {
	if len(objects) == 0 {
		return nil
	}

	ns, whitelist, content, err := constructApplyData(objects, groupID, generation)
	if err != nil {
		return err
	}

	content, err = injectConduit(content)
	if err != nil {
		return err
	}

	return execApply(ns, whitelist, content, groupID)
}

func injectConduit(content []byte) ([]byte, error) {
	inBuf := bytes.NewBuffer(content)
	outBuf := &bytes.Buffer{}
	err := cmd.InjectYAML(inBuf, outBuf, "v0.4.1")
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}

func execApply(ns string, whitelist map[string]bool, content []byte, groupID string) error {
	output := &bytes.Buffer{}
	errOutput := &bytes.Buffer{}
	cmd := exec.Command("kubectl", "-n", ns, "apply", "--force", "--grace-period", "120", "--prune", "-l", "apply.cattle.io/groupID="+groupID, "-o", "json", "-f", "-")
	for group := range whitelist {
		cmd.Args = append(cmd.Args, "--prune-whitelist="+group)
	}
	cmd.Stdin = bytes.NewReader(content)
	cmd.Stdout = output
	cmd.Stderr = errOutput

	if err := cmd.Run(); err != nil {
		logrus.Errorf("Failed to apply %s: %s, input: %s", errOutput.String(), string(content))

		return fmt.Errorf("failed to apply: %s", errOutput.String())
	}

	if logrus.GetLevel() <= logrus.DebugLevel {
		fmt.Printf("Applied: %s", output.String())
	}

	return nil
}

func constructApplyData(objects []runtime.Object, groupID string, generation int64) (string, map[string]bool, []byte, error) {
	buffer := &bytes.Buffer{}
	whitelist := map[string]bool{}
	ns := ""
	for i, obj := range objects {
		if i > 0 {
			buffer.WriteString("\n---\n")
		}

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

		gvk := fmt.Sprintf("%s/%s", objType.GetAPIVersion(), objType.GetKind())
		if len(strings.Split(gvk, "/")) < 3 {
			gvk = "/" + gvk
		}
		whitelist[gvk] = true

		bytes, err := yaml.Marshal(obj)
		if err != nil {
			return "", nil, nil, errors.Wrapf(err, "failed to encode %s/%s/%s/%s", objType.GetAPIVersion(), objType.GetKind(), metaObj.GetNamespace(), metaObj.GetName())
		}
		buffer.Write(bytes)
	}

	return ns, whitelist, buffer.Bytes(), nil
}
