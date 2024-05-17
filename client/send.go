package client

import (
	"net"

	"github.com/swxctx/plex"
	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
)

// addInnerClient
func (c *plexClient) addInnerClient(conn net.Conn) {
	clientArg := &innerClient{
		id:          conn.LocalAddr().String(),
		conn:        conn,
		messageChan: make(chan string),
	}

	c.mutex.Lock()
	c.innerClients = append(c.innerClients, clientArg)
	c.mutex.Unlock()

	// start write message job
	go c.startWriteMessagesToInnerServer(clientArg)
}

// removeInnerClient
func (c *plexClient) removeInnerClient(id string) {
	c.mutex.Lock()
	for i, ic := range c.innerClients {
		if ic.id == id {
			close(ic.messageChan)
			c.innerClients = append(c.innerClients[:i], c.innerClients[i+1:]...)
			break
		}
	}
	c.mutex.Unlock()
}

// broadcastMessage
func (c *plexClient) broadcastMessage(sendMessage string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, ct := range c.innerClients {
		ct.messageChan <- sendMessage
	}
}

// startWriteMessagesToInnerServer
func (c *plexClient) startWriteMessagesToInnerServer(clientArg *innerClient) {
	plog.Infof("start write message to inner server job ,id-> %s", clientArg.id)
	for msgStr := range clientArg.messageChan {
		// pack message
		msgByte, err := pack.Pack(packInnerSendMsg(msgStr))
		if err != nil {
			plog.Errorf("failed to pack message err-> %v, message-> %s", err, msgStr)
			continue
		}

		// send
		_, err = clientArg.conn.Write(msgByte)
		if err != nil {
			plog.Errorf("failed to send message to client err-> %v", err)
			continue
		}
	}
}

// packInnerSendMsg
func packInnerSendMsg(outerMsg string) *pack.Message {
	return &pack.Message{
		Seq:  plex.GetSeq(),
		URI:  send_msg_uri,
		Body: outerMsg,
	}
}
