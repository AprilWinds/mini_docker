package main

import (
	"log/slog"
	"mini_docker/cmd"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mini_docker",
		Usage: "mini_docker is simple docker clone",
		Commands: []*cli.Command{
			cmd.Exec,
			cmd.Export,
			cmd.Logs,
			cmd.PS,
			cmd.RM,
			cmd.Run,
			cmd.Setup,
			cmd.Stop,
			cmd.Net,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error("run mini_docker error", "err", err)
	}
}
