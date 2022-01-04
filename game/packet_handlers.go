package game

import (
	"github.com/gitfyu/mable/internal/protocol/packet"
	inbound "github.com/gitfyu/mable/internal/protocol/packet/inbound/play"
)

// HandlePacket processes a packet sent by the player. This function should only be used by the server itself.
func (p *Player) HandlePacket(pk packet.Inbound) {
	p.packetLock.Lock()
	defer p.packetLock.Unlock()

	p.world.Schedule(func() {
		p.handlePacket(pk)
	})
}

func (p *Player) handlePacket(pk packet.Inbound) {
	switch pk.(type) {
	case *inbound.KeepAlive:
		p.handleKeepAlive(pk.(*inbound.KeepAlive))
	case *inbound.Update:
		p.handleUpdate(pk.(*inbound.Update))
	}
}

func (p *Player) handleKeepAlive(pk *inbound.KeepAlive) {
}

func (p *Player) handleUpdate(pk *inbound.Update) {
	if pk.HasPos {
		oldChunkPos := ChunkPosFromWorldCoords(p.pos.X, p.pos.Z)
		newChunkPos := ChunkPosFromWorldCoords(pk.X, pk.Z)
		p.pos.X, p.pos.Y, p.pos.Z = pk.X, pk.Y, pk.Z

		if oldChunkPos != newChunkPos {
			p.updateChunks()
		}
	}
}
