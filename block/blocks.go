package block

// ID represents a 12-bit block ID
type ID uint16

const (
	// maxID is the maximum value that can be stored in ID
	maxID = 1<<12 - 1

	// maxMetadata is the maximum value that can be stored in block metadata
	maxMetadata = 1<<4 - 1
)

const (
	Stone ID = 1
)

// Data encodes a block ID with metadata
type Data uint16

// ToUint16 encodes a Data value to an uint16, to be used in packets
func (d Data) ToUint16() uint16 {
	return uint16(d)
}

// Type returns the type of block stored in this Data
func (d Data) Type() ID {
	return ID(d >> 4)
}

// Metadata returns the metadata stored in this Data
func (d Data) Metadata() uint8 {
	return uint8(d) & 15
}

// ToData creates a Data value for this block without any metadata
func (id ID) ToData() Data {
	return Data(id << 4)
}

// ToDataWithMetadata creates a Data value for this block with specified metadata
func (id ID) ToDataWithMetadata(metadata uint8) Data {
	return Data(id)<<4 | Data(metadata)&15
}
