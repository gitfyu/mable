package mable

import (
	"fmt"
	"github.com/gitfyu/mable/protocol/packet"
)

// TODO implement a way to properly generate the JSON response in the future
const defaultResponse = `{"version":{"name":"1.7.6-1.8.9","protocol":47},"players":{"max":0,"online":0},"description":{"text":"Hello world"}}`

func handleStatus(c *connHandler) error {
	if err := readStatusRequest(c); err != nil {
		return err
	}
	if err := writeStatusResponse(c); err != nil {
		return err
	}

	time, err := readStatusPing(c)
	if err != nil {
		return err
	}

	return writeStatusPong(c, time)
}

func readStatusRequest(c *connHandler) error {
	id, _, err := c.readPacket()
	if err != nil {
		return err
	}
	if id != packet.StatusRequest {
		return fmt.Errorf("expected packet %d, got %d", packet.StatusRequest, id)
	}

	return nil
}

func writeStatusResponse(c *connHandler) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteString(defaultResponse)
	return c.WritePacket(packet.StatusResponse, buf)
}

func readStatusPing(c *connHandler) (int64, error) {
	id, buf, err := c.readPacket()
	if err != nil {
		return 0, err
	}
	if id != packet.StatusPing {
		return 0, fmt.Errorf("expected packet %d, got %d", packet.StatusPing, id)
	}

	time, err := buf.ReadLong()
	if err != nil {
		return 0, err
	}

	return time, nil
}

func writeStatusPong(c *connHandler, time int64) error {
	buf := packet.AcquireBuffer()
	defer packet.ReleaseBuffer(buf)

	buf.WriteLong(time)
	return c.WritePacket(packet.StatusPong, buf)
}
