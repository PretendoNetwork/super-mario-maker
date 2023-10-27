package nex_datastore_super_mario_maker

import (
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func SuggestedCourseSearchObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStoreSearchParam, extraData []string) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	// * This method is called when a course is completed
	// * to show the scrolling courses at the bottom of the
	// * screen. extraData[0] is the DataID of the current
	// * course. extraData[1] and extraData[2] are usually
	// * 2 and 6 respectively. extraData[3] and extraData[4]
	// * are always 0? extraData seems to not have any
	// * effect on the NUMBER of courses returned but likely
	// * does act as a filter of some kind? Maybe it has to
	// * do with difficulty? Or ratings?

	_, err = strconv.ParseUint(extraData[0], 0, 64)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.InvalidArgument
	}

	// TODO - Use extraData for filtering
	pRankingResults, errCode := datastore_smm_db.GetRandomCoursesWithLimit(int(param.ResultRange.Length))
	if errCode != 0 {
		return errCode
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

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

	globals.SecureServer.Send(responsePacket)

	return 0
}
