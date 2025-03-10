package util

import (
	"fmt"
	"net"
)

func GetCIDRGateway(subnet string) string {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		LogAndExit("failed to parse subnet", err)
	}
	ipnet.IP = ipnet.IP.To4()

	if ipnet.IP[3] == 0 {
		ipnet.IP[3] = 1
	}
	mask, _ := net.IPMask(ipnet.Mask).Size()
	return fmt.Sprintf("%s/%d", ipnet.IP.String(), mask)
}
