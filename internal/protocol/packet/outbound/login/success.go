package login

import (
	"github.com/gitfyu/mable/internal/protocol"
	"github.com/google/uuid"
)

type Success struct {
	UUID     uuid.UUID
	Username string
}

func (Success) PacketID() uint {
	return 0x02
}

func (s *Success) MarshalPacket(w protocol.Writer) error {
	if err := protocol.WriteString(w, s.UUID.String()); err != nil {
		return err
	}
	return protocol.WriteString(w, s.Username)
}
