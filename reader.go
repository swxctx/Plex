package plex

import (
	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
	"io"
	"net"
	"time"
)

// startReaderRoutine
func (s *plexServer) newPlexConnection(conn net.Conn) {
	// 开启startReader时，默认是没有授权的，同时也不做存储
	authTimeout := time.Now().Add(time.Duration(s.cfg.AuthTimeout) * time.Second)

	// 设置超时时间
	conn.SetReadDeadline(authTimeout)

	// new connection
	connection := plexConnection{
		plexServer: s,
		storeInfo: &storeInfo{
			conn:      conn,
			heartbeat: time.Now().Unix(),
		},
	}

	// start routine
	go connection.startReaderRoutine()
}

// startReaderRoutine
func (c *plexConnection) startReaderRoutine() {
	remoteAddr := c.storeInfo.conn.RemoteAddr().String()
	plog.Infof("accept conn-> %s", remoteAddr)

	// set auth timeout
	c.storeInfo.conn.SetReadDeadline(time.Now().Add(time.Duration(c.plexServer.cfg.AuthTimeout) * time.Second))

	for {
		// unpack message
		message, err := pack.Unpack(c.storeInfo.conn)
		if err != nil {
			// exception close
			if err == io.EOF {
				plog.Tracef("conn closed, remote-> %s, uid-> %s", remoteAddr, c.uid)
				c.close()
				break
			}

			// ReadDeadline timeout
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				plog.Warnf("auth timeout, remote-> %s", remoteAddr)
				c.close()
				break
			}
			continue
		}
		if message == nil {
			continue
		}

		plog.Tracef("receiver msg-> %#v", message)
	}
}

// receiveMsgHandler
func (c *plexConnection) receiveMsgHandler(message *pack.Message) {
	switch message.URI {
	case auth_uri:
		// auth
		if !c.plexServer.authFunc(message.Body) {
			// auth failed
			return
		}
		// auth success
		c.responseOnlyUri(auth_success_uri)
		break
	case heartbeat_uri:
		// heartbeat
	}
}
