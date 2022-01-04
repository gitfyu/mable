package block

type ID uint16

const (
	Stone ID = 1
)

// Data represents a block type with metadata
type Data struct {
	// ID is the type of the block
	ID ID

	// Metadata contains additional data for the block
	Metadata uint8
}

// ToUint16 encodes a Data value to an uint16, to be used in packets
func (b Data) ToUint16() uint16 {
	return uint16(b.ID)<<4 | uint16(b.Metadata)&16
}

// ToData creates a Data value for this block without any metadata
func (id ID) ToData() Data {
	return Data{ID: id}
}

// ToDataWithMetadata creates a Data value for this block with specified metadata
func (id ID) ToDataWithMetadata(metadata uint8) Data {
	return Data{ID: id, Metadata: metadata}
}
