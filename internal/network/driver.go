package network

import (
	"errors"
	"fmt"

	"github.com/vishvananda/netlink"
)

func createBridge(name string, gateway string) error {
	_, err := netlink.LinkByName(name)
	if err == nil {
		return errors.New("bridge already exists: " + name)
	}

	la := netlink.NewLinkAttrs()
	la.Name = name
	br := &netlink.Bridge{LinkAttrs: la}

	// 创建桥接设备
	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("failed to create bridge: %w", err)
	}

	// 解析地址
	addr, err := netlink.ParseAddr(gateway)
	if err != nil {
		return fmt.Errorf("failed to parse subnet: %w", err)
	}

	// 添加地址
	if err := netlink.AddrAdd(br, addr); err != nil {
		return fmt.Errorf("failed to add subnet to bridge: %w", err)
	}

	// 设置桥接设备为启动状态
	if err := netlink.LinkSetUp(br); err != nil {
		return fmt.Errorf("failed to set bridge up: %w", err)
	}

	return nil
}

func deleteDriver(name string) error {
	br, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("failed to get bridge: %w", err)
	}
	if err := netlink.LinkDel(br); err != nil {
		return fmt.Errorf("failed to delete bridge: %w", err)
	}
	return nil
}
