package packet

import (
	"bufio"
	"errors"
	"github.com/gitfyu/mable/protocol"
	"io"
)

var errTooLarge = errors.New("packet exceeds maximum size")

// ReaderConfig is used to configure settings for a Reader
type ReaderConfig struct {
	// MaxSize is the maximum size in bytes per packet. Larger packets will be rejected.
	MaxSize int
}

// Reader is used to read packets
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

// ReadPacket reads a single packet
func (r *Reader) ReadPacket(buf *Buffer) (ID, error) {
	var size protocol.VarInt
	if err := protocol.ReadVarInt(r.reader, &size); err != nil {
		return 0, err
	}
	if int(size) > r.cfg.MaxSize {
		return 0, errTooLarge
	}

	data := make([]byte, int(size))
	if _, err := io.ReadFull(r.reader, data); err != nil {
		return 0, err
	}

	_, _ = buf.Write(data)

	id, err := buf.ReadVarInt()
	if err != nil {
		return 0, err
	}

	return ID(id), nil
}
