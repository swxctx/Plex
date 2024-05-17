package plex

import (
	"fmt"
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
		// message logic
		if err := c.handleMessage(); err != nil {
			if err == io.EOF {
				plog.Tracef("connection closed, remote-> %s, uid-> %s", c.remoteAddr, c.uid)
				break
			}

			// ReadDeadline timeout
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				plog.Warnf("connection timeout, remote-> %s", remoteAddr)
				c.close()
				break
			}
			plog.Errorf("error handling message: %s", err)
			continue
		}
	}
}

// handleMessage
func (c *plexConnection) handleMessage() error {
	message, err := pack.Unpack(c.storeInfo.conn)
	if err != nil {
		return err
	}
	if message == nil {
		return nil
	}

	plog.Tracef("Server received message from remote-> %s, msg-> %#v", c.remoteAddr, message)
	return c.receiveMsgHandler(message)
}

// receiveMsgHandler
func (c *plexConnection) receiveMsgHandler(message *pack.Message) error {
	if !c.isAuthenticated && message.URI != auth_uri && message.URI != inner_auth_uri {
		plog.Errorf("unauthorized access attempt, remote-> %s", c.remoteAddr)
		// auth failed
		err := c.responseOnlyUri(auth_failed_uri)
		if err != nil {
			return err
		}
		return fmt.Errorf("unauthorized access")
	}

	switch message.URI {
	case auth_uri:
		// auth
		return c.handlerAuth(message)
	case heartbeat_uri:
		// heartbeat
		return c.handleHeartbeat(message)
	case inner_auth_uri:
		// inner tcp
		return c.handlerInnerAuth(message)
	case send_msg_uri:
		// send msg to client
		return c.handleSendMessage(message)
	}
	return nil
}

// setReadDeadline
func (c *plexConnection) setReadDeadline(duration int64) {
	c.storeInfo.conn.SetReadDeadline(time.Now().Add(time.Duration(duration) * time.Second))
}

// handlerAuth
func (c *plexConnection) handlerAuth(message *pack.Message) error {
	// already auth
	if c.isAuthenticated {
		plog.Infof("already authenticated, remote-> %s", c.remoteAddr)
		return nil
	}

	// auth func
	authSuccess, uid := c.plexServer.authFunc(message.Body)
	if !authSuccess {
		plog.Errorf("auth failed, remote-> %s, body-> %s", c.remoteAddr, message.Body)
		c.responseOnlyUri(auth_failed_uri)
		return fmt.Errorf("auth failed")
	}
	c.isAuthenticated = true

	plog.Infof("auth success, remote-> %s, body-> %s", c.remoteAddr, message.Body)

	// save conn cache
	c.plexServer.store.set(uid, c.storeInfo)
	c.uid = uid

	// heartbeat timeout
	c.setReadDeadline(c.plexServer.cfg.HeartbeatTimeout)

	// auth success
	c.responseOnlyUri(auth_success_uri)

	return nil
}

// handlerInnerAuth
func (c *plexConnection) handlerInnerAuth(message *pack.Message) error {
	// already auth
	if c.isAuthenticated {
		plog.Infof("already authenticated, remote-> %s", c.remoteAddr)
		return nil
	}

	// auth func
	if message.Body != c.plexServer.cfg.InnerPassword {
		plog.Errorf("inner auth failed, remote-> %s, body-> %d", c.remoteAddr, message.Body)
		// auth failed
		c.responseOnlyUri(auth_failed_uri)
		return fmt.Errorf("inner auth failed")
	}
	c.isAuthenticated = true

	plog.Infof("inner auth success, remote-> %s, body-> %s", c.remoteAddr, message.Body)

	// heartbeat timeout
	c.setReadDeadline(c.plexServer.cfg.HeartbeatTimeout)

	// auth success
	return c.responseOnlyUri(auth_success_uri)
}

// handlerAuth
func (c *plexConnection) handleHeartbeat(message *pack.Message) error {
	// heartbeat
	c.storeInfo.heartbeat = time.Now().Unix()

	if len(c.uid) > 0 {
		c.plexServer.store.set(c.uid, c.storeInfo)
	}

	// next heartbeat timeout
	c.setReadDeadline(c.plexServer.cfg.HeartbeatTimeout)

	// response success
	return c.responseOnlyUri(heartbeat_uri)
}

// handlerAuth
func (c *plexConnection) handleSendMessage(message *pack.Message) error {
	go c.plexServer.sendMessageOuterClient(message)
	return nil
}
