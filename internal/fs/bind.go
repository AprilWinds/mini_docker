package fs

import (
	"log/slog"
	"mini_docker/internal/util"
	"os"
	"syscall"
)

func BindMount(srcDir, targetDir string) {
	os.MkdirAll(srcDir, 0777)
	os.MkdirAll(targetDir, 0777)
	err := syscall.Mount(srcDir, targetDir, "bind", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		util.LogAndExit("bind mount failed:", err)
	}
}

func UnmountBind(targetDir string) {
	err := syscall.Unmount(targetDir, 0)
	if err != nil {
		slog.Error("unmount volume failed", "err", err)
	}
	os.RemoveAll(targetDir)
}
