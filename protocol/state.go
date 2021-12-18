package protocol

type State uint8

const (
	StateHandshake State = iota
	StateStatus
	StateLogin
	StatePlay
)
