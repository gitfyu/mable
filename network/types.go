package network

type (
	Error   string
	VarInt  int32
	VarLong int64
)

const (
	varIntMaxBytes   = 5
	varLongMaxBytes  = 10
	ErrVarIntTooBig  = Error("VarInt too big")
	ErrVarLongTooBig = Error("VarLong too big")
)

func (e Error) Error() string {
	return string(e)
}
