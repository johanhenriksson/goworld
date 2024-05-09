package srv

import "fmt"

// Identity is a unique identifier for a pool element
// 0-7: type id
// 8-31: generation
// 32-63: id
type Identity uint64

const None Identity = 0

func NewID(typeID, generation, index int) Identity {
	return Identity(uint64(typeID)<<56 | uint64(generation)<<32 | uint64(index))
}

func (id Identity) TypeID() int {
	return int((id >> 56) & 0xFF)
}

func (id Identity) Generation() int {
	return int((id >> 32) & 0xFFFFFF)
}

func (id Identity) Index() int {
	return int(id & 0xFFFFFFFF)
}

func (id Identity) String() string {
	return fmt.Sprintf("ID(%d, %d, %d)", id.TypeID(), id.Generation(), id.Index())
}
