package login

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/internal/protocol"
)

type Disconnect struct {
	Reason *chat.Msg
}

func (_ Disconnect) PacketID() uint {
	return 0x00
}

func (s *Disconnect) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteChat(s.Reason)
}