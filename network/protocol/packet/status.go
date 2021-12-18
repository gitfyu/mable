package packet

// Client->Server

const (
	StatusRequest ID = iota
	StatusPing
)

// Server->Client

const (
	StatusResponse ID = iota
	StatusPong
)
