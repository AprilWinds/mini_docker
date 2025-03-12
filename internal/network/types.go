package network

const stroageRootDir = "/var/mini_docker/network"

type Network struct {
	Name    string `json:"name"`
	Subnet  string `json:"subnet"`
	Driver  string `json:"driver"`
	Gateway string `json:"gateway"`
	IPM     *IPMgr `json:"ip_mgr"`
}

type IPMgr struct {
	StartIP uint32              `json:"start_ip"`
	EndIP   uint32              `json:"end_ip"`
	Used    map[uint32]struct{} `json:"used"`
}
