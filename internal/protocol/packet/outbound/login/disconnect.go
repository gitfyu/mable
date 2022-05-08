package login

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/internal/protocol"
)

type Disconnect struct {
	Reason *chat.Msg
}

func (Disconnect) PacketID() uint {
	return 0x00
}

func (d *Disconnect) MarshalPacket(w protocol.Writer) error {
	return protocol.WriteChat(w, d.Reason)
}
