package packet

import (
	"bytes"
	"io"

	"github.com/gitfyu/mable/internal/protocol"
)

// Writer is used to write packets.
type Writer struct {
	out io.Writer
	buf bytes.Buffer
}

// NewWriter constructs a new Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		out: w,
	}
}

// WritePacket adds a single packet to the internal buffer, which will be written the next
// time that Flush is called.
func (w *Writer) WritePacket(pk Outbound) {
	data := protocol.AcquireWriteBuffer()
	defer protocol.ReleaseWriteBuffer(data)

	data.Reset()
	data.WriteVarInt(int32(pk.PacketID()))
	pk.MarshalPacket(data)

	// Write length
	protocol.WriteVarInt(&w.buf, int32(data.Len()))
	// Write ID + data
	w.buf.Write(data.Bytes())
}

// Writes the internal buffer to the io.Writer that was used to construct this Writer.
func (w *Writer) Flush() error {
	_, err := w.out.Write(w.buf.Bytes())
	w.buf.Reset()
	return err
}
