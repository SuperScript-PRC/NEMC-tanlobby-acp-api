package mc_packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"
)

type SubChunkEntry struct {
	Result       byte
	SubChunkPosX int32
	SubChunkPosY int16
	SubChunkPosZ int32
	NBTData      []byte
	BlobHash     uint64
}

func (s *SubChunkEntry) Marshal(io protocol.IO) {
	io.Uint8(&s.Result)
	io.Int32(&s.SubChunkPosX)
	io.Int16(&s.SubChunkPosY)
	io.Int32(&s.SubChunkPosZ)
	protocol.FuncSliceUint32Length(io, &s.NBTData, io.Uint8)
	io.Uint64(&s.BlobHash)
}

type SubChunk struct {
	Dimension    byte
	Entries      []SubChunkEntry
	CacheEnabled uint8
}

func (SubChunk) ID() uint8 {
	return IDSubChunk
}

func (SubChunk) Name() string {
	return "SubChunk"
}

func (s *SubChunk) Filter(from packet.Packet, extra map[string]any) {
	pk := from.(*packet.SubChunk)

	s.Dimension = byte(pk.Dimension)
	s.Entries = make([]SubChunkEntry, len(pk.SubChunkEntries))

	for index, value := range pk.SubChunkEntries {
		s.Entries[index] = SubChunkEntry{
			Result:       value.Result,
			SubChunkPosX: pk.Position[0] + int32(value.Offset[0]),
			SubChunkPosY: int16(pk.Position[1] + int32(value.Offset[1])),
			SubChunkPosZ: pk.Position[2] + int32(value.Offset[2]),
			NBTData:      value.RawPayload,
			BlobHash:     value.BlobHash,
		}
	}

	if pk.CacheEnabled {
		s.CacheEnabled = 1
	} else {
		s.CacheEnabled = 0
	}
}

func (*SubChunk) ToMCPacket() packet.Packet {
	panic("Not implement")
}

func (s *SubChunk) Marshal(io protocol.IO) {
	io.Uint8(&s.Dimension)
	protocol.SliceUint16Length(io, &s.Entries)
	io.Uint8(&s.CacheEnabled)
}
