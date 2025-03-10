package container

import (
	"mini_docker/internal/fs"
	"mini_docker/internal/util"
	"os"
)

func RM(containerID string) {
	c, err := getContainer(containerID)
	if err != nil {
		util.LogAndExit("get container failed:", err)
	}
	if c.Setting.Volume != nil && len(c.Setting.Volume) > 1 {
		fs.UnmountBind(getMountDir(containerID, c.Setting.Volume[1]))
	}
	fs.UnmountOverlay(getRootDir(containerID))
	os.RemoveAll(getRootDir(containerID))
}
