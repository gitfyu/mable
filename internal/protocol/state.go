package protocol

// State represents the state of the protocol.
type State uint8

const (
	// StateHandshake is the initial state.
	StateHandshake State = iota

	// StateStatus is used when a client wants to request information to display in the server list menu.
	StateStatus

	// StateLogin is used when the client wants to authenticate.
	StateLogin

	// StatePlay is used after successful authentication and a player instance has been created.
	StatePlay
)
