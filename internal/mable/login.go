package mable

import (
	"errors"
	"github.com/gitfyu/mable/protocol/packet/login"
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
	return username, id, writeLoginSuccess(c, username, id)
}

func readLoginStart(c *conn) (string, error) {
	pk, err := c.readPacket()
	if err != nil {
		return "", err
	}
	l, ok := pk.(*login.Start)
	if !ok {
		return "", errors.New("expected login start")
	}

	return l.Username, nil
}

func writeLoginSuccess(c *conn, username string, id uuid.UUID) error {
	pk := login.Success{
		UUID:     id,
		Username: username,
	}
	return c.WritePacket(&pk)
}
