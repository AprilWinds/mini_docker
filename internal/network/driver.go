package network

import (
	"fmt"
	"mini_docker/internal/util"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type driverType string

func createBridge(name string, rawIp string) error {
	if _, err := netlink.LinkByName(name); err == nil {
		return fmt.Errorf("driver already exists: %s", name)
	}

	bridge := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
	}

	if err := netlink.LinkAdd(bridge); err != nil {
		return fmt.Errorf("failed to add link: %w", err)
	}

	if err := setLinkIP(bridge, rawIp); err != nil {
		return fmt.Errorf("failed to setup IP: %w", err)
	}

	if err := netlink.LinkSetUp(bridge); err != nil {
		return fmt.Errorf("failed to set link up: %w", err)
	}

	return nil
}

func createVeth(name string) (peerName string, err error) {
	if _, err := netlink.LinkByName(name); err == nil {
		return "", fmt.Errorf("driver already exists: %s", name)
	}

	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
		PeerName: util.ReverseStr(name),
	}

	if err := netlink.LinkAdd(veth); err != nil {
		return "", err
	}

	if err := netlink.LinkSetUp(veth); err != nil {
		return "", err
	}

	return veth.PeerName, nil
}

func setLinkIP(link netlink.Link, rawIp string) error {
	addr, err := netlink.ParseAddr(rawIp)
	if err != nil {
		return err
	}

	err = netlink.AddrAdd(link, addr)
	if err != nil {
		return err
	}

	return nil
}

func connectBridge(bridgeName string, vethName string) error {
	bridge, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	veth, err := netlink.LinkByName(vethName)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(veth, bridge); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(veth); err != nil {
		return err
	}

	return nil
}

func setPeerIP(vethName string, rawIp string) error {
	peerVeth, err := netlink.LinkByName(vethName)
	if err != nil {
		return fmt.Errorf("failed to get peer veth: %w", err)
	}
	if err := setLinkIP(peerVeth, rawIp); err != nil {
		return fmt.Errorf("failed to setup IP: %w", err)
	}
	if err = netlink.LinkSetUp(peerVeth); err != nil {
		return fmt.Errorf("failed to set peer veth up: %w", err)
	}
	return nil
}

func movePeerToNS(vethName string, ns netns.NsHandle) error {
	peerVeth, err := netlink.LinkByName(vethName)
	if err != nil {
		return fmt.Errorf("failed to get peer veth: %w", err)
	}
	if err = netlink.LinkSetNsFd(peerVeth, int(ns)); err != nil {
		return fmt.Errorf("failed to set peer veth to ns: %w", err)
	}
	return nil
}
