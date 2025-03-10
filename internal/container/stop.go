package container

import (
	"mini_docker/internal/util"
)

func Stop(containerID string) {
	c, err := getContainer(containerID)
	if err != nil {
		util.LogAndExit("get container failed:", err)
	}
	process, err := getProcess(c.PId)
	if err != nil {
		util.LogAndExit("container not running", nil)
	}
	err = process.Kill()
	if err != nil {
		util.LogAndExit("kill container failed:", err)
	}
}
