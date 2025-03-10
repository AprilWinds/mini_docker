package container

import (
	"mini_docker/internal/cgroup"
	"os/exec"
	"time"
)

const StorageRootDir = "/var/mini_docker/container/"

type container struct {
	Id         string    `json:"id"`
	PId        int       `json:"pid"`
	CreateTime time.Time `json:"create_time"`
	Setting    *Setting  `json:"setting"`
	child      *exec.Cmd
}

type Setting struct {
	ImageName string         `json:"image_name"`
	Name      string         `json:"name"`
	It        bool           `json:"it"`
	CMD       []string       `json:"cmd"`
	Volume    []string       `json:"volume"`
	CgroupCfg *cgroup.Config `json:"cgroup_cfg"`
	Env       []string       `json:"env"`
}
