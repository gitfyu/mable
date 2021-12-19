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

// ReadPacket reads a single packet. If a nil error is returned, then the returned Buffer should be released using
// ReleaseBuffer.
func (r *Reader) ReadPacket() (ID, *Buffer, error) {
	var size protocol.VarInt
	if err := protocol.ReadVarInt(r.reader, &size); err != nil {
		return 0, nil, err
	}
	if int(size) > r.cfg.MaxSize {
		return 0, nil, errTooLarge
	}

	data := make([]byte, int(size))
	if _, err := io.ReadFull(r.reader, data); err != nil {
		return 0, nil, err
	}

	buf := AcquireBuffer()
	_, _ = buf.Write(data)

	id, err := buf.ReadVarInt()
	if err != nil {
		return 0, nil, err
	}

	return ID(id), buf, nil
}
