package plex

import (
	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
	"net"
	"time"
)

// plexConnection
type plexConnection struct {
	// 用户标识
	uid string
	// conn cache info
	storeInfo *storeInfo
	// bind server info
	plexServer *plexServer
}

// storeInfo
type storeInfo struct {
	// connection client
	conn net.Conn
	// last heartbeat ts
	heartbeat int64
}

// send
func (c *plexConnection) send(data []byte) {
	c.storeInfo.conn.Write(data)
}

// close connection
func (c *plexConnection) close() {
	// close conn
	c.storeInfo.conn.Close()

	// del store
	c.plexServer.store.Del(c.uid)
}

// responseOnlyUri
func (c *plexConnection) responseOnlyUri(uri string) {
	// pack message
	message, err := pack.Pack(&pack.Message{
		Seq: time.Now().Unix(),
		URI: uri,
	})
	if err != nil {
		plog.Errorf("responseAuth: err-> %v, uri-> %s", err, uri)
		return
	}

	// send message
	c.send(message)
}
