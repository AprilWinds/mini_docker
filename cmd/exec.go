package cmd

import (
	"log/slog"
	"mini_docker/internal/container"

	"github.com/urfave/cli/v2"
)

var Exec = &cli.Command{
	Name:  "exec",
	Usage: "exec a command in a container",
	Action: func(c *cli.Context) error {
		if len(c.Args().Slice()) < 2 {
			slog.Error("missing container id and command argument")
			return nil
		}
		container.Exec(c.Args().Get(0), c.Args().Slice()[1:])
		return nil
	},
}
