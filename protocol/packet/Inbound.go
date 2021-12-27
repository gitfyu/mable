package packet

import (
	"github.com/gitfyu/mable/protocol"
)

type Inbound interface {
	UnmarshalPacket(r *protocol.ReadBuffer)
}
