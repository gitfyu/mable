package packet

import (
	"github.com/gitfyu/mable/protocol"
	"io"
	"sync"
)

var writeBufPool = sync.Pool{
	New: func() interface{} {
		return new(protocol.WriteBuffer)
	},
}

// Writer is used to write packets
type Writer struct {
	writer io.Writer
}

// NewWriter constructs a new Writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) writeVarInt(v protocol.VarInt) error {
	n := protocol.VarIntSize(v)
	b := make([]byte, n)
	protocol.WriteVarInt(b, v)

	if _, err := w.writer.Write(b); err != nil {
		return err
	}

	return nil
}

func (w *Writer) WritePacket(pk Outbound) error {
	buf := writeBufPool.Get().(*protocol.WriteBuffer)
	defer writeBufPool.Put(buf)

	buf.Reset()
	pk.MarshalPacket(buf)

	if err := w.writeVarInt(protocol.VarInt(buf.Len())); err != nil {
		return err
	}

	if _, err := w.writer.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
