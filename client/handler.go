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
	for {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			plog.Errorf("failed to connect to server err-> %v", err)
			time.Sleep(time.Duration(3) * time.Second)
			continue
		}

		remoteAddr := conn.RemoteAddr().String()
		plog.Infof("inner client connected, remote-> %s", serverAddr)

		authSuccess := false

		// 处理连接
		c.handleConnection(conn, remoteAddr, &authSuccess)

		// 连接关闭后，尝试重新连接
		plog.Infof("attempting to reconnect to server in %d", 3)
		time.Sleep(time.Duration(3) * time.Second)
	}
}

// handleConnection 负责处理与服务器的连接和通信
func (c *plexClient) handleConnection(conn net.Conn, remoteAddr string, authSuccess *bool) {
	defer conn.Close()
	for {
		if !*authSuccess {
			// 发送认证消息
			if err := c.sendAuthMessage(conn); err != nil {
				plog.Errorf("inner client send auth message err-> %v", err)
				c.removeInnerClient(remoteAddr)
				return
			}
		}

		// 解包消息
		message, err := pack.Unpack(conn)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				plog.Infof("connection closed by server")
				c.removeInnerClient(remoteAddr)
				return
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
			*authSuccess = true

			// 启动心跳
			go c.startInnerClientHeartbeat(conn)

			// 添加客户端
			c.addInnerClient(conn)
		case auth_failed_uri:
			plog.Infof("inner client auth failed, check inner password")
			c.removeInnerClient(remoteAddr)
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
	ticker := time.NewTicker(time.Duration(c.cfg.Heartbeat) * time.Second)
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
