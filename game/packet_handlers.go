package game

import (
	"github.com/gitfyu/mable/internal/protocol/packet"
	inbound "github.com/gitfyu/mable/internal/protocol/packet/inbound/play"
)

func (p *Player) HandlePacket(pk packet.Inbound) {
	p.worldLock.RLock()
	defer p.worldLock.RUnlock()

	p.world.Schedule(func() {
		p.handlePacket(pk)
	})
}

func (p *Player) handlePacket(pk packet.Inbound) {
	switch pk.(type) {
	case *inbound.KeepAlive:
		p.handleKeepAlive(pk.(*inbound.KeepAlive))
	case *inbound.Update:
		p.handlePlayer(pk.(*inbound.Update))
	}
}

func (p *Player) handleKeepAlive(pk *inbound.KeepAlive) {
}

func (p *Player) handlePlayer(pk *inbound.Update) {

}
