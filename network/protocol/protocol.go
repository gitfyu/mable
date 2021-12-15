package protocol

type State int

const (
	StateHandshake State = iota
	StateStatus
)
