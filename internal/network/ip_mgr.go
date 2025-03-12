package network

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

func NewIPMgr(cidr string) (*IPMgr, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	start := ipToUint32(ipNet.IP)
	mask := binary.BigEndian.Uint32(ipNet.Mask)
	end := start | ^mask

	start++
	end--

	allocator := &IPMgr{
		StartIP: start,
		EndIP:   end,
		Used:    make(map[uint32]struct{}),
	}
	return allocator, nil
}

func ipToUint32(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func uint32ToIP(ip uint32) net.IP {
	ipv4 := make(net.IP, 4)
	binary.BigEndian.PutUint32(ipv4, ip)
	return ipv4
}

func calculateMask(availableIPs uint32) net.IPMask {
	hostBits := math.Ceil(math.Log2(float64(availableIPs)))
	netBits := 32 - int(hostBits)
	return net.CIDRMask(netBits, 32)
}

func (i *IPMgr) allocate() (string, error) {
	availableIPs := 0
	for ip := i.StartIP; ip <= i.EndIP; ip++ {
		if _, exists := i.Used[ip]; !exists {
			availableIPs++
		}
	}
	if availableIPs == 0 {
		return "", fmt.Errorf("no available IP addresses in the range")
	}

	mask := calculateMask(uint32(availableIPs))
	ones, _ := mask.Size()

	for ip := i.StartIP; ip <= i.EndIP; ip++ {
		if _, exists := i.Used[ip]; !exists {
			i.Used[ip] = struct{}{}
			allocatedIP := uint32ToIP(ip)
			return fmt.Sprintf("%s/%d", allocatedIP.String(), ones), nil
		}
	}
	return "", fmt.Errorf("no available IP addresses in the range")
}

func (i *IPMgr) release(ipStr string) error {
	ipAddr, _, err := net.ParseCIDR(ipStr)
	if err != nil {
		ipAddr = net.ParseIP(ipStr)
	}
	if ipAddr == nil {
		return fmt.Errorf("invalid IP address")
	}
	ip := ipToUint32(ipAddr)
	if _, exists := i.Used[ip]; !exists {
		return fmt.Errorf("IP address %s is not allocated", ipStr)
	}
	delete(i.Used, ip)
	return nil
}
