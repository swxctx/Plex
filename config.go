package plex

import (
	"github.com/swxctx/plex/plog"
)

// Config
type Config struct {
	// TCP监听端口
	Port string
	// Http api 监听端口号
	HttpPort string
	// 多机部署，服务器地址["123.123.45.67:9587", "123.123.45.68:9588"]
	OuterServers []string
	// 内部通信密码
	InnerPassword string
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

	if len(cfg.HttpPort) <= 0 {
		cfg.HttpPort = "9579"
	}

	if len(cfg.OuterServers) <= 0 {
		cfg.OuterServers = []string{"127.0.0.1:" + cfg.Port}
	}

	if len(cfg.InnerPassword) <= 0 {
		cfg.InnerPassword = "plex-inner"
	}

	if cfg.ShowTrace {
		plog.SetLevel("trace")
	}
	if cfg.AuthTimeout <= 0 {
		cfg.AuthTimeout = 30
	}
	if cfg.Heartbeat <= 0 {
		cfg.Heartbeat = 60
	}
	if cfg.HeartbeatTimeout <= 0 {
		cfg.HeartbeatTimeout = 120
	}

	plog.Infof("--- config start ----")
	plog.Infof("Port: %s", cfg.Port)
	plog.Infof("HttpPort: %s", cfg.HttpPort)
	plog.Infof("OuterServers: %v", cfg.OuterServers)
	plog.Infof("InnerPassword: %v", cfg.InnerPassword)
	plog.Infof("ShowTrace: %v", cfg.ShowTrace)
	plog.Infof("AuthTimeout(s): %d", cfg.AuthTimeout)
	plog.Infof("Heartbeat(s): %d", cfg.Heartbeat)
	plog.Infof("HeartbeatTimeout(s): %d", cfg.HeartbeatTimeout)
	plog.Infof("--- config end ----")
	return cfg
}
