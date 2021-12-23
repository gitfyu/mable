package mable

import (
	"errors"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/google/uuid"
)

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
	id, buf, err := c.readPacket()
	if err != nil {
		return "", err
	}
	if id != packet.LoginStart {
		return "", errors.New("expected login start")
	}

	return buf.ReadString()
}

func writeLoginSuccess(c *conn, username string, id uuid.UUID) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteString(id.String())
	buf.WriteString(username)

	return c.WritePacket(packet.LoginSuccess, buf)
}
