package status

import (
	"github.com/gitfyu/mable/protocol"
)

type Response struct {
	Content string
}

func (_ Response) PacketID() uint {
	return 0x00
}

func (r *Response) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteString(r.Content)
}
