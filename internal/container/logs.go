package container

import (
	"fmt"
	"log/slog"
	"os"
)

func Logs(containerID string) {
	logs, err := os.ReadFile(getLogFile(containerID))
	if err != nil {
		slog.Error("read logs failed:", "err", err)
	}
	fmt.Println(string(logs))
}
