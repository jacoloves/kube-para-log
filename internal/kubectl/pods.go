package kubectl

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type PodList struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	} `json:"items"`
}

// FindMatchingPods runs `kubectl get pods -o json` and returns pod names containing the keyword
func FindMatchingPods(keyword string, namespace string) ([]string, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run kubectl: %w", err)
	}

	var podList PodList
	if err := json.Unmarshal(output, &podList); err != nil {
		return nil, fmt.Errorf("failed to parse kubectl output: %w", err)
	}

	var matched []string
	for _, item := range podList.Items {
		name := item.Metadata.Name
		if strings.Contains(name, keyword) {
			matched = append(matched, name)
		}
	}

	return matched, nil
}
