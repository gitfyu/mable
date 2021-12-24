package entity

import (
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/rs/zerolog/log"
)

// HandlePacket processes an incoming packet for the player
func (p *Player) HandlePacket(id packet.ID, data *packet.Buffer) error {
	switch id {
	case packet.PlayClientKeepAlive:
		return p.handleKeepAlive(data)
	case packet.PlayClientPos:
		return p.handleFlying(true, false, data)
	case packet.PlayClientPosAndLook:
		return p.handleFlying(true, true, data)
	default:
		return nil
	}
}

func (p *Player) handleKeepAlive(data *packet.Buffer) error {
	i, err := data.ReadVarInt()
	if err != nil {
		return err
	}

	log.Debug().Int("id", int(i)).Msg("KeepAlive")
	return nil
}

// handleFlying handles one of several 'flying' packets
func (p *Player) handleFlying(hasPos bool, hasLook bool, data *packet.Buffer) error {
	if hasPos {
		var coords [3]float64

		for i := 0; i < 3; i++ {
			var err error
			coords[i], err = data.ReadDouble()
			if err != nil {
				return err
			}
		}

		p.setCoords(coords[0], coords[1], coords[2])
	}

	// TODO process packet further

	return nil
}
