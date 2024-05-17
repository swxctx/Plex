package plex

import (
	"net"

	"github.com/swxctx/plex/plog"
)

var (
	server *plexServer
)

// plexServer
type plexServer struct {
	// config
	cfg *Config
	// tcp listen
	listener net.Listener
	// conn store cache
	store *connStore
	// tcp auth func
	authFunc func(body string) (bool, string)
}

// Start
func Start(config *Config, fn ...func(body string) (bool, string)) {
	// reload config
	cfg := reloadConfig(config)

	plog.Infof("reload config success.")

	// new server
	server = &plexServer{
		cfg:   cfg,
		store: newConnStore(cfg.MaxConnection),
	}

	// auth handler
	if len(fn) > 0 {
		server.authFunc = fn[0]
	}
	plog.Infof("new plex server success.")

	plog.Infof("--- server start begin ---")

	// start http
	go server.startHttpServer()

	// start tcp
	server.startTcpServer()

	plog.Infof("--- server start end ---")
}

// SetAuthFunc func set outer auth func
func SetAuthFunc(fn func(body string) (bool, string)) {
	server.authFunc = fn
}
