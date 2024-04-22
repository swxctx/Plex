package main

import (
	"github.com/swxctx/plex"
	"github.com/swxctx/plex/plog"
)

func main() {
	// new server
	plex.NewServer(&plex.Config{
		ShowTrace: true,
	})

	// start
	if err := plex.Start(); err != nil {
		plog.Errorf("start err-> %v", err)
		return
	}
}
