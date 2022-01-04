package server

import (
	"errors"
	inbound "github.com/gitfyu/mable/internal/protocol/packet/inbound/status"
	outbound "github.com/gitfyu/mable/internal/protocol/packet/outbound/status"
)

// TODO implement a way to properly generate the JSON response in the future
const defaultResponse = `{"version":{"name":"1.7.6-1.8.9","protocol":47},"players":{"max":0,"online":0},"description":{"text":"Hello world"}}`

// handleStatus processes the status flow.
func handleStatus(c *conn) error {
	if err := readStatusRequest(c); err != nil {
		return err
	}

	c.WritePacket(&outbound.Response{
		Content: defaultResponse,
	})

	time, err := readStatusPing(c)
	if err != nil {
		return err
	}

	c.WritePacket(&outbound.Pong{
		Time: time,
	})
	return nil
}

func readStatusRequest(c *conn) error {
	pk, err := c.readPacket()
	if err != nil {
		return err
	}
	if _, ok := pk.(*inbound.Request); ok {
		return errors.New("expected status request")
	}

	return nil
}

func readStatusPing(c *conn) (int64, error) {
	pk, err := c.readPacket()
	if err != nil {
		return 0, err
	}

	ping, ok := pk.(*inbound.Ping)
	if !ok {
		return 0, errors.New("expected ping")
	}

	return ping.Time, nil
}
