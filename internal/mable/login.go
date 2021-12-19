package mable

import (
	"errors"
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/gitfyu/mable/protocol/packet"
	"github.com/google/uuid"
)

func handleLogin(c *connHandler) error {
	if c.version != protocol.Version_1_7_6 && c.version != protocol.Version_1_8 {
		return cancelLogin(c, "Please use Minecraft 1.7.6-1.8.9!")
	}

	username, err := readLoginStart(c)
	if err != nil {
		return err
	}

	// TODO implement authenticated login

	id := generateOfflineUUID(username)
	return writeLoginSuccess(c, username, id)
}

func cancelLogin(c *connHandler, reason string) error {
	msg := chat.Msg{
		Text:  reason,
		Color: chat.ColorRed,
	}
	return c.Disconnect(&msg)
}

func readLoginStart(c *connHandler) (string, error) {
	id, buf, err := c.readPacket()
	if err != nil {
		return "", err
	}
	if id != packet.LoginStart {
		return "", errors.New("expected login start")
	}

	return buf.ReadString()
}

func writeLoginSuccess(c *connHandler, username string, id uuid.UUID) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteString(id.String())
	buf.WriteString(username)

	return c.WritePacket(packet.LoginSuccess, buf)
}
