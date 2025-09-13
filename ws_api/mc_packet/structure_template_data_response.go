package mc_packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/nbt"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"
)

type StructureTemplateDataResponse struct {
	StructureName     string
	Success           bool
	ResponseType      byte
	StructureTemplate map[string]any
}

func (StructureTemplateDataResponse) ID() uint8 {
	return IDStructureTemplateDataResponse
}

func (StructureTemplateDataResponse) Name() string {
	return "StructureTemplateDataResponse"
}

func (s *StructureTemplateDataResponse) Filter(from packet.Packet, extra map[string]any) {
	pk := from.(*packet.StructureTemplateDataResponse)
	s.StructureName = pk.StructureName
	s.Success = pk.Success
	s.ResponseType = pk.ResponseType
	s.StructureTemplate = pk.StructureTemplate
}

func (*StructureTemplateDataResponse) ToMCPacket() packet.Packet {
	panic("Not implement")
}

func (s *StructureTemplateDataResponse) Marshal(io protocol.IO) {
	io.StringUTF(&s.StructureName)
	io.Bool(&s.Success)
	io.Uint8(&s.ResponseType)
	io.NBT(&s.StructureTemplate, nbt.LittleEndian)
}
