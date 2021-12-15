package mable

const defaultResponse = `{"version":{"name":"1.7.6-1.8.9","protocol":47},"players":{"max":0,"online":0},"description":{"text":"Hello world"}}`

var statusHandlers = idToPacketHandler{
	handleStatusRequest,
}

func handleStatusRequest(_ int, h *connHandler) error {
	ok := h.enc.WriteVarInt(0x00) &&
		h.enc.WriteString(defaultResponse) &&
		h.enc.WritePacket(true)
	if !ok {
		return h.enc.LastError()
	}

	return nil
}
