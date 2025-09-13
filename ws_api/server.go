package ws_api

import (
	"fmt"
	"net/http"

	"github.com/pterm/pterm"
	"github.com/vmihailenco/msgpack/v5"
)

func StartWSServer(port int) {
	// 开启LanGame的 WebSocket 接口服务
	http.HandleFunc("/", handleWSCliConnection)
	go handleBroadcasts()

	pterm.Info.Println("已在", port, "端口开放 WebSocket 接口")
	pterm.Error.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// 处理来自单个 WebSocket Client 的连接
func handleWSCliConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		pterm.Error.Println("客户端连接预处理出错:", err)
		return
	}
	defer conn.Close()

	client := &WSAPIClient{conn: conn, ready: false}
	clis[client] = true

	pterm.Info.Println("客户端", conn.RemoteAddr().String(), "已连接到LanGame")

	if botDatasReady {
		sendBotBasicIDAndSetClientReady(client)
	} else {
		pterm.Info.Println("客户端", conn.RemoteAddr().String(), "正在等待玩家信息初始化")
	}
	sendUpdateUQ(client)

	// 读取消息
	for {
		var msg Message
		_, content, err := conn.ReadMessage()
		if err != nil {
			pterm.Warning.Println("客户端", conn.RemoteAddr().String(), "连接中断:", err)
			delete(clis, client)
			break
		}
		err = msgpack.Unmarshal(content, &msg)
		if err != nil {
			pterm.Warning.Println("客户端", conn.RemoteAddr().String(), "读取 msgpack 出现问题:", err)
			delete(clis, client)
			break
		} else {
			// 暂时将所有消息统一处理
			receiver <- msg
		}
	}
}
