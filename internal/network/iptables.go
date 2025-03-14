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

func setSNAT(rawIp string) error {
	args := fmt.Sprintf("-t nat -A POSTROUTING -s %s -o eth0 -j MASQUERADE", rawIp)
	return execIptablesCMD(args)
}

func deleteSNAT(rawIp string) error {
	args := fmt.Sprintf("-t nat -D POSTROUTING -s %s -o eth0 -j MASQUERADE", rawIp)
	return execIptablesCMD(args)
}

func setDNAT(srcPort string, destIp string, destPort string) error {
	args := fmt.Sprintf("-t nat -A PREROUTING -p tcp  --dport %s -j DNAT --to-destination %s:%s", srcPort, destIp, destPort)
	if err := execIptablesCMD(args); err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}
	// 本机进程走的output链
	args = fmt.Sprintf("-t nat -A OUTPUT -p tcp  --dport %s -j DNAT --to-destination %s:%s", srcPort, destIp, destPort)
	if err := execIptablesCMD(args); err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}
	return nil
}

func deleteDNAT(srcPort string, destIp string, destPort string) error {
	args := fmt.Sprintf("-t nat -D PREROUTING -p tcp  --dport %s -j DNAT --to-destination %s:%s", srcPort, destIp, destPort)
	if err := execIptablesCMD(args); err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}

	args = fmt.Sprintf("-t nat -D OUTPUT -p tcp  --dport %s -j DNAT --to-destination %s:%s", srcPort, destIp, destPort)
	if err := execIptablesCMD(args); err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}
	return nil
}

func setForward(bridgeName string) error {
	args := fmt.Sprintf("-A FORWARD -i %s -o eth0 -j ACCEPT", bridgeName)
	if err := execIptablesCMD(args); err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}
	args = fmt.Sprintf("-A FORWARD -i eth0 -o %s -j ACCEPT", bridgeName)
	if err := execIptablesCMD(args); err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}
	return nil
}

func deleteForward(bridgeName string) error {
	args := fmt.Sprintf("-D FORWARD -i %s -o eth0 -j ACCEPT", bridgeName)
	if err := execIptablesCMD(args); err != nil {
		return fmt.Errorf("failed to set iptables: %w", err)
	}
	args = fmt.Sprintf("-D FORWARD -i eth0 -o %s -j ACCEPT", bridgeName)
	return execIptablesCMD(args)
}
