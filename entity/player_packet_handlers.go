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
		return p.handlePos(data)
	case packet.PlayClientPosAndLook:
		return p.handlePosAndLook(data)
	default:
		return nil
	}
}

// readFlyingCoords reads the coordinates from a flying packet
func readFlyingCoords(data *packet.Buffer) (float64, float64, float64, error) {
	var coords [3]float64

	for i := 0; i < 3; i++ {
		var err error
		coords[i], err = data.ReadDouble()
		if err != nil {
			return 0, 0, 0, err
		}
	}

	return coords[0], coords[1], coords[2], nil
}

func (p *Player) handleKeepAlive(data *packet.Buffer) error {
	i, err := data.ReadVarInt()
	if err != nil {
		return err
	}

	log.Debug().Int("id", int(i)).Msg("KeepAlive")
	return nil
}

// handlePos handles a packet.PlayClientPos
func (p *Player) handlePos(data *packet.Buffer) error {
	x, y, z, err := readFlyingCoords(data)
	if err != nil {
		return err
	}

	// TODO process packet further

	p.setCoords(x, y, z)
	return nil
}

// handlePosAndLook handles a packet.PlayClientPosAndLook
func (p *Player) handlePosAndLook(data *packet.Buffer) error {
	x, y, z, err := readFlyingCoords(data)
	if err != nil {
		return err
	}

	// TODO process packet further

	p.setCoords(x, y, z)
	return nil
}
