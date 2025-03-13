package network

import (
	"mini_docker/internal/util"

	"github.com/vishvananda/netns"
)

func Apply(pid int, networkName string, mappingPort []string) {

	n, err := getNetwork(networkName)
	if err != nil {
		util.LogAndExit("failed to get network", err)
	}

	vethName := util.GetRandomStr()
	peerVethName, err := createVeth(vethName)
	if err != nil {
		util.LogAndExit("failed to create veth", err)
	}
	if err := connectBridge(networkName, vethName); err != nil {
		util.LogAndExit("failed to connect bridge", err)
	}
	cns, err := netns.GetFromPid(pid)
	if err != nil {
		util.LogAndExit("failed to get container netns", err)
	}
	if err := movePeerToNS(peerVethName, cns); err != nil {
		util.LogAndExit("failed to move peer veth to container netns", err)
	}

	ons, err := netns.Get()
	if err != nil {
		util.LogAndExit("failed to get host netns", err)
	}
	defer netns.Set(ons)
	if err := netns.Set(cns); err != nil {
		util.LogAndExit("failed to set container netns", err)
	}

	ip, err := n.IPM.allocate()
	if err != nil {
		util.LogAndExit("failed to allocate ip", err)
	}
	if err := setPeerIP(peerVethName, ip); err != nil {
		util.LogAndExit("failed to set peer ip", err)
	}
	if err := setRoute(peerVethName, n.Gateway); err != nil {
		util.LogAndExit("failed to set route", err)
	}
}
