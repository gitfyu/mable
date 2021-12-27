package packet

import (
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet/handshake"
	"github.com/gitfyu/mable/protocol/packet/login"
	"github.com/gitfyu/mable/protocol/packet/play"
	"github.com/gitfyu/mable/protocol/packet/status"
)

func decodeInbound(s protocol.State, id protocol.VarInt) Inbound {
	switch s {
	case protocol.StateHandshake:
		switch id {
		case 0x00:
			return &handshake.Handshake{}
		}
	case protocol.StateStatus:
		switch id {
		case 0x00:
			return status.Request{}
		case 0x01:
			return &status.Ping{}
		}
	case protocol.StateLogin:
		switch id {
		case 0x00:
			return &login.Start{}
		}
	case protocol.StatePlay:
		switch id {
		case 0x00:
			return &play.InKeepAlive{}
		}
	}

	return nil
}
