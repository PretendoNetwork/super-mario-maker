package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func RateCustomRanking(err error, client *nex.Client, callID uint32, params []*datastore_super_mario_maker_types.DataStoreRateCustomRankingParam) uint32 {

	/*
		This has to change. We need to figure out what DataStoreRateCustomRankingParam.applicationId means.
		This method sends a set number of DataStoreRateCustomRankingParam's in the params list depending
		on what action the user does.

		- Course upload: 4 params
			- DataStoreRateCustomRankingParam.dataId is always the course ID
			- DataStoreRateCustomRankingParam.score is always 0
			- DataStoreRateCustomRankingParam.period is always 36500
			- DataStoreRateCustomRankingParam.applicationId is:
				- 0 (seen in other requests relating to Mii data?)
				- 2400
				- 200000000
				- 200002400

		- Course clear: 2 params
			- DataStoreRateCustomRankingParam.dataId is always the course ID
			- DataStoreRateCustomRankingParam.score is always 1
			- DataStoreRateCustomRankingParam.period is always 36500
			- DataStoreRateCustomRankingParam.applicationId is:
				- 200000000
				- 200002400

		- Course star: 12 params
			- DataStoreRateCustomRankingParam.dataId is:
				- course ID
				- course ID
				- course ID
				- course ID
				- course ID
				- course ID
				- uploader PID
				- uploader PID
				- uploader PID
				- uploader PID
				- uploader PID
				- uploader PID
			- DataStoreRateCustomRankingParam.score is always 1
			- DataStoreRateCustomRankingParam.period is:
				- 36500
				- 36500
				- 1
				- 1
				- 1
				- 1
				- 36500
				- 36500
				- 1
				- 1
				- 1
				- 1
			- DataStoreRateCustomRankingParam.applicationId is:
				- 0 (seen in other requests relating to Mii data?)
				- 2400
				- 119700101
				- 119702501
				- 1497730516
				- 1497732916
				- 300000000 (seen in other requests relating to Mii data?)
				- 300002400
				- 419700101
				- 419702501
				- 1797730516
				- 1797732916

			It seems like many applicationId's are actually a different ID with 2400 added to it?

			Since it always sends a set number of params per action, a quick and dirty way to tell them apart is
			to check the param list length. However this should NOT be relied on forever and we NEED to learn more
			about applicationId, as it IS used in many other requests
	*/

	if len(params) == 2 {
		// assume "course clear" action
		database.IncrementCourseClearCount(params[0].DataID)
	}

	if len(params) == 12 {
		// assume "star course" action
		database.IncrementCourseStarCount(params[0].DataID)
	}

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodRateCustomRanking, nil)

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

	return 0
}
