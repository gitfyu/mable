package mable

import (
	chat2 "github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol/packet"
)

var loginHandlers = newPacketHandlerLookup(
	packetHandlers{
		packet.LoginStart: handleLoginStart,
	},
)

func handleLoginStart(h *connHandler, data *network.PacketData) error {
	return h.Disconnect(&chat2.Msg{Text: "TODO", Color: chat2.ColorYellow})
}
