package network

import (
	"encoding/json"
	"fmt"
	"mini_docker/internal/util"
	"os"
	"path/filepath"
)

func Create(name, subnet string) {
	if _, err := getNetwork(name); err == nil {
		util.LogAndExit("network already exists", nil)
	}

	os.MkdirAll(stroageRootDir, 0755)
	if _, err := craeteNetwork(name, subnet); err != nil {
		util.LogAndExit("failed to create network", err)
	}
}

func craeteNetwork(name, subnet string) (*Network, error) {
	ipMgr, err := NewIPMgr(subnet)
	if err != nil {
		return nil, fmt.Errorf("failed to create ip manager: %w", err)
	}
	gateway, err := ipMgr.allocate()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate ip: %w", err)
	}

	if err := createBridge(name, gateway); err != nil {
		return nil, fmt.Errorf("failed to create bridge: %w", err)
	}
	if err := setHostSNAT(name, gateway); err != nil {
		return nil, fmt.Errorf("failed to set host snat: %w", err)
	}

	n := Network{
		Name:    name,
		Subnet:  subnet,
		Driver:  "bridge",
		Gateway: gateway,
		IPM:     ipMgr,
	}
	f, err := os.OpenFile(filepath.Join(stroageRootDir, name+".json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		deleteHostSNAT(name, gateway)
		return nil, fmt.Errorf("failed to create network metadata file: %w", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(n)
	return &n, nil
}

func getNetwork(networkName string) (*Network, error) {
	f, err := os.Open(filepath.Join(stroageRootDir, networkName+".json"))
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
