package mc_packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"
)

type SubChunkRequest struct {
	Dimension    uint8
	SubChunkPosX int32
	SubChunkPosY int16
	SubChunkPosZ int32
	Offsets      []protocol.SubChunkOffset
}

func (SubChunkRequest) ID() uint8 {
	return IDSubChunkRequest
}

func (SubChunkRequest) Name() string {
	return "SubChunkRequest"
}

func (*SubChunkRequest) Filter(from packet.Packet, extra map[string]any) {
	panic("Not implement")
}

func (s *SubChunkRequest) ToMCPacket() packet.Packet {
	return &packet.SubChunkRequest{
		Dimension: int32(s.Dimension),
		Position: [3]int32{
			s.SubChunkPosX,
			int32(s.SubChunkPosY),
			s.SubChunkPosZ,
		},
		Offsets: s.Offsets,
	}
}

func (s *SubChunkRequest) Marshal(io protocol.IO) {
	io.Uint8(&s.Dimension)
	io.Int32(&s.SubChunkPosX)
	io.Int16(&s.SubChunkPosY)
	io.Int32(&s.SubChunkPosZ)
	protocol.SliceUint16Length(io, &s.Offsets)
}
