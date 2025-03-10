package network

import (
	"encoding/json"
	"mini_docker/internal/util"
	"os"
	"path/filepath"
)

func Create(name, subnet string) {
	gateway := util.GetCIDRGateway(subnet)
	n := Network{
		Name:    name,
		Subnet:  subnet,
		Driver:  "bridge",
		Gateway: gateway,
	}

	if err := createBridge(name, gateway); err != nil {
		util.LogAndExit("failed to create bridge", err)
	}
	os.MkdirAll(stroageRootDir, 0755)
	f, err := os.OpenFile(filepath.Join(stroageRootDir, name+".json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		util.LogAndExit("failed to create network metadata file", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(n)
}
