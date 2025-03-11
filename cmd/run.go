package cmd

import (
	"errors"
	"log/slog"
	"mini_docker/internal/cgroup"
	"mini_docker/internal/container"
	"mini_docker/internal/util"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var Run = &cli.Command{
	Name:  "run",
	Usage: "run a container with image name, eg: mini_docker run [options] image_name [command] [args...]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		&cli.BoolFlag{
			Name:  "it",
			Usage: "enable interactive mode, eg: -it",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "volume, eg: -v /host/path:/container/path",
		},
		&cli.Float64Flag{
			Name:  "cpus",
			Usage: "cpu limit, eg: -cpus 0.5",
		},
		&cli.StringFlag{
			Name:  "mem",
			Usage: "mem limit, eg: -mem 100m",
		},
		&cli.StringFlag{
			Name:  "e",
			Usage: "environment variables, eg: -e key=value",
		},
		&cli.StringFlag{
			Name:  "p",
			Usage: "port mapping, eg: -p 8080:80",
		},
	},
	Action: func(c *cli.Context) error {
		if len(c.Args().Slice()) < 2 {
			slog.Error("missing image name argument and command")
			return nil
		}

		setting := &container.Setting{
			ImageName: c.Args().Get(0),
			Name:      c.String("name"),
			It:        c.Bool("it"),
			CMD:       c.Args().Slice()[1:],
			Volume:    parseVolume(c),
			Env:       parseEnv(c),
			Port:      parsePort(c),
			CgroupCfg: parseCgroup(c),
		}
		slog.Info("run container", "setting", setting)
		container.Run(setting)
		return nil
	},
}

func parseCgroup(c *cli.Context) *cgroup.Config {
	cfg := &cgroup.Config{}
	cpuFloat := c.Float64("cpus")
	if cpuFloat < 0 {
		util.LogAndExit("invalid cpu limit", errors.New("invalid cpu limit"))
	}
	cfg.CPUS = float32(cpuFloat)

	memStr := c.String("mem")
	if memStr != "" {
		if !strings.HasSuffix(memStr, "m") {
			util.LogAndExit("invalid memory limit", errors.New("invalid memory limit"))
		}
		memNum, err := strconv.ParseInt(memStr[:len(memStr)-1], 10, 64)
		if err != nil && memNum < 0 {
			util.LogAndExit("invalid memory limit", err)
		}
		cfg.MemoryLimit = uint64(memNum) * 1024 * 1024
	}
	return cfg
}

func parseVolume(c *cli.Context) []string {
	volumes := []string{}
	volumeStr := c.String("v")
	if volumeStr != "" {
		v := strings.Split(volumeStr, ":")
		if len(v) != 2 {
			util.LogAndExit("invalid volume", errors.New("invalid volume"))
		}
		volumes = append(volumes, v...)
	}
	return volumes
}

func parseEnv(c *cli.Context) []string {
	envStr := c.String("e")
	if envStr != "" {
		envs := strings.Split(envStr, " ")
		for _, env := range envs {
			if !strings.Contains(env, "=") {
				util.LogAndExit("invalid env", nil)
			}
		}
		return envs
	}
	return nil
}

func parsePort(c *cli.Context) []string {
	portStr := c.String("p")
	if portStr != "" && strings.Contains(portStr, ":") {
		return strings.Split(portStr, ":")
	}
	return []string{}
}
