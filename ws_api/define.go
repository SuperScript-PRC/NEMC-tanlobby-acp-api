package ws_api

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
)

// WebSocket 通信部分

// 基本消息类型
const (
	WSMSG_SERVER_PACKET             = "ServerMCPacket"         // [接入点->WSCli] 来自租赁服服务端的数据包
	WSMSG_SERVER_PACKET_JSON        = "ServerMCPacketJson"     // [接入点->WSCli] 来自租赁服服务端的数据包Json
	WSMSG_SERVER_CUSTOM_PACKET      = "ServerCustomMCPacket"   // [接入点->WSCli] 租赁服->客户端的自定义字节数据包
	WSMSG_SET_BOT_BASIC_INFO        = "SetBotBasicInfo"        // [接入点->WSCli] 设置玩家基本信息 (RuntimeID 等)
	WSMSG_SET_SERVER_LISTEN_PACKETS = "SetServerListenPackets" // [WSCli->接入点] 设置需要监听的来自服务器的数据包
	WSMSG_UPDATE_UQ                 = "UpdateUQ"               // [LanGame->WSCli] 更新客户端玩家数据
	WSMSG_UPDATE_ABILITIES          = "UpdateAbilities"        // [接入点->WSCli] 更新客户端玩家能力数据
	WSMSG_GAME_READY                = "GameReady"              // [接入点->WSCli] 就绪
)

// 额外的特殊消息类型
const (
	WSMSG_BreakBlock = "BreakBlock" // [WSCli->LanGame] 请求挖掘方块
)

type Message struct {
	// 消息类型
	Type string `msgpack:"type"`
	// 消息正文
	Content any `msgpack:"content"`
}

// UQ 部分

// LanGame所操控的玩家的基本信息。
type BotBasicData struct {
	Name           string `msgpack:"bot_name"`
	UUID           string `msgpack:"uuid"`
	EntityUniqueID int64  `msgpack:"bot_entity_unique_id"`
	RuntimeID      uint64 `msgpack:"bot_runtime_id"`
}

// LanGame所在租赁服的玩家的基本信息。
type PlayerBasicInfo struct {
	Name      string `msgpack:"name"`
	UUID      string `msgpack:"uuid"`
	XUID      string `msgpack:"xuid"`
	UniqueID  int64  `msgpack:"uniqueID"`
	Abilities any    `msgpack:"abilities"`
}

type PlayersUQMap map[string]*PlayerBasicInfo

func (pb *PlayerBasicInfo) SetAbilities(abilities protocol.AbilityData) {
	pb.Abilities = abilities
}
