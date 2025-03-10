package cmd

import (
	"mini_docker/internal/container"

	"github.com/urfave/cli/v2"
)

var PS = &cli.Command{
	Name:  "ps",
	Usage: "list all containers",
	Action: func(c *cli.Context) error {
		container.PS()
		return nil
	},
}
