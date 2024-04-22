package plex

import (
	"github.com/swxctx/plex/plog"
)

// Config
type Config struct {
	// 监听端口
	Port string
	// 显示详细运行日志
	ShowTrace bool
	// 最大连接数
	MaxConnection int
	// 等待鉴权超时时间，超时未鉴权即会断开连接(默认20秒)
	AuthTimeout int64
	// 心跳时间间隔(默认60秒)
	Heartbeat int64
	// 心跳超时时间，超时即会断开连接(默认120秒)
	HeartbeatTimeout int64
	// 是否是小端字节序(默认是大端字节序)
	LittleEndian bool
}

// reloadConfig
func reloadConfig(cfgArg *Config) *Config {
	cfg := cfgArg
	if len(cfg.Port) <= 0 {
		cfg.Port = "9578"
	}
	if cfg.ShowTrace {
		plog.SetLevel("trace")
	}
	if cfg.AuthTimeout <= 0 {
		cfg.AuthTimeout = 20
	}
	if cfg.Heartbeat <= 0 {
		cfg.Heartbeat = 60
	}
	if cfg.HeartbeatTimeout <= 0 {
		cfg.HeartbeatTimeout = 120
	}

	return cfg
}
