package packet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/gitfyu/mable/internal/protocol"
)

var (
	errTooLarge = errors.New("packet exceeds maximum size")
)

// ReaderConfig is used to configure settings for a Reader.
type ReaderConfig struct {
	// MaxSize is the maximum size in Bytes per packet. Larger packets will be rejected.
	MaxSize int
}

// Reader is used to read packets.
type Reader struct {
	reader  *bufio.Reader
	cfg     ReaderConfig
	readBuf bytes.Buffer
}

// NewReader constructs a new Reader that reads from the provided io.Reader.
func NewReader(r io.Reader, cfg ReaderConfig) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
		cfg:    cfg,
	}
}

// ReadPacket reads a single packet. It returns nil for unknown packets.
func (r *Reader) ReadPacket(state protocol.State) (pk Inbound, err error) {
	size, err := protocol.ReadVarInt(r.reader)
	if err != nil {
		return nil, fmt.Errorf("reading packet size: %w", err)
	}
	if int(size) > r.cfg.MaxSize {
		return nil, errTooLarge
	}

	id, err := protocol.ReadVarInt(r.reader)
	if err != nil {
		return nil, fmt.Errorf("reading packet ID: %w", err)
	}
	if err := r.readToBuf(int64(size) - int64(protocol.VarIntSize(id))); err != nil {
		return nil, fmt.Errorf("reading packet body: %w", err)
	}

	pk = createInbound(state, uint(id))
	if pk == nil {
		return nil, nil
	}
	if err := pk.UnmarshalPacket(&r.readBuf); err != nil {
		return nil, fmt.Errorf("bad packet 0x%x: %w", id, err)
	}
	return pk, nil
}

// readToBuf reads exactly n bytes to the internal buffer.
func (r *Reader) readToBuf(n int64) error {
	r.readBuf.Reset()

	lr := io.LimitedReader{
		R: r.reader,
		N: n,
	}
	_, err := r.readBuf.ReadFrom(&lr)
	return err
}
