package container

import (
	"io"
	"log/slog"
	"mini_docker/internal/fs"
	"mini_docker/internal/util"
	"os"
	"strings"
	"syscall"
)

func readCMDFromParent() []string {
	reader := os.NewFile(3, "pipe")
	command, err := io.ReadAll(reader)
	if err != nil {
		util.LogAndExit("read command from parent failed:", err)
	}
	return strings.Split(string(command), " ")
}

func Setup() {
	fs.SetupMountInfo()
	argv := readCMDFromParent()
	slog.Info("exec container command", "command", argv)
	err := syscall.Exec(argv[0], argv, os.Environ())
	if err != nil {
		util.LogAndExit("exec container command failed:", err)
	}
}
