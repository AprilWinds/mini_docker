package network

import (
	"encoding/json"
	"mini_docker/internal/util"
	"os"
	"path/filepath"
)

func Create(name, subnet string) {
	if _, err := getNetwork(name); err == nil {
		util.LogAndExit("network already exists", nil)
	}

	os.MkdirAll(stroageRootDir, 0755)

	ipMgr, err := NewIPMgr(subnet)
	if err != nil {
		util.LogAndExit("failed to create ip manager", err)
	}
	gateway, err := ipMgr.allocate()
	if err != nil {
		util.LogAndExit("failed to allocate ip", err)
	}

	if err := createBridge(name, gateway); err != nil {
		util.LogAndExit("failed to create bridge", err)
	}
	if err := setHostSNAT(name, gateway); err != nil {
		util.LogAndExit("failed to set host snat", err)
	}

	n := Network{
		Name:    name,
		Subnet:  subnet,
		Driver:  "bridge",
		Gateway: gateway,
		IPM:     ipMgr,
	}
	os.MkdirAll(filepath.Join(stroageRootDir, name), 0755)
	f, err := os.OpenFile(filepath.Join(stroageRootDir, name, ".json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		util.LogAndExit("failed to create network metadata file", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(n)
}

func getNetwork(name string) (*Network, error) {
	f, err := os.Open(filepath.Join(stroageRootDir, name, ".json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var n Network
	if err := json.NewDecoder(f).Decode(&n); err != nil {
		return nil, err
	}
	return &n, nil
}
