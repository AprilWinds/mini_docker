package fs

import (
	"mini_docker/internal/util"
	"os"
	"path/filepath"
	"syscall"
)

func SetupMountInfo() {
	mountPrivate()
	mountRoot()
	mountProc()
}

func mountPrivate() {
	// 先将根目录重新挂载为私有，防止挂载事件传播到宿主机
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		util.LogAndExit("remount root private failed:", err)
	}
}

func mountProc() {
	// MS_NOEXEC 在本文件系统允许运行程序
	// MS_NOSUID 在本系统中运行程序的时候，允许set-user-ID set-group-ID
	// MS_NOD 这个参数是自 Linux 2.4，所有 mount 的系统都会默认设定的参数
	flags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(flags), "")
	if err != nil {
		util.LogAndExit("mount proc failed:", err)
	}
}

func mountRoot() {
	wd, err := os.Getwd()
	if err != nil {
		util.LogAndExit("get wd failed:", err)
	}

	// 为当前目录创建一个挂载点
	BindMount(wd, wd)

	// new_root：新根文件系统的路径（必须是一个挂载点）。
	// old_old：旧根文件系统挂载的目标路径（必须位于 new_root 下）。
	pivotDir := filepath.Join(wd, ".old_root")
	err = os.MkdirAll(pivotDir, 0777)
	if err != nil {
		util.LogAndExit("create pivot dir failed:", err)
	}

	err = syscall.PivotRoot(wd, pivotDir)
	if err != nil {
		util.LogAndExit("pivot root failed:", err)
	}

	err = syscall.Chdir("/")
	if err != nil {
		util.LogAndExit("chdir failed:", err)
	}

	// 更新根文件系统后，卸载旧的挂载点
	pivotDir = filepath.Join("/", ".old_root")
	err = syscall.Unmount(pivotDir, syscall.MNT_DETACH)
	if err != nil {
		util.LogAndExit("unmount pivot dir failed:", err)
	}

	err = os.RemoveAll(pivotDir)
	if err != nil {
		util.LogAndExit("remove pivot dir failed:", err)
	}
}
