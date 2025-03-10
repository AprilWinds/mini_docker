package container

import (
	"fmt"
	"mini_docker/internal/util"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"
)

func getProcess(pid int) (*os.Process, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("find process failed: %w", err)
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return nil, fmt.Errorf("process not running: %w", err)
	}
	return process, nil
}

func PS() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID\tNAME\tIMAGE\tCMD\tCREATED\tSTATUS")

	dir, err := os.ReadDir(StorageRootDir)
	if err != nil {
		util.LogAndExit("read root dir failed:", err)
	}
	for _, d := range dir {
		if d.IsDir() {

			c, err := getContainer(d.Name())
			if err != nil {
				continue
			}

			status := "stop"
			if _, err := getProcess(c.PId); err == nil {
				status = "running"
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				c.Id,
				c.Setting.Name,
				c.Setting.ImageName,
				strings.Join(c.Setting.CMD, " "),
				c.CreateTime.Format("2006-01-02 15:04:05"),
				status,
			)
		}
	}
	w.Flush()
}
