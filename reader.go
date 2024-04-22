package plex

import (
	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
	"net"
)

// startReaderRoutine
func (s *plexServer) startReaderRoutine(conn net.Conn) {
	/*// 开启startReader时，默认是没有授权的，同时也不做存储
	authTimeout := time.Now().Add(time.Duration(s.cfg.AuthTimeout) * time.Second)

	// 设置超时时间
	conn.SetReadDeadline(authTimeout)*/

	for {
		// unpack message
		message, err := pack.Unpack(conn)
		if err != nil {
			continue
		}
		if message == nil {
			continue
		}
		plog.Infof("receiver msg-> %#v", message)
	}
}
