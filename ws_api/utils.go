package ws_api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"
	"github.com/Happy2018new/nemc-tan-lobby-solver/ws_api/mc_packet"
)

// 获取数据包的 ID
func GetPacketID(pk_bytes []byte) (pkID uint32) {
	protocol.Varuint32(bytes.NewReader(pk_bytes), &pkID)
	return
}

func ToFloat64(n any) (float64, bool) {
	switch v := n.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func ToInt32(n any) (int32, bool) {
	switch v := n.(type) {
	case float64:
		return int32(v), true
	case uint8:
		return int32(v), true
	case int8:
		return int32(v), true
	case uint16:
		return int32(v), true
	case int16:
		return int32(v), true
	case int32:
		return v, true
	case uint32:
		return int32(v), true
	case uint64:
		return int32(v), true
	case int64:
		return int32(v), true
	default:
		return 0, false
	}
}

func ToByteArray(n any) ([]byte, bool) {
	switch v := n.(type) {
	case []byte:
		return v, true
	default:
		return nil, false
	}
}

// 将 json array 序列转换为 []uint32 序列。
// 由于 json 数组中的 int 似乎被认为是 float64,
// 因此需要进一步的转换为 uint32
func convertToIntArr(a any) ([]uint32, bool) {
	newArr := []uint32{}
	arr, ok := a.([]any)
	if !ok {
		return nil, false
	}
	for _, e := range arr {
		e1, ok := ToInt32(e)
		if !ok {
			return nil, false
		}
		newArr = append(newArr, uint32(e1))
	}
	return newArr, true
}

func ConstructPacketFromContent(content map[string]any, custom bool, unmarshaler func([]byte, interface{}) error) (pk1 packet.Packet, err error) {
	pkID, ok := ToInt32(content["ID"])
	if !ok {
		err = fmt.Errorf("数据包解析错误: 无效ID: %d", content["ID"])
		return
	}
	if custom {
		defer func() {
			if err := recover(); err != nil {
				err = fmt.Errorf("自定义数据包解析错误: 无效数据包结构: %v", err)
			}
		}()
		bs, ok := ToByteArray(content["Content"])
		if !ok {
			err = fmt.Errorf("自定义数据包解析错误: 无效数据包体: %v", content["Content"])
			return
		}
		pool := mc_packet.NewPool()
		if pkID >= int32(len(pool)) {
			err = fmt.Errorf("自定义数据包解析错误: 无效ID: %d", content["ID"])
		}
		pk := pool[uint32(pkID)]
		pk.Marshal(protocol.NewReader(bytes.NewReader(bs), 0, false))
		pk1 = pk.ToMCPacket()
	} else {
		pk_structure, found := pool[uint32(pkID)]
		if !found {
			err = fmt.Errorf("数据包解析错误: 未知ID: %d", content["ID"])
			return
		}

		pk := pk_structure()
		pk_bts, ok := ToByteArray(content["Content"])
		if !ok {
			err = fmt.Errorf("数据包解析错误: 无效数据包体: %v", content["Content"])
			return

		}
		err = unmarshaler(pk_bts, &pk)
		if err != nil {
			var m map[string]any
			unmarshaler(pk_bts, &m)
			err = fmt.Errorf("数据包解析错误: 无效数据包结构: %v: %#v", err, m)
		}
		pk1 = pk
	}
	return
}

func ConstructContentFromServerPacket(pkID uint32, pk packet.Packet) Message {
	if ok, cid := mc_packet.IsCustomPacket(pkID); ok {
		cp := mc_packet.NewPool()[cid]
		buf := bytes.NewBuffer(nil)
		writer := protocol.NewWriter(buf, 0)
		cp.Marshal(writer)
		return Message{
			Type: WSMSG_SERVER_CUSTOM_PACKET,
			Content: map[string]any{
				"ID":      pkID,
				"Content": buf.Bytes(),
			},
		}
	} else {
		return Message{
			Type: WSMSG_SERVER_PACKET,
			Content: map[string]any{
				"ID":      pkID,
				"Content": pk,
			},
		}
	}
}

func ConstructJsonContentFromServerPacket(pkID uint32, pk packet.Packet) Message {
	if ok, cid := mc_packet.IsCustomPacket(pkID); ok {
		cp := mc_packet.NewPool()[cid]
		buf := bytes.NewBuffer(nil)
		writer := protocol.NewWriter(buf, 0)
		cp.Marshal(writer)
		return Message{
			Type: WSMSG_SERVER_CUSTOM_PACKET,
			Content: map[string]any{
				"ID":      pkID,
				"Content": buf.Bytes(),
			},
		}
	} else {
		bs, _ := json.Marshal(pk)
		return Message{
			Type: WSMSG_SERVER_PACKET_JSON,
			Content: map[string]any{
				"ID":      pkID,
				"Content": string(bs),
			},
		}
	}
}
