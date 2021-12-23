package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/gitfyu/mable/chat"
	"github.com/gitfyu/mable/protocol"
	"github.com/google/uuid"
	"math"
	"sync"
)

var errInvalidLength = errors.New("invalid length received")

// Buffer is a container for packet data
type Buffer struct {
	buf bytes.Buffer
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &Buffer{}
	},
}

// AcquireBuffer returns a Buffer, which must be released afterwards using ReleaseBuffer
func AcquireBuffer() *Buffer {
	b := bufferPool.Get().(*Buffer)
	b.buf.Reset()
	return b
}

// ReleaseBuffer releases a Buffer, after which you must no longer use it
func ReleaseBuffer(b *Buffer) {
	bufferPool.Put(b)
}

func (b *Buffer) Reset() {
	b.buf.Reset()
}

func (b *Buffer) WriteUnsignedByte(v uint8) {
	b.buf.WriteByte(v)
}

func (b *Buffer) WriteSignedByte(v int8) {
	b.buf.WriteByte(byte(v))
}

func (b *Buffer) WriteBool(v bool) {
	if v {
		b.buf.WriteByte(1)
	} else {
		b.buf.WriteByte(0)
	}
}

func (b *Buffer) WriteUnsignedShort(v uint16) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, v)
	b.buf.Write(data)
}

func (b *Buffer) WriteUnsignedShortLittleEndian(v uint16) {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, v)
	b.buf.Write(data)
}

func (b *Buffer) ReadUnsignedShort() (uint16, error) {
	data := make([]byte, 2)
	if _, err := b.buf.Read(data); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(data), nil
}

func (b *Buffer) WriteInt(v int32) {
	b.WriteUnsignedInt(uint32(v))
}

func (b *Buffer) WriteUnsignedInt(v uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, v)
	b.buf.Write(data)
}

func (b *Buffer) WriteVarInt(v protocol.VarInt) {
	data := make([]byte, protocol.VarIntSize(v))
	protocol.WriteVarInt(data, v)
	b.buf.Write(data)
}

func (b *Buffer) ReadVarInt() (protocol.VarInt, error) {
	var v protocol.VarInt
	if err := protocol.ReadVarInt(&b.buf, &v); err != nil {
		return 0, err
	}

	return v, nil
}

func (b *Buffer) ReadLong() (int64, error) {
	data := make([]byte, 8)
	if _, err := b.buf.Read(data); err != nil {
		return 0, err
	}

	return int64(binary.BigEndian.Uint64(data)), nil
}

func (b *Buffer) WriteLong(v int64) {
	b.WriteUnsignedLong(uint64(v))
}

func (b *Buffer) WriteUnsignedLong(v uint64) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, v)
	b.buf.Write(data)
}

func (b *Buffer) WriteFloat(v float32) {
	b.WriteUnsignedInt(math.Float32bits(v))
}

func (b *Buffer) WriteDouble(v float64) {
	b.WriteUnsignedLong(math.Float64bits(v))
}

func (b *Buffer) WritePosition(x, y, z int32) {
	b.WriteUnsignedLong((uint64(x&0x3FFFFFF) << 38) | (uint64(y&0xFFF) << 26) | (uint64(z) & 0x3FFFFFF))
}

func (b *Buffer) WriteString(s string) {
	data := []byte(s)
	b.WriteVarInt(protocol.VarInt(len(data)))
	b.buf.Write(data)
}

func (b *Buffer) WriteStringFromBytes(data []byte) {
	b.WriteVarInt(protocol.VarInt(len(data)))
	b.buf.Write(data)
}

func (b *Buffer) ReadString() (string, error) {
	n, err := b.ReadVarInt()
	if err != nil {
		return "", err
	}

	if n < 0 {
		return "", errInvalidLength
	} else if n == 0 {
		return "", nil
	}

	data := make([]byte, n)
	if _, err := b.buf.Read(data); err != nil {
		return "", err
	}

	return string(data), nil
}

func (b *Buffer) WriteUUID(uuid uuid.UUID) {
	buf, _ := uuid.MarshalBinary()
	b.buf.Write(buf)
}

func (b *Buffer) WriteMsg(msg *chat.Msg) error {
	str, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	b.WriteStringFromBytes(str)
	return nil
}

// Write writes the given data to the buffer. It never returns an error.
func (b *Buffer) Write(data []byte) (int, error) {
	return b.buf.Write(data)
}

// Bytes returns the currently stored packet data. It is valid until the next read/write call.
func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

// Len returns the number of bytes currently stored in the Buffer
func (b *Buffer) Len() int {
	return b.buf.Len()
}
