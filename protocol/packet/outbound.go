package packet

import (
	"github.com/gitfyu/mable/protocol"
)

type Outbound interface {
	MarshalPacket(w *protocol.WriteBuffer)
}
