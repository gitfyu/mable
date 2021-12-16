package network

type (
	Error string
)

const (
	// The protocol allows for much longer strings, but to prevent abuse this is capped at a much lower arbitrary value
	// TODO: ensure that this limit is high enough to not cause issues with legitimate clients
	stringMaxBytes = 1024

	ErrStringNegLen = Error("string has negative length")
	ErrStringTooBig = Error("string too big")
)

func (e Error) Error() string {
	return string(e)
}
