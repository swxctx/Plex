package client

import (
	"github.com/swxctx/plex/plog"
)

// Config
type Config struct {
	// 多机部署，服务器地址["127.0.0.1:9587", "127.0.0.1:9588"]
	InnerServers []string
	// 内部通信密码
	InnerPassword string
	// 显示详细运行日志
	ShowTrace bool
	// 是否是小端字节序(默认是大端字节序)
	LittleEndian bool
}

// reloadConfig
func reloadConfig(cfgArg *Config) *Config {
	cfg := cfgArg

	if len(cfg.InnerServers) <= 0 {
		cfg.InnerServers = []string{"127.0.0.1:9578"}
	}

	if len(cfg.InnerPassword) <= 0 {
		cfg.InnerPassword = "plex-inner"
	}

	if cfg.ShowTrace {
		plog.SetLevel("trace")
	}

	plog.Infof("--- config start ----")
	plog.Infof("InnerServers: %v", cfg.InnerServers)
	plog.Infof("InnerPassword: %v", cfg.InnerPassword)
	plog.Infof("ShowTrace: %v", cfg.ShowTrace)
	plog.Infof("LittleEndian(s): %d", cfg.LittleEndian)
	plog.Infof("--- config end ----")
	return cfg
}
