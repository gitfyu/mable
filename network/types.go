package network

type (
	Error   string
	VarInt  int32
	VarLong int64
)

const (
	varIntMaxBytes  = 5
	varLongMaxBytes = 10
	// The protocol allows for much longer strings, but to prevent abuse this is capped at a much lower arbitrary value
	// TODO: ensure that this limit is high enough to not cause issues with legitimate clients
	stringMaxBytes = 1024

	ErrVarIntTooBig = Error("VarInt too big")
	ErrStringNegLen = Error("String has negative length")
	ErrStringTooBig = Error("String too big")
)

func (e Error) Error() string {
	return string(e)
}
