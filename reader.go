package plex

import (
	"io"
	"net"
	"time"

	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
)

// startReaderRoutine
func (c *plexConnection) startReaderRoutine() {
	remoteAddr := c.storeInfo.conn.RemoteAddr().String()
	c.remoteAddr = remoteAddr
	plog.Infof("plex server accept conn-> %s", remoteAddr)

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

		plog.Tracef("server receiver remote-> %s, msg-> %#v", c.remoteAddr, message)
		c.receiveMsgHandler(message)
	}
}

// receiveMsgHandler
func (c *plexConnection) receiveMsgHandler(message *pack.Message) {
	switch message.URI {
	case auth_uri:
		// auth
		authSuccess, uid := c.plexServer.authFunc(message.Body)
		if !authSuccess {
			plog.Errorf("auth failed, remote-> %s, body-> %d", c.remoteAddr, message.Body)
			// auth failed
			c.responseOnlyUri(auth_failed_uri)
			return
		}

		plog.Infof("auth success, remote-> %s, body-> %s", c.remoteAddr, message.Body)
		// save conn cache
		c.plexServer.store.set(uid, c.storeInfo)
		c.uid = uid

		// heartbeat timeout
		c.setReadDeadline(c.plexServer.cfg.HeartbeatTimeout)

		// auth success
		c.responseOnlyUri(auth_success_uri)
		break
	case heartbeat_uri:
		// heartbeat
		c.storeInfo.heartbeat = time.Now().Unix()

		if len(c.uid) > 0 {
			c.plexServer.store.set(c.uid, c.storeInfo)
		}

		// next heartbeat timeout
		c.setReadDeadline(c.plexServer.cfg.HeartbeatTimeout)

		// response success
		go c.responseOnlyUri(heartbeat_uri)
		break
	case inner_auth_uri:
		// inner tcp
		if message.Body != c.plexServer.cfg.InnerPassword {
			plog.Errorf("inner auth failed, remote-> %s, body-> %d", c.remoteAddr, message.Body)
			// auth failed
			c.responseOnlyUri(auth_failed_uri)
			return
		}

		plog.Infof("inner auth success, remote-> %s, body-> %s", c.remoteAddr, message.Body)

		// heartbeat timeout
		c.setReadDeadline(c.plexServer.cfg.HeartbeatTimeout)

		// auth success
		c.responseOnlyUri(auth_success_uri)
		break
	case send_msg_uri:
		// send msg to client
		go c.plexServer.sendMessageOuterClient(message)
		break
	}
}

// setReadDeadline
func (c *plexConnection) setReadDeadline(duration int64) {
	c.storeInfo.conn.SetReadDeadline(time.Now().Add(time.Duration(duration) * time.Second))
}
