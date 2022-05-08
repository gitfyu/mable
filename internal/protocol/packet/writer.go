package packet

import (
	"bytes"
	"io"

	"github.com/gitfyu/mable/internal/protocol"
)

// Writer is used to write packets.
type Writer struct {
	out io.Writer
	// buf is a buffer holding packets that are ready to be flushed.
	buf bytes.Buffer
	// dataBuf is a buffer used to store encoded packet data,
	// cached for performance.
	dataBuf bytes.Buffer
}

// NewWriter constructs a new Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		out: w,
	}
}

// WritePacket adds a single packet to the internal buffer, which will be written the next
// time that Flush is called.
func (w *Writer) WritePacket(pk Outbound) error {
	w.dataBuf.Reset()

	// 1. Encode the packet ID + content
	protocol.WriteVarInt(&w.dataBuf, int32(pk.PacketID()))
	if err := pk.MarshalPacket(&w.dataBuf); err != nil {
		return err
	}

	// 2. Write the length of the encoded data, followed by the data itself
	protocol.WriteVarInt(&w.buf, int32(w.dataBuf.Len()))
	w.buf.Write(w.dataBuf.Bytes())
	return nil
}

// Writes the internal buffer to the io.Writer that was used to construct this Writer.
func (w *Writer) Flush() error {
	_, err := w.out.Write(w.buf.Bytes())
	w.buf.Reset()
	return err
}
