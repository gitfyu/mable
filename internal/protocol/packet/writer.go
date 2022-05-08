package packet

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gitfyu/mable/internal/protocol"
)

// Writer is used to write packets.
type Writer struct {
	writer    io.Writer
	varIntBuf bytes.Buffer
}

// NewWriter constructs a new Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) writeVarInt(v int32) error {
	w.varIntBuf.Reset()
	protocol.WriteVarInt(&w.varIntBuf, v)
	if _, err := w.writer.Write(w.varIntBuf.Bytes()); err != nil {
		return err
	}

	return nil
}

// WritePacket writes a single packet, including its length and id.
func (w *Writer) WritePacket(pk Outbound) error {
	buf := protocol.AcquireWriteBuffer()
	defer protocol.ReleaseWriteBuffer(buf)

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
