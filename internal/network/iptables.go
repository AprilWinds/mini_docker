package network

import (
	"fmt"
	"os/exec"
	"strings"
)

func execIptablesCMD(args string) error {
	cmd := exec.Command("iptables", strings.Split(args, " ")...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}
	return nil
}

func setHostSNAT(bridgeName string, rawIp string) error {
	args := fmt.Sprintf("-t nat -A POSTROUTING -s %s -o %s -j MASQUERADE", rawIp, bridgeName)
	return execIptablesCMD(args)
}

func deleteHostSNAT(bridgeName string, rawIp string) error {
	args := fmt.Sprintf("-t nat -D POSTROUTING -s %s -o %s -j MASQUERADE", rawIp, bridgeName)
	return execIptablesCMD(args)
}

func setContainerMapping(srcIp string, dstIp string) error {
	args := fmt.Sprintf("-t nat -A PREROUTING -d %s -j DNAT --to-destination %s", srcIp, dstIp)
	return execIptablesCMD(args)
}
