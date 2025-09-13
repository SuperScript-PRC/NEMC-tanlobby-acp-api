package ws_api

import "github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"

var server_to_client_listen_packets = map[uint32]bool{
	packet.IDStartGame:        true,
	packet.IDUpdateAttributes: true,
	packet.IDUpdateAbilities:  true,
	packet.IDAddPlayer:        true,
	packet.IDPlayerList:       true,
}

// 由服务端发往客户端的特定数据包是否需要监听
func IsListenedS2CPacket(id uint32) bool {
	_, ok := server_to_client_listen_packets[id]
	return ok
}

// 设置要监听的由服务端发往 Minecraft 客户端的数据包
func SetServerToClientListenPackets(pk_ids []uint32) {
	server_to_client_listen_packets = map[uint32]bool{}
	for _, v := range pk_ids {
		server_to_client_listen_packets[v] = true
	}
}
