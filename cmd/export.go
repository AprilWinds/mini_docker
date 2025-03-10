package cmd

import (
	"log/slog"
	"mini_docker/internal/container"

	"github.com/urfave/cli/v2"
)

var Export = &cli.Command{
	Name:  "export",
	Usage: "export a container to an image with name",
	Action: func(c *cli.Context) error {
		if len(c.Args().Slice()) < 2 {
			slog.Error("missing container id and imageName argument")
			return nil
		}
		container.Export(c.Args().Get(0), c.Args().Get(1))
		return nil
	},
}
