package client

import (
	"sync"

	"github.com/swxctx/plex/plog"
)

var (
	client *plexClient
)

// plexClient
type plexClient struct {
	// config
	cfg *Config
	// inner client
	innerClients []*innerClient
	// lock
	mutex sync.Mutex
}

// Start for start inner client
func Start(config *Config) {
	// reload config
	cfg := reloadConfig(config)

	plog.Infof("reload config success.")

	// new inner client
	client = &plexClient{
		cfg:          cfg,
		innerClients: make([]*innerClient, 0),
	}

	plog.Infof("new plex server success.")

	plog.Infof("--- client init begin ---")

	// inner client
	go client.startInnerClient()

	plog.Infof("--- client init end ---")
}

// Send func send message to client
func Send(sendMessage *SendMessageArgs) error {
	msg, err := marshalSendMessage(sendMessage)
	if err != nil {
		return err
	}

	// publish
	client.broadcastMessage(msg)

	return nil
}
