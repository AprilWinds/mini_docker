package cmd

import (
	"log/slog"
	"mini_docker/internal/container"

	"github.com/urfave/cli/v2"
)

var RM = &cli.Command{
	Name:  "rm",
	Usage: "remove a container by id",
	Action: func(c *cli.Context) error {
		if len(c.Args().Slice()) < 1 {
			slog.Error("missing container id argument")
			return nil
		}
		container.RM(c.Args().Get(0))
		return nil
	},
}
