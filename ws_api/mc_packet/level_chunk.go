package mc_packet

import (
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft/protocol/packet"
)

type LevelChunk struct {
	Dimension            uint8
	ChunkPosX            int32
	ChunkPosZ            int32
	HighestSubChunkIndex uint8
	CacheEnabled         byte
}

func (LevelChunk) ID() uint8 {
	return IDLevelChunk
}

func (LevelChunk) Name() string {
	return "LevelChunk"
}

func (l *LevelChunk) Filter(from packet.Packet, extra map[string]any) {
	pk := from.(*packet.LevelChunk)
	l.Dimension = uint8(pk.Dimension)
	l.ChunkPosX = pk.Position[0]
	l.ChunkPosZ = pk.Position[1]
	l.HighestSubChunkIndex = uint8(pk.HighestSubChunk)
	if pk.CacheEnabled {
		l.CacheEnabled = 1
	} else {
		l.CacheEnabled = 0
	}
}

func (*LevelChunk) ToMCPacket() packet.Packet {
	panic("Not implement")
}

func (l *LevelChunk) Marshal(io protocol.IO) {
	io.Uint8(&l.Dimension)
	io.Int32(&l.ChunkPosX)
	io.Int32(&l.ChunkPosZ)
	io.Uint8(&l.HighestSubChunkIndex)
	io.Uint8(&l.CacheEnabled)
}
