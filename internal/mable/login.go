package mable

import (
	"errors"
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/entity"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/google/uuid"
)

func handleLogin(c *conn) error {
	if c.version != protocol.Version_1_7_6 && c.version != protocol.Version_1_8 {
		return cancelLogin(c, "Please use Minecraft 1.7.6-1.8.9!")
	}

	username, err := readLoginStart(c)
	if err != nil {
		return err
	}

	// TODO implement authenticated login

	id := generateOfflineUUID(username)
	c.player = &player{
		name: username,
		uid:  id,
		id:   entity.GenId(),
	}
	return writeLoginSuccess(c, username, id)
}

func cancelLogin(c *conn, reason string) error {
	msg := chat.Msg{
		Text:  reason,
		Color: chat.ColorRed,
	}
	return c.Disconnect(&msg)
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
