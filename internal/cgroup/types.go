package cgroup

type Config struct {
	Name        string
	MemoryLimit uint64  // 字节
	CPUS        float32 // 占用cpu的百分比
}

type Cgroup struct {
	cfg *Config
}

func New(cfg *Config) *Cgroup {
	cg := &Cgroup{
		cfg: cfg,
	}
	return cg
}
