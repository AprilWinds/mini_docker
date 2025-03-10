package cmd

import (
	"log/slog"
	"mini_docker/internal/container"

	"github.com/urfave/cli/v2"
)

var Logs = &cli.Command{
	Name:  "logs",
	Usage: "logs a container",
	Action: func(c *cli.Context) error {
		if len(c.Args().Slice()) < 1 {
			slog.Error("missing container id argument")
			return nil
		}
		container.Logs(c.Args().Get(0))
		return nil
	},
}
