package main

import "github.com/PretendoNetwork/nex-go"

type CourseMetadata struct {
	DataID      uint64
	OwnerPID    uint32
	Size        uint32
	Name        string
	MetaBinary  []byte
	Stars       uint32
	Attempts    uint32
	Failures    uint32
	Completions uint32
	Flag        uint32
	DataType    uint16
	Period      uint16
}

type CourseWorldRecord struct {
	FirstPID    uint32
	BestPID     uint32
	CreatedTime *nex.DateTime
	UpdatedTime *nex.DateTime
	Score       int32
}
