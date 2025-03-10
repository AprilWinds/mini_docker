package fs

import (
	"fmt"
	"mini_docker/internal/util"
	"os"
	"path/filepath"
	"syscall"
)

func MountOverlay(rootPath string) {
	os.MkdirAll(rootPath, 0755)
	lowerDir := filepath.Join(rootPath, "lower")
	upperDir := filepath.Join(rootPath, "upper")
	workDir := filepath.Join(rootPath, "work")
	mergedDir := filepath.Join(rootPath, "merged")

	os.MkdirAll(rootPath, 0755)
	os.MkdirAll(lowerDir, 0755)
	os.MkdirAll(upperDir, 0755)
	os.MkdirAll(workDir, 0755)
	os.MkdirAll(mergedDir, 0755)

	err := syscall.Mount("overlay", mergedDir, "overlay", 0, fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir))
	if err != nil {
		util.LogAndExit("mount overlay fs failed:", err)
	}
}

func UnmountOverlay(rootPath string) {
	mergedDir := filepath.Join(rootPath, "merged")
	syscall.Unmount(mergedDir, 0)
}
