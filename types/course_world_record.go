package types

import "github.com/PretendoNetwork/nex-go"

type CourseWorldRecord struct {
	FirstPID    uint32
	BestPID     uint32
	CreatedTime *nex.DateTime
	UpdatedTime *nex.DateTime
	Score       int32
}
