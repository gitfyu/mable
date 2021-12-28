package player

import (
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/rs/zerolog/log"
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
	case *play.InKeepAlive:
		p.handleKeepAlive(pk.(*play.InKeepAlive))
	case *play.InPlayer:
		p.handlePlayer(pk.(*play.InPlayer))
	}
}

func (p *Player) handleKeepAlive(pk *play.InKeepAlive) {
	log.Debug().Int("id", pk.ID).Msg("KeepAlive")
}

func (p *Player) handlePlayer(pk *play.InPlayer) {

}
