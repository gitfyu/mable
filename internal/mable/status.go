package mable

import (
	"github.com/gitfyu/mable/protocol/packet"
)

// TODO implement a way to properly generate the JSON response in the future
const defaultResponse = `{"version":{"name":"1.7.6-1.8.9","protocol":47},"players":{"max":0,"online":0},"description":{"text":"Hello world"}}`

var statusHandlers = newPacketHandlerLookup(
	packetHandlers{
		packet.StatusRequest: handleStatusRequest,
		packet.StatusPing:    handlePing,
	},
)

func handleStatusRequest(h *connHandler, _ *packet.Packet) error {
	buf := packet.AcquireBuilder()
	defer packet.ReleaseBuilder(buf)

	buf.Init(packet.StatusResponse).
		PutString(defaultResponse)

	return h.WritePacket(buf)
}

func handlePing(h *connHandler, p *packet.Packet) error {
	time := p.GetLong()

	buf := packet.AcquireBuilder()
	defer packet.ReleaseBuilder(buf)

	buf.Init(packet.StatusPong).
		PutLong(time)

	if err := h.WritePacket(buf); err != nil {
		return err
	}

	return h.Close()
}
