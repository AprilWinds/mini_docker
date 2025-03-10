package container

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"mini_docker/internal/cgroup"
	"mini_docker/internal/fs"
	"mini_docker/internal/util"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func (c *container) startChild() {
	reader, writer, err := os.Pipe()
	if err != nil {
		util.LogAndExit("create pipe failed:", err)
	}

	child := exec.Command("/proc/self/exe", "setup")
	child.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	child.ExtraFiles = []*os.File{reader}
	if c.Setting.It {
		child.Stdin = os.Stdin
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr
	} else {
		logFile, err := os.OpenFile(getLogFile(c.Id), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			util.LogAndExit("open log file failed:", err)
		}
		child.Stdout = logFile
		child.Stderr = logFile
	}
	child.Dir = getWorkDir(c.Id)
	child.Env = c.Setting.Env
	if err := child.Start(); err != nil {
		util.LogAndExit("start child failed:", err)
	}

	c.child = child
	c.PId = child.Process.Pid
	childCommand := strings.Join(c.Setting.CMD, " ")
	writer.WriteString(childCommand)
	writer.Close()
	slog.Info("start child", "id", c.Id, "pid", c.PId)
}

func (c *container) loadImage() {
	imageTar := "./" + c.Setting.ImageName + ".tar.gz"
	lowerDir := getReadOnlyDir(c.Id)
	_, err := exec.Command("tar", "-xzf", imageTar, "-C", lowerDir).CombinedOutput()
	if err != nil {
		util.LogAndExit("tar -xzf image failed:", err)
	}
}

func (c *container) setupOverlayFS() {
	fs.MountOverlay(getRootDir(c.Id))
	if len(c.Setting.Volume) == 2 {
		hostPath := c.Setting.Volume[0]
		targetPath := getMountDir(c.Id, c.Setting.Volume[1])
		fs.BindMount(hostPath, targetPath)
	}
}

func (c *container) saveMetadata() {
	outFile := getMatadataFile(c.Id)
	data, err := json.Marshal(c)
	if err != nil {
		util.LogAndExit("json marshal failed:", err)
	}
	os.WriteFile(outFile, data, 0644)
}

func Run(s *Setting) {
	id := uuid.New()
	c := &container{
		Id:         base64.RawURLEncoding.EncodeToString(id[:]),
		CreateTime: time.Now(),
		Setting:    s,
	}
	c.setupOverlayFS()
	c.loadImage()
	c.startChild()
	c.saveMetadata()

	cg := cgroup.New(s.CgroupCfg)
	if err := cg.Apply(c.PId); err != nil {
		util.LogAndExit("apply cgroup failed:", err)
	}

	if c.Setting.It {
		c.child.Wait()
	}
}
