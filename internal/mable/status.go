package mable

import (
	"github.com/gitfyu/mable/network"
	"github.com/gitfyu/mable/network/protocol/packet"
)

const defaultResponse = `{"version":{"name":"1.7.6-1.8.9","protocol":47},"players":{"max":0,"online":0},"description":{"text":"Hello world"}}`

var statusHandlers = idToPacketHandler{
	handleStatusRequest,
}

func handleStatusRequest(h *connHandler, _ *network.PacketData) error {
	buf := network.AcquirePacketBuilder()
	defer network.ReleasePacketBuilder(buf)

	buf.Init(packet.StatusResponse).
		PutString(defaultResponse)

	return h.WritePacket(buf)
}
