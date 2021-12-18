package mable

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol/packet"
)

var loginHandlers = newPacketHandlerLookup(
	packetHandlers{
		packet.LoginStart: handleLoginStart,
	},
)

func handleLoginStart(h *connHandler, data *network.PacketData) error {
	reason := chat.NewBuilder("TODO: ").
		Bold().
		Color(chat.ColorGold).
		Append("not yet implemented.").
		NotBold().
		Color(chat.ColorYellow).
		Build()

	return h.Disconnect(reason)
}
