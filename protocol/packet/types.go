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

	PlayServerKeepAlive   ID = 0x00
	PlayServerJoinGame    ID = 0x01
	PlayServerChatMessage ID = 0x02
	PlayServerPosAndLook  ID = 0x08
	PlayServerChunkData   ID = 0x21
	PlayServerDisconnect  ID = 0x40
)
