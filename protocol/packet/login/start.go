package login

import (
	"github.com/gitfyu/mable/protocol"
)

type Start struct {
	Username string
}

func (s *Start) UnmarshalPacket(r *protocol.ReadBuffer) {
	s.Username = r.ReadString()
}
