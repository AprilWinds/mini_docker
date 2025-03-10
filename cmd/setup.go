package cmd

import (
	"mini_docker/internal/container"

	"github.com/urfave/cli/v2"
)

var Setup = &cli.Command{
	Name:  "setup",
	Usage: "setup a container",
	Action: func(c *cli.Context) error {
		container.Setup()
		return nil
	},
}
