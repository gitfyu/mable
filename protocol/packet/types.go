package packet

type ID uint8

// Client -> Server

const (
	Handshake ID = 0x00

	StatusRequest ID = 0x00
	StatusPing    ID = 0x01

	LoginStart   ID = 0x00
	LoginSuccess ID = 0x02
)

// Server -> Client

const (
	StatusResponse ID = 0x00
	StatusPong     ID = 0x01

	LoginDisconnect ID = 0x00

	PlayDisconnect ID = 0x40
)
