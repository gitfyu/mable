package status

import (
	"github.com/gitfyu/mable/internal/protocol"
)

type Response struct {
	Content string
}

func (Response) PacketID() uint {
	return 0x00
}

func (r *Response) MarshalPacket(w protocol.Writer) error {
	return protocol.WriteString(w, r.Content)
}
