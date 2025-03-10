package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"mini_docker/internal/fs"
)

func getContainer(containerID string) (*container, error) {
	rootDir := getRootDir(containerID)
	metadata, err := os.ReadFile(getMatadataFile(containerID))
	if err != nil {
		fs.UnmountOverlay(rootDir)
		os.RemoveAll(rootDir)
		return nil, fmt.Errorf("read metadata failed: %w", err)
	}
	var c container
	if err := json.Unmarshal(metadata, &c); err != nil {
		return nil, fmt.Errorf("unmarshal metadata failed: %w", err)
	}
	return &c, nil
}

// 获取容器根目录
func getRootDir(containerID string) string {
	return filepath.Join(StorageRootDir, containerID)
}

func getLogFile(containerID string) string {
	return filepath.Join(StorageRootDir, containerID, "logs")
}

func getMatadataFile(containerID string) string {
	return filepath.Join(StorageRootDir, containerID, "metadata.json")
}

func getWorkDir(containerID string) string {
	return filepath.Join(StorageRootDir, containerID, "merged")
}

func getReadOnlyDir(containerID string) string {
	return filepath.Join(StorageRootDir, containerID, "lower")
}
func getMountDir(containerID string, mountDir string) string {
	return filepath.Join(StorageRootDir, containerID, "merged", mountDir)
}
