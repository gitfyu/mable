package login

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/google/uuid"
)

type Success struct {
	UUID     uuid.UUID
	Username string
}

func (s *Success) MarshalPacket(w *protocol.WriteBuffer) {
	w.WriteVarInt(0x02)
	w.WriteString(s.UUID.String())
	w.WriteString(s.Username)
}
