package cmd

import (
	"mini_docker/internal/network"
	"mini_docker/internal/util"
	"net"

	"github.com/urfave/cli/v2"
)

var Net = &cli.Command{
	Name:  "net",
	Usage: "manage network",
	Subcommands: []*cli.Command{
		{
			Name:  "create",
			Usage: "create a network",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet",
				},
			},
			Action: func(c *cli.Context) error {
				subnet := c.String("subnet")
				if subnet == "" || !isIPv4(subnet) {
					util.LogAndExit("missing subnet", nil)
				}

				if c.Args().Len() < 1 {
					util.LogAndExit("missing network name", nil)
				}
				network.Create(c.Args().Get(0), subnet)
				return nil
			},
		},

		{
			Name:  "rm",
			Usage: "remove a network",
			Action: func(c *cli.Context) error {
				if len(c.Args().Slice()) < 1 {
					util.LogAndExit("missing network name", nil)
				}
				network.RM(c.Args().Get(0))
				return nil
			},
		},

		{
			Name:  "ls",
			Usage: "list all networks",
			Action: func(c *cli.Context) error {
				network.LS()
				return nil
			},
		},
	},
}

func isIPv4(subnet string) bool {
	ip, _, err := net.ParseCIDR(subnet)
	return err == nil && ip.To4() != nil
}
