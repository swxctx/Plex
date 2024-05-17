package plex

import (
	"net"

	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
)

// plexConnection
type plexConnection struct {
	// 用户标识
	uid string
	// remote
	remoteAddr string
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

// close connection
func (c *plexConnection) close() {
	// close conn
	c.storeInfo.conn.Close()

	// del store
	c.plexServer.store.del(c.uid)
}

// responseOnlyUri
func (c *plexConnection) responseOnlyUri(uri string) {
	// pack message
	message, err := pack.Pack(&pack.Message{
		Seq: GetSeq(),
		URI: uri,
	})
	if err != nil {
		plog.Errorf("responseOnlyUri: err-> %v, uri-> %s", err, uri)
		return
	}

	// send message
	c.storeInfo.send(message)
}

// send
func (s *storeInfo) send(data []byte) {
	s.conn.Write(data)
}
