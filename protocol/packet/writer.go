package packet

import (
	"github.com/gitfyu/mable/protocol"
	"io"
)

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

func (w *Writer) WritePacket(id ID, data *Buffer) error {
	size := protocol.VarIntSize(protocol.VarInt(id)) + data.Len()
	if err := w.writeVarInt(protocol.VarInt(size)); err != nil {
		return err
	}

	if err := w.writeVarInt(protocol.VarInt(id)); err != nil {
		return err
	}

	if _, err := w.writer.Write(data.Bytes()); err != nil {
		return err
	}

	return nil
}
