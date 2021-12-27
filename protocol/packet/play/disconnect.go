package play

import (
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
)

type OutDisconnect struct {
	Reason *chat.Msg
}

func (_ OutDisconnect) PacketID() uint {
	return 0x40
}

func (s *OutDisconnect) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteChat(s.Reason)
}
