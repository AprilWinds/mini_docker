package network

import (
	"fmt"
	"net"

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
		return err
	}

	if err := setLinkIP(bridge, rawIp); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(bridge); err != nil {
		return err
	}

	return nil
}

func createVeth(name string) error {
	if _, err := netlink.LinkByName(name); err == nil {
		return fmt.Errorf("driver already exists: %s", name)
	}

	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
		PeerName: "P" + name,
	}

	if err := netlink.LinkAdd(veth); err != nil {
		return err
	}

	return netlink.LinkSetUp(veth)
}

func setLinkIP(link netlink.Link, rawIp string) error {
	addr, err := netlink.ParseAddr(rawIp)
	if err != nil {
		return err
	}

	return netlink.AddrAdd(link, addr)
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
		return err
	}

	if err := setLinkIP(peerVeth, rawIp); err != nil {
		return err
	}

	return netlink.LinkSetUp(peerVeth)
}

func movePeerToNS(vethName string, ns netns.NsHandle) error {
	peerVeth, err := netlink.LinkByName(vethName)

	if err != nil {
		return err
	}

	return netlink.LinkSetNsFd(peerVeth, int(ns))
}

func setRoute(vethName string, rawIp string) error {
	peerVeth, err := netlink.LinkByName(vethName)
	if err != nil {
		return err
	}
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")
	addr, err := netlink.ParseAddr(rawIp)
	defaultRoute := &netlink.Route{
		LinkIndex: peerVeth.Attrs().Index,
		Gw:        addr.IP,
		Dst:       cidr,
	}
	return netlink.RouteAdd(defaultRoute)
}
