package tmux

import (
	"fmt"
	"os"
	"os/exec"
)

var PaneColors = []string{
	"fg=colour15",
	"fg=colour196",
	"fg=colour208",
	"fg=colour226",
	"fg=colour46",
	"fg=colour51",
	"fg=colour27",
	"fg=colour129",
}

func StartTmuxWithLogs(sessionName string, podNames []string, namespace string) error {
	if len(podNames) == 0 {
		return fmt.Errorf("no pods to display logs for")
	}

	insideTmux := isInsideTmux()

	if !insideTmux {
		// kill session if it exists
		_ = exec.Command("tmux", "kill-session", "-t", sessionName).Run()

		// 1. create new tmux session (detached)
		if err := exec.Command("tmux", "new-session", "-d", "-s", sessionName).Run(); err != nil {
			return fmt.Errorf("failed to create tmux session: %w", err)
		}
	}

	// use session name or current if inside tmux
	target := sessionName
	if insideTmux {
		target = "."
	}

	// 2. create panes and run logs
	for i, pod := range podNames {
		if i > 0 {
			// split pane, select new one, run command
			splitDirection := "-h"
			if i%2 == 0 {
				splitDirection = "-v"
			}
			exec.Command("tmux", "split-window", splitDirection, "-t", target).Run()
			exec.Command("tmux", "select-layout", "-t", target, "tiled").Run()
		}

		// setting pane name to pod name
		exec.Command("tmux", "select-pane", "-t", target, "-T", pod).Run()

		// setting pane color
		color := PaneColors[i%len(PaneColors)]
		exec.Command("tmux", "select-pane", "-t", target, "-P", color).Run()

		cmd := fmt.Sprintf("kubectl logs -f %s -n %s", pod, namespace)
		if err := exec.Command("tmux", "send-keys", "-t", target, cmd, "C-m").Run(); err != nil {
			return fmt.Errorf("failed to send command to tmux: %w", err)
		}
	}

	if !insideTmux {
		return exec.Command("tmux", "attach-session", "-t", sessionName).Run()
	}

	return nil // inside tmux: no attach needed
}

// isInsideTmux checks if we're already inside a tmux session
func isInsideTmux() bool {
	return len(os.Getenv("TMUX")) > 0
}
