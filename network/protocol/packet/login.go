package packet

// Client->Server

const (
	LoginStart ID = iota
)

// Server->Client

const (
	LoginDisconnect ID = iota
)
