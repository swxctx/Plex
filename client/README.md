## Plex Client

- `Plex Client` 为业务端提供快捷的消息发送客户端

- 业务端可以快速集成 `Plex Client`，之后可以快速调用方法进行消息发送

## 配置说明

```yaml
type Config struct {
	// 多机部署，服务器地址["127.0.0.1:9587", "127.0.0.1:9588"]
	InnerServers []string
	// 内部通信密码，需要与PlexServer保持一致
	InnerPassword string
	// 心跳时间间隔(默认60秒)
	Heartbeat int64
	// 显示详细运行日志，默认false
	ShowTrace bool
	// 是否是小端字节序(默认是大端字节序)，需要与PlexServer保持一致
	LittleEndian bool
}
```

## 使用案例

- 如下代码所示，代码启动了一个 `HTTP` 服务监听，同时初始化了 `PlexClient`
- 调用接口 `/plexApi/v1/send` 后服务调用了 `PlexClient` 的消息发送
- 这里只是一个简单的演示案例，在实际应用中可以与自己的服务架构相结合使用


```
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
```