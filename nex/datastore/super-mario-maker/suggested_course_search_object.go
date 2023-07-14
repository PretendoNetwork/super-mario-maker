package nex_datastore_super_mario_maker

import (
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/utility"
)

func SuggestedCourseSearchObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreSearchParam, extraData []string) {
	// TODO: complete this

	courseID, _ := strconv.ParseUint(extraData[0], 0, 64)

	if utility.UserNotOwnCourse(courseID, client.PID()) {
		database.IncrementCourseAttemptCount(courseID) // We also know this is when a user attempts a course
	}

	pRankingResults := make([]*datastore_super_mario_maker_types.DataStoreCustomRankingResult, 0)

	courseMetadatas := database.GetCourseMetadatasByLimit(4) // In PCAPs param.minimalRatingFrequency is 4 but is 0 here?

	for _, courseMetadata := range courseMetadatas {
		pRankingResults = append(pRankingResults, utility.CourseMetadataToDataStoreCustomRankingResult(courseMetadata))
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodSuggestedCourseSearchObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
