package mable

import "github.com/rs/zerolog/log"

var statusHandlers = idToPacketHandler{
	handleStatusRequest,
}

func handleStatusRequest(_ int, _ *connHandler) error {
	log.Debug().Msg("status request")
	return nil
}
