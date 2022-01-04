package packet

import (
	"fmt"
	"github.com/gitfyu/mable/internal/protocol"
	"io"
	"sync"
)

var writeBufPool = sync.Pool{
	New: func() interface{} {
		return new(protocol.WriteBuffer)
	},
}

// Writer is used to write packets.
type Writer struct {
	writer io.Writer
}

// NewWriter constructs a new Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) writeVarInt(v int32) error {
	n := protocol.VarIntSize(v)
	b := make([]byte, n)
	protocol.WriteVarInt(b, v)

	if _, err := w.writer.Write(b); err != nil {
		return err
	}

	return nil
}

// WritePacket writes a single packet, including its length and id.
func (w *Writer) WritePacket(pk Outbound) error {
	buf := writeBufPool.Get().(*protocol.WriteBuffer)
	defer writeBufPool.Put(buf)

	buf.Reset()
	buf.WriteVarInt(int32(pk.PacketID()))
	pk.MarshalPacket(buf)

	if err := w.writeVarInt(int32(buf.Len())); err != nil {
		return fmt.Errorf("writing packet size: %w", err)
	}

	if _, err := w.writer.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("writing packet body: %w", err)
	}

	return nil
}
