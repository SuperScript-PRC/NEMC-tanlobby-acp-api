package ws_api

// 处理由LanGame发送到 WebSocket Client 的所有消息

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"

	"github.com/pterm/pterm"
)

var pool = packet.ListAllPackets()

// 处理第一次获取到的 StartGame 数据包
// 以获取到当前LanGame所使用的用户的 EntityRuntimeID 以及 EntityUniqueID
func handleStartGame(pk packet.StartGame) {
	botRuntimeID = pk.EntityRuntimeID
	botUniqueID = pk.EntityUniqueID
}

// 处理第一次获取到的 PlayerList 数据包
// 现在已不起实际作用。
func handleFirstPlayerList(conn *minecraft.Conn, pk packet.PlayerList) {
	for _, entry := range pk.Entries {
		if entry.EntityUniqueID == conn.GameData().EntityUniqueID {
			botRuntimeID = conn.GameData().EntityRuntimeID
			botUniqueID = entry.EntityUniqueID
			botName = entry.Username
			botUUID = entry.UUID
			botDatasReady = true
			handoutBotBasicInfo()
		}
	}
}

// 处理 UpdateAbilities 数据包
// 以更新玩家能力
func handleAbilitySet(pk protocol.AbilityData) {
	playername, found := GetPlayerNameByUniqueID(pk.EntityUniqueID)
	if found {
		uqmap := GetUQMap()
		playerobj := uqmap[playername]
		playerobj.SetAbilities(pk)
		SetUQMapPlayer(uqmap, playerobj)
	} else {
		pterm.Error.Println("未找到", pk.EntityUniqueID, "所对应的玩家")
		return
	}
	BroadcastMessageToEndpoints(Message{
		Type:    WSMSG_UPDATE_UQ,
		Content: simple_uq_map,
	})
}

// 过滤掉被拦截的数据包
func FilterBlockingPackets(
	pks []packet.Packet,
	filterFunc func(uint32) bool,
) []packet.Packet {
	new_pks := []packet.Packet{}
	for _, pk := range pks {
		if !filterFunc(pk.ID()) {
			new_pks = append(new_pks, pk)
		}
	}
	return new_pks
}

// 判断 租赁服 -> Minecraft客户端 的数据包是否应该被 ws_api 处理
// 并转发至 WebSocket 客户端
func HandleServerPacketsToEndpoints(conn *minecraft.Conn, pk packet.Packet) {
	if !IsListenedS2CPacket(pk.ID()) {
		return
	}
	switch pk1 := pk.(type) {
	case *packet.StartGame:
		if !botDatasReady {
			handleStartGame(*pk1)
		}
	case *packet.PlayerList:
		handlePlayerList(*pk1)
		if !botDatasReady {
			handleFirstPlayerList(conn, *pk1)
		}
	case *packet.UpdateAbilities:
		handleAbilitySet(pk1.AbilityData)
	case *packet.AddPlayer:
		handleAbilitySet(pk1.AbilityData)
	}
	BroadcastMessageToEndpoints(ConstructJsonContentFromServerPacket(pk.ID(), pk))

}

func SetReady() {
	BroadcastMessageToEndpoints(Message{
		Type:    WSMSG_GAME_READY,
		Content: nil,
	})
}
