package client

import (
	"io"
	"net"
	"time"

	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
)

// innerClient
type innerClient struct {
	// client id
	id string
	// conn to plex server
	conn net.Conn
	// send message channel
	messageChan chan string
}

// startInnerClient
func (c *plexClient) startInnerClient() {
	plog.Infof("starting inner client...")
	for _, sd := range c.cfg.InnerServers {
		go c.innerClientConnection(sd)
	}
	plog.Infof("inner client all started...")
}

// innerClientConnection
func (c *plexClient) innerClientConnection(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		plog.Errorf("failed to connect to server err-> %v", err)
		return
	}
	defer conn.Close()

	plog.Infof("inner client connected, remote-> %s", serverAddr)

	authSuccess := false

	// read data
	for {
		if !authSuccess {
			// auth
			if err := c.sendAuthMessage(conn); err != nil {
				plog.Errorf("inner client send auth message err-> %v", err)
				return
			}
		}

		// unpack message
		message, err := pack.Unpack(conn)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				plog.Infof("connection closed by server")
				break
			}
			plog.Errorf("inner client unpack err-> %v", err)
			continue
		}
		if message == nil {
			continue
		}

		plog.Tracef("inner client receive message-> %#v", message)

		switch message.URI {
		case auth_success_uri:
			plog.Infof("inner client auth success")
			authSuccess = true

			// heartbeat
			go c.startInnerClientHeartbeat(conn)

			// add inner client
			c.addInnerClient(conn)
		case auth_failed_uri:
			plog.Infof("inner client auth failed, check inner password")
			return
		}
	}
}

// sendAuthMessage
func (c *plexClient) sendAuthMessage(conn net.Conn) error {
	writeData, err := pack.Pack(&pack.Message{
		URI:  inner_auth_uri,
		Body: c.cfg.InnerPassword,
	})
	if err != nil {
		return err
	}
	_, err = conn.Write(writeData)
	return err
}

// startInnerClientHeartbeat
func (c *plexClient) startInnerClientHeartbeat(conn net.Conn) {
	msg, err := pack.Pack(&pack.Message{
		URI: heartbeat_uri,
	})
	if err != nil {
		plog.Errorf("start inner client heartbeat err-> %v", err)
		return
	}

	// time ticket
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if _, err := conn.Write(msg); err != nil {
				plog.Errorf("inner client heartbeat send err-> %v", err)
				return
			}
		}
	}
}
