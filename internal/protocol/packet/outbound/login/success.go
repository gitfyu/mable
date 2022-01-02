package login

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/google/uuid"
)

type Success struct {
	UUID     uuid.UUID
	Username string
}

func (_ Success) PacketID() uint {
	return 0x02
}

func (s *Success) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteString(s.UUID.String())
	w.WriteString(s.Username)
}
