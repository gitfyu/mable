package packet

import (
	"github.com/gitfyu/mable/protocol"
)

type Outbound interface {
	PacketID() uint
	MarshalPacket(w *protocol.WriteBuffer)
}
