package mable

import (
	"errors"
	inbound "github.com/gitfyu/mable/protocol/packet/inbound/login"
	outbound "github.com/gitfyu/mable/protocol/packet/outbound/login"
	"github.com/google/uuid"
)

// handleLogin processes the login sequence. Currently, it only supports offline ('cracked') mode.
func handleLogin(c *conn) (string, uuid.UUID, error) {
	username, err := readLoginStart(c)
	if err != nil {
		return "", uuid.Nil, err
	}

	// TODO implement authenticated login

	id := generateOfflineUUID(username)
	c.WritePacket(&outbound.Success{
		UUID:     id,
		Username: username,
	})
	return username, id, nil
}

func readLoginStart(c *conn) (string, error) {
	pk, err := c.readPacket()
	if err != nil {
		return "", err
	}
	l, ok := pk.(*inbound.Start)
	if !ok {
		return "", errors.New("expected login start")
	}

	return l.Username, nil
}
