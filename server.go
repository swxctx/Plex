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
}

// NewServer
func NewServer(config *Config) {
	// reload config
	cfg := reloadConfig(config)

	plog.Infof("reload config success.")

	// new server
	server = &plexServer{
		cfg:   cfg,
		store: newConnStore(cfg.MaxConnection),
	}
	plog.Infof("new plex server success.")
}

// Start
func Start() error {
	plog.Infof("start plex server.")

	// tcp listen
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.cfg.Port))
	if err != nil {
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

		// 开始读取监听
		server.startReaderRoutine(conn)
	}
}
