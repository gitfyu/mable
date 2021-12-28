package mable

import (
	"errors"
	"github.com/gitfyu/mable/protocol/packet/status"
)

// TODO implement a way to properly generate the JSON response in the future
const defaultResponse = `{"version":{"name":"1.7.6-1.8.9","protocol":47},"players":{"max":0,"online":0},"description":{"text":"Hello world"}}`

// handleStatus completes the 'status' flow, which is used by the client to retrieve information to display in the
// server list
func handleStatus(c *conn) error {
	if err := readStatusRequest(c); err != nil {
		return err
	}

	c.WritePacket(&status.Response{
		Content: defaultResponse,
	})

	time, err := readStatusPing(c)
	if err != nil {
		return err
	}

	c.WritePacket(&status.Pong{
		Time: time,
	})
	return nil
}

func readStatusRequest(c *conn) error {
	pk, err := c.readPacket()
	if err != nil {
		return err
	}
	if _, ok := pk.(*status.Request); ok {
		return errors.New("expected status request")
	}

	return nil
}

func readStatusPing(c *conn) (int64, error) {
	pk, err := c.readPacket()
	if err != nil {
		return 0, err
	}

	ping, ok := pk.(*status.Ping)
	if !ok {
		return 0, errors.New("expected ping")
	}

	return ping.Time, nil
}
