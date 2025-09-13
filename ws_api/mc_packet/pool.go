package mc_packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"
)

const (
	IDLevelChunk uint8 = iota
	IDSubChunk
	IDSubChunkRequest
	IDStructureTemplateDataResponse
)

var MCPacketIDToCustomPacketID map[uint32]uint8 = map[uint32]uint8{
	packet.IDLevelChunk:                    IDLevelChunk,
	packet.IDSubChunk:                      IDSubChunk,
	packet.IDSubChunkRequest:               IDSubChunkRequest,
	packet.IDStructureTemplateDataResponse: IDStructureTemplateDataResponse,
}

type MCPacket interface {
	ID() uint8
	Name() string
	Filter(from packet.Packet, extra map[string]any)
	ToMCPacket() packet.Packet
	Marshal(io protocol.IO)
}

func NewPool() []MCPacket {
	return []MCPacket{
		&LevelChunk{},
		&SubChunk{},
		&SubChunkRequest{},
		&StructureTemplateDataResponse{},
	}
}

func IsCustomPacket(pk_id uint32) (ok bool, cid uint8) {
	cid, ok = MCPacketIDToCustomPacketID[pk_id]
	return
}
