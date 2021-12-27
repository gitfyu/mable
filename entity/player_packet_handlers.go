package entity

import (
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/rs/zerolog/log"
)

// HandlePacket processes an incoming packet for the player
func (p *Player) HandlePacket(pk packet.Inbound) error {
	switch pk.(type) {
	case *play.InKeepAlive:
		return p.handleKeepAlive(pk.(*play.InKeepAlive))
	default:
		return nil
	}
}

func (p *Player) handleKeepAlive(pk *play.InKeepAlive) error {
	log.Debug().Int("id", pk.ID).Msg("KeepAlive")
	return nil
}
