package plex

import (
	"fmt"
	"github.com/swxctx/plex/plog"
	"net"
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
	authFunc func(body string) bool
}

// NewServer
func NewServer(config *Config, fn ...func(body string) bool) {
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
}

// Start
func Start() error {
	plog.Infof("start plex server.")

	if server.authFunc == nil {
		return fmt.Errorf("plex auth func is nil")
	}

	// tcp listen
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.cfg.Port))
	if err != nil {
		plog.Errorf("plex tcp server listen err-> %v", err)
		return fmt.Errorf("plex listen err-> %v", err)
	}
	defer listener.Close()

	plog.Infof("plex server is starting...")
	plog.Infof("start accept job.")

	for {
		// listen and accept
		conn, err := listener.Accept()
		if err != nil {
			plog.Errorf("listener accept err-> %v", err)
			continue
		}

		// start conn logic
		server.newPlexConnection(conn)
	}
}

// SetAuthFunc
func SetAuthFunc(fn func(body string) bool) {
	server.authFunc = fn
}
