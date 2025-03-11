package network

import (
	"mini_docker/internal/util"

	"github.com/vishvananda/netns"
)

func Apply(pid int, networkName string, mappingPort []string) {

	vethName := util.GetRandomStr()
	peerVethName, err := createVeth(vethName)
	if err != nil {
		util.LogAndExit("failed to create veth", err)
	}
	if err := connectBridge(networkName, vethName); err != nil {
		util.LogAndExit("failed to connect bridge", err)
	}

	ons, err := netns.Get()
	if err != nil {
		util.LogAndExit("failed to get host netns", err)
	}
	cns, err := netns.GetFromPid(pid)
	if err != nil {
		util.LogAndExit("failed to get container netns", err)
	}

	if err := movePeerToNS(peerVethName, cns, "192.168.1.5/24"); err != nil {
		util.LogAndExit("failed to move peer veth to container netns", err)
	}

	if err := netns.Set(cns); err != nil {
		util.LogAndExit("failed to set container netns", err)
	}
	defer netns.Set(ons)
	if err := setContainerMapping(mappingPort[0], mappingPort[1]); err != nil {
		util.LogAndExit("failed to set container mapping", err)
	}
}
