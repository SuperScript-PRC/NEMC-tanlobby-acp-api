package ws_api

import (
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type WSAPIClient struct {
	// 来自客户端的连接
	conn *websocket.Conn
	// 客户端通信是否已就绪。
	// 如果未就绪, LanGame会尝试向客户端更新机器人基本信息以及 UQHolder 信息。
	ready bool
}

func (cli *WSAPIClient) SendMessage(msg Message) error {
	if cli.conn != nil {
		bts, err := msgpack.Marshal(msg)
		if err != nil {
			return err
		}
		return cli.conn.WriteMessage(websocket.BinaryMessage, bts)
	} else {
		return nil
	}
}

// 设置 WebSocket 客户端的状态为已初始化
func (cli *WSAPIClient) Ready() {
	cli.ready = true
}
