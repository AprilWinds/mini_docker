package container

import (
	"fmt"
	"mini_docker/internal/util"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func Exec(containerID string, command []string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	c, err := getContainer(containerID)
	if err != nil {
		util.LogAndExit("get container failed:", err)
	}
	if _, err := getProcess(c.PId); err != nil {
		util.LogAndExit("container not running", nil)
	}

	args := []string{"--target", strconv.Itoa(c.PId), "--all", "env"}
	args = append(args, c.Setting.Env...)
	args = append(args, command...)
	fmt.Println(args)
	cmd := exec.Command("nsenter", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		util.LogAndExit("exec failed:", err)
	}
}
