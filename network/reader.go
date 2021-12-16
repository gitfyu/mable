package network

import (
	"bufio"
	"errors"
	"github.com/gitfyu/mable/network/protocol"
	"github.com/gitfyu/mable/network/protocol/packet"
	"io"
)

var errPacketTooLarge = errors.New("packet exceeds maximum size")

// ReaderConfig is used to configure settings for a Reader
type ReaderConfig struct {
	// MaxPacketSize is the maximum size in bytes per packet. Larger packets will be rejected
	MaxPacketSize int
}

// Reader is used to read packets one at a time
type Reader struct {
	reader *bufio.Reader
	cfg    ReaderConfig
}

// NewReader constructs a new Reader that reads from the provided io.Reader
func NewReader(r io.Reader, cfg ReaderConfig) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
		cfg:    cfg,
	}
}

// ReadPacket reads a single packet. The buf parameter will be used to store the packet if it is large enough, otherwise
// a new buffer will be allocated. This function returns a buffer, which might be the buffer passed to this function,
// holding the packets contents.
func (r *Reader) ReadPacket(buf []byte) (packet.ID, []byte, error) {
	var size protocol.VarInt
	if err := protocol.ReadVarInt(r.reader, &size); err != nil {
		return 0, nil, err
	}
	if int(size) > r.cfg.MaxPacketSize {
		return 0, nil, errPacketTooLarge
	}

	var id protocol.VarInt
	if err := protocol.ReadVarInt(r.reader, &id); err != nil {
		return 0, nil, err
	}

	dataSize := int(size) - protocol.VarIntSize(id)
	if cap(buf) >= dataSize {
		buf = buf[:dataSize]
	} else {
		buf = make([]byte, dataSize)
	}

	if _, err := io.ReadFull(r.reader, buf); err != nil {
		return 0, nil, err
	}

	return packet.ID(id), buf, nil
}
