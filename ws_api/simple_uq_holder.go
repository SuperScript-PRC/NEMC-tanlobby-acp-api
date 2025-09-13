package ws_api

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"

	"github.com/google/uuid"
)

// 维护一个简单的UqHolder, 储存全局玩家的基本信息

var simple_uq_map = make(PlayersUQMap)

// 机器人本身的信息
var botDatasReady bool = false
var botName string
var botUniqueID int64
var botUUID uuid.UUID
var botRuntimeID uint64

// 处理 PlayerList 以更新 UQ 表
func handlePlayerList(pk packet.PlayerList) {
	if pk.ActionType == 0 {
		for _, entry := range pk.Entries {
			addPlayerUQ(entry)
		}
	} else {
		for _, entry := range pk.Entries {
			removePlayerUQ(entry)
		}
	}
	BroadcastMessageToEndpoints(Message{
		Type:    WSMSG_UPDATE_UQ,
		Content: simple_uq_map,
	})
}

func GetUQMap() PlayersUQMap {
	return simple_uq_map
}

// 设置 (更新) UQ 表中的玩家实例
func SetUQMapPlayer(uqmap PlayersUQMap, playerobj *PlayerBasicInfo) {
	uqmap[playerobj.Name] = playerobj
}

// 从玩家的 Unique ID 获取玩家名
func GetPlayerNameByUniqueID(uqid int64) (string, bool) {
	for k, v := range GetUQMap() {
		if v.UniqueID == uqid {
			return k, true
		}
	}
	return "", false
}

func getBotBasicInfo() BotBasicData {
	return BotBasicData{
		Name:           botName,
		UUID:           botUUID.String(),
		EntityUniqueID: botUniqueID,
		RuntimeID:      botRuntimeID,
	}
}

// 向 UQ 表添加玩家对象
func addPlayerUQ(entry protocol.PlayerListEntry) {
	simple_uq_map[entry.Username] = &PlayerBasicInfo{
		Name:     entry.Username,
		UUID:     entry.UUID.String(),
		UniqueID: entry.EntityUniqueID,
		XUID:     entry.XUID,
	}
}

// 从 UQ 表移除玩家对象
func removePlayerUQ(entry protocol.PlayerListEntry) {
	delete(simple_uq_map, entry.Username)
}
