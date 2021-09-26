package main

type CourseMetadata struct {
	Stars      uint32
	DataID     uint64
	OwnerPID   uint32
	Size       uint32
	Name       string
	MetaBinary []byte
	Flag       uint32
	DataType   uint16
	Period     uint16
}
