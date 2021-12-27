package packet

import (
	"bufio"
	"errors"
	"github.com/gitfyu/mable/protocol"
	"io"
)

var (
	errTooLarge  = errors.New("packet exceeds maximum size")
	errBadPacket = errors.New("bad packet")
)

// ReaderConfig is used to configure settings for a Reader
type ReaderConfig struct {
	// MaxSize is the maximum size in Bytes per packet. Larger packets will be rejected.
	MaxSize int
}

// Reader is used to read packets
type Reader struct {
	reader  *bufio.Reader
	cfg     ReaderConfig
	readBuf protocol.ReadBuffer
}

// NewReader constructs a new Reader that reads from the provided io.Reader
func NewReader(r io.Reader, cfg ReaderConfig) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
		cfg:    cfg,
	}
}

// ReadPacket reads a single packet. It returns nil for unknown packets.
func (r *Reader) ReadPacket(state protocol.State) (pk Inbound, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errBadPacket
			}
		}
	}()

	var size protocol.VarInt
	if err := protocol.ReadVarInt(r.reader, &size); err != nil {
		return nil, err
	}
	if int(size) > r.cfg.MaxSize {
		return nil, errTooLarge
	}

	var id protocol.VarInt
	if err := protocol.ReadVarInt(r.reader, &id); err != nil {
		return nil, err
	}

	if err := r.readBuf.ReadAll(r.reader, int(size)-protocol.VarIntSize(id)); err != nil {
		return nil, err
	}

	pk = decodeInbound(state, id)
	if pk == nil {
		return nil, nil
	}

	pk.UnmarshalPacket(&r.readBuf)
	return pk, nil
}
