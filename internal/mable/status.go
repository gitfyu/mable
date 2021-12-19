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

func handleStatusRequest(h *connHandler, _ *packet.Buffer) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteString(defaultResponse)
	return h.WritePacket(packet.StatusResponse, buf)
}

func handlePing(h *connHandler, data *packet.Buffer) error {
	time, err := data.ReadLong()
	if err != nil {
		return err
	}

	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteLong(time)
	if err := h.WritePacket(packet.StatusPong, buf); err != nil {
		return err
	}

	return h.Close()
}
