package status

import (
	"github.com/gitfyu/mable/protocol"
)

type Response struct {
	Content string
}

func (r *Response) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(0x00)
	w.WriteString(r.Content)
}
