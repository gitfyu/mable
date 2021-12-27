package entity

import (
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/rs/zerolog/log"
)

// HandlePacket processes an incoming packet for the player
func (p *Player) HandlePacket(pk packet.Inbound) {
	switch pk.(type) {
	case *play.InKeepAlive:
		p.handleKeepAlive(pk.(*play.InKeepAlive))
	}
}

func (p *Player) handleKeepAlive(pk *play.InKeepAlive) {
	log.Debug().Int("id", pk.ID).Msg("KeepAlive")
}
