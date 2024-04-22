package plex

import (
	"net"
)

// connection
type connection struct {
	// connection client
	conn net.Conn
	// last heartbeat ts
	heartbeat int64
}
