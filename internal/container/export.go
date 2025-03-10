package container

import (
	"mini_docker/internal/util"
	"os/exec"
	"path/filepath"
)

func Export(containerID string, imageName string) {
	rootPath := filepath.Join(StorageRootDir, containerID)
	imageTar := "./" + imageName + ".tar.gz"
	_, err := exec.Command("tar", "-czf", imageTar, "-C", rootPath, ".").CombinedOutput()
	if err != nil {
		util.LogAndExit("tar -czf image failed:", err)
	}
}
