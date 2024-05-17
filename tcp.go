package plex

import (
	"fmt"
	"net"
	"time"

	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
)

// startTcpServer
func (s *plexServer) startTcpServer() {
	plog.Infof("start plex server.")

	if s.authFunc == nil {
		panic(fmt.Errorf("plex auth func is nil"))
	}

	// tcp listen
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		plog.Errorf("plex tcp server listen err-> %v", err)
		panic(err)
	}
	defer listener.Close()

	plog.Infof("plex tcp server is starting...")
	plog.Infof("start tcp accept job.")

	for {
		// listen and accept
		conn, err := listener.Accept()
		if err != nil {
			plog.Errorf("listener accept err-> %v", err)
			continue
		}

		// start conn logic
		s.newPlexConnection(conn)
	}
}

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

// sendMessageOuterClient
func (s *plexServer) sendMessageOuterClient(innerMessage *pack.Message) {
	// unpack outer message
	outerMessage := unmarshalSendMessage(innerMessage.Body)
	if outerMessage == nil {
		return
	}

	// get client cache
	store, exists := s.store.get(outerMessage.Uid)
	if !exists {
		return
	}
	if store == nil || store.conn == nil {
		return
	}

	// pack message
	message, err := pack.Pack(&pack.Message{
		Seq:  GetSeq(),
		URI:  outerMessage.Uri,
		Body: outerMessage.Body,
	})
	if err != nil {
		plog.Errorf("responseOuterClient: err-> %v, uri-> %s", err, outerMessage.Uri)
		return
	}

	// send message
	store.send(message)
}
