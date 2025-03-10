package network

const stroageRootDir = "/var/mini_docker/network"

type Network struct {
	Name    string `json:"name"`
	Subnet  string `json:"subnet"`
	Driver  string `json:"driver"`
	Gateway string `json:"gateway"`
}

type IPMgr struct {
}
