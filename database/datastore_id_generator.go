package database

import (
	"encoding/binary"
)

// DataStoreIDGenerator represents a safe, structured, ID for unique NEX DataStore objects
type DataStoreIDGenerator struct {
	inUse    bool
	regionID uint8
	nodeID   uint8
	Value    uint32
}

// InUse is used to check if a DataStoreIDGenerator Node is currently generating a number
func (dataStoreIDGenerator *DataStoreIDGenerator) InUse() bool {
	return dataStoreIDGenerator.inUse
}

// SetInUse is used to check if a DataStoreIDGenerator Node is currently generating a number
func (dataStoreIDGenerator *DataStoreIDGenerator) SetInUse(inUse bool) {
	dataStoreIDGenerator.inUse = inUse
}

// Next is used to get the next value from the DataStoreIDGenerator Node
func (dataStoreIDGenerator *DataStoreIDGenerator) Next() uint64 {
	if dataStoreIDGenerator.inUse {
		panic("Cannot call DataStoreIDGenerator.Next while in use")
	}

	dataStoreIDGenerator.Value += 1

	data := make([]byte, 8)

	data[2] = dataStoreIDGenerator.regionID
	data[3] = dataStoreIDGenerator.nodeID
	binary.BigEndian.PutUint32(data[4:], dataStoreIDGenerator.Value)

	return binary.BigEndian.Uint64(data)
}

func NewDataStoreIDGenerator(regionID, nodeID uint8, value uint32) *DataStoreIDGenerator {
	return &DataStoreIDGenerator{
		regionID: regionID,
		nodeID:   nodeID,
		Value:    value,
	}
}
