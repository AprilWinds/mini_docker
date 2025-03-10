package cgroup

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const rootPath = "/sys/fs/cgroup"

func (c *Cgroup) Apply(pid int) error {
	if c.cfg.CPUS != 0.0 {
		return c.apply(pid, "cpu")
	}
	if c.cfg.MemoryLimit != 0 {
		return c.apply(pid, "memory")
	}
	return nil
}

func (c *Cgroup) Destroy() {
	os.RemoveAll(filepath.Join(rootPath, "cpu", c.cfg.Name))
	os.RemoveAll(filepath.Join(rootPath, "memory", c.cfg.Name))
}

func (c *Cgroup) apply(pid int, typeName string) error {
	cgroupPath := filepath.Join(rootPath, typeName, c.cfg.Name)
	os.MkdirAll(cgroupPath, 0755)

	defer func() error {
		processFile := filepath.Join(cgroupPath, "cgroup.procs")
		slog.Info("apply cgroup", "processFile", processFile, "pid", pid)
		return os.WriteFile(processFile, []byte(strconv.Itoa(pid)), 0644)
	}()

	if typeName == "cpu" {
		period, err := os.ReadFile(filepath.Join(cgroupPath, "cpu.cfs_period_us"))
		if err != nil {
			return fmt.Errorf("read cpu.cfs_period_us failed: %w", err)
		}
		periodInt, err := strconv.ParseInt(strings.TrimSpace(string(period)), 10, 32)
		if err != nil {
			return fmt.Errorf("parse cpu.cfs_period_us failed: %w", err)
		}

		quotaNum := int64(c.cfg.CPUS * float32(periodInt))
		cpuFile := filepath.Join(cgroupPath, "cpu.cfs_quota_us")

		if err := os.WriteFile(cpuFile, []byte(strconv.FormatInt(quotaNum, 10)), 0644); err != nil {
			return fmt.Errorf("write cpu.cfs_quota_us failed: %w", err)
		}
	}

	if typeName == "memory" {
		memLimitFile := filepath.Join(cgroupPath, "memory.limit_in_bytes")
		return os.WriteFile(memLimitFile, []byte(strconv.FormatUint(c.cfg.MemoryLimit, 10)), 0644)
	}

	return nil
}
