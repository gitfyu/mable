package packet

// Client->Server

const (
	StatusRequest ID = iota
	Ping
)

// Server->Client

const (
	StatusResponse ID = iota
	Pong
)
