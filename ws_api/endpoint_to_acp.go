package ws_api

import (
	"encoding/json"
	"net/http"

	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"
	"github.com/pterm/pterm"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/gorilla/websocket"
)

var clis = make(map[*WSAPIClient]bool)
var broadcaster = make(chan Message)
var receiver = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 向所有 WebSocket Client 广播消息
func BroadcastMessageToEndpoints(msg Message) {
	broadcaster <- msg
}

// 统一处理来自所有 WebSocket Client 的通信消息
func HandleEndpointsMessages(
	writePacket func(pk packet.Packet) error,
) {
	var writeErr error
	for {
		msg := <-receiver
		content, ok := msg.Content.(map[string]any)
		if !ok {
			pterm.Error.Println("无效 WS 信息:", msg.Type)
			continue
		}
		switch msg.Type {
		case WSMSG_SERVER_PACKET:
			pk1, err := ConstructPacketFromContent(content, false, msgpack.Unmarshal)
			if err != nil {
				pterm.Error.Printfln("%v", err)
				continue
			}
			writeErr = writePacket(pk1)
		case WSMSG_SERVER_PACKET_JSON:
			pk1, err := ConstructPacketFromContent(content, false, json.Unmarshal)
			if err != nil {
				pterm.Error.Printfln("%v", err)
				continue
			}
			writeErr = writePacket(pk1)
		case WSMSG_SERVER_CUSTOM_PACKET:
			pk1, err := ConstructPacketFromContent(content, true, msgpack.Unmarshal)
			if err != nil {
				pterm.Error.Printfln("%v", err)
				continue
			}
			writeErr = writePacket(pk1)
		case WSMSG_SET_SERVER_LISTEN_PACKETS:
			pkIDs, ok := convertToIntArr(content["PacketsID"])
			if !ok {
				pterm.Error.Println("无法识别监听数据包请求:", content["PacketsID"])
				continue
			}
			// pterm.Info.Println("设置监听服务端的数据包:", pkIDs)
			SetServerToClientListenPackets(pkIDs)
		default:
			pterm.Warning.Println("无效的消息类型:", msg.Type)
		}
		if writeErr != nil {
			pterm.Error.Println("writePacketErr:", writeErr)
		}
	}
}

// 处理来自LanGame向所有 WSAPIClient 广播的消息
func handleBroadcasts() {
	for {
		msg := <-broadcaster
		for client := range clis {
			err := client.SendMessage(msg)
			if err != nil {
				pterm.Error.Println(err)
				client.conn.Close()
				delete(clis, client)
			}

		}
	}
}

// 分发LanGame机器人自身的基本信息, 如玩家名, UQ 信息等
// 然后设置所有 WebSocket 客户端连接状态为已初始化
func handoutBotBasicInfo() {
	for cli := range clis {
		if !cli.ready {
			sendBotBasicIDAndSetClientReady(cli)
		}
	}
}

func sendBotBasicIDAndSetClientReady(cli *WSAPIClient) {
	// 向 WebSocket 客户端发送 BotBasicID 并使其得到初始化
	cli.SendMessage(Message{
		Type:    WSMSG_SET_BOT_BASIC_INFO,
		Content: getBotBasicInfo(),
	})
	cli.Ready()
}

func sendUpdateUQ(cli *WSAPIClient) {
	// 向 WebSocket 客户端发送全局玩家 UQ 更新信息
	pterm.Info.Println("向", cli.conn.RemoteAddr().String(), "发送全局玩家UQ更新信息")
	cli.SendMessage(Message{
		Type:    WSMSG_UPDATE_UQ,
		Content: simple_uq_map,
	})
}
