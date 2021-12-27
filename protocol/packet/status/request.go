package status

import (
	"github.com/gitfyu/mable/protocol"
)

type Request struct{}

func (_ Request) UnmarshalPacket(_ *protocol.ReadBuffer) {
}
