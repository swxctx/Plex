package plex

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/swxctx/plex/plog"
)

// outerServerData
type outerServerData struct {
	Host string `json:"host,omitempty"`
}

// startHttpServer
func (s *plexServer) startHttpServer() {
	plog.Infof("plex http server is starting...")

	http.HandleFunc("/plex/v1/host", s.hostHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.HttpPort), nil); err != nil {
		plog.Errorf("start http listen err-> %v", err)
	}

	plog.Infof("plex http server is started...")
}

// hostHandler
func (s *plexServer) hostHandler(w http.ResponseWriter, r *http.Request) {
	plog.Tracef("hostHandler get hosts, ip-> %s, method-> %s", r.RemoteAddr, r.Method)

	// method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cfgHosts := s.cfg.OuterServers
	hosts := outerServerData{
		Host: shuffleHosts(cfgHosts),
	}

	// response
	data, err := json.Marshal(hosts)
	if err != nil {
		plog.Errorf("hostHandler: Marshal hosts err-> %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(data)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// shuffleHosts
func shuffleHosts(hosts []string) string {
	return hosts[rand.Intn(len(hosts))]
}
