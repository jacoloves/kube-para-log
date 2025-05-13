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

func GuessBestContainerName(pod string, namesapce string) (string, error) {
	cmd := exec.Command("kubectl", "get", "pod", pod, "-n", namesapce, "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get pod json: %w", err)
	}

	var result struct {
		Status struct {
			ContainerStatuses []struct {
				Name string `json:"name"`
			} `json:"containerStatuses"`
		} `json:"status"`
	}

	if err := json.Unmarshal(out, &result); err != nil {
		return "", fmt.Errorf("failed to parse pod json: %w", err)
	}

	for _, cs := range result.Status.ContainerStatuses {
		if isSidecarContainer(cs.Name) {
			return cs.Name, nil
		}
	}

	return "", fmt.Errorf("no non-sidecar container found")
}

func isSidecarContainer(name string) bool {
	sidecars := []string{
		"envoy", "fluentd", "datadog",
	}
	for _, sc := range sidecars {
		if strings.Contains(name, sc) {
			return true
		}
	}

	return false
}
