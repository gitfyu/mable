package mable

import (
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol/packet"
)

var loginHandlers = newPacketHandlerLookup(
	packetHandlers{
		packet.LoginStart: handleLoginStart,
	},
)

func handleLoginStart(h *connHandler, data *network.PacketData) error {
	return h.Disconnect(`{"text":"TODO","color":"yellow"}`)
}
