package main

import (
	"encoding/json"
	"net/http"

	"github.com/swxctx/plex/client"
	"github.com/swxctx/plex/plog"
)

func main() {
	client.Start(&client.Config{
		InnerServers: []string{"10.60.76.224:9578"},
		ShowTrace:    true,
	})
	plog.Infof("client started...")

	startHttpServer()
	plog.Infof("logic send server api started...")
}

// startHttpServer
func startHttpServer() {
	plog.Infof("plex http server is starting...")

	http.HandleFunc("/plexApi/v1/send", apiHandler)
	if err := http.ListenAndServe("0.0.0.0:9501", nil); err != nil {
		plog.Errorf("start http listen err-> %v", err)
	}

	plog.Infof("plex api is started...")
}

// RequestData
type RequestData struct {
	Uid  string `json:"uid"`
	Body string `json:"body"`
	Uri  string `json:"uri"`
}

// hostHandler
func apiHandler(w http.ResponseWriter, r *http.Request) {
	plog.Tracef("hostHandler get hosts, ip-> %s, method-> %s", r.RemoteAddr, r.Method)

	var (
		data RequestData
	)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		// error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	client.Send(&client.SendMessageArgs{
		Uid:  data.Uid,
		Body: data.Body,
		Uri:  data.Uri,
	})
	w.Write([]byte("{\"msg\": \"success\"}"))
}
