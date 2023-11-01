package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetCustomRankingByDataID(err error, packet nex.PacketInterface, callID uint32, param *datastore_super_mario_maker_types.DataStoreGetCustomRankingByDataIDParam) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	client := packet.Sender()

	pRankingResult := datastore_smm_db.GetCustomRankingsByDataIDs(param.ApplicationID, param.DataIDList)
	pResults := make([]*nex.Result, 0, len(param.DataIDList))

	for _, rankingResult := range pRankingResult {
		// * This is kind of backwards.
		// * The database pulls this data
		// * by default, so it can be done
		// * in a single query. So instead
		// * of checking if a flag *IS*
		// * set, and conditionally *ADDING*
		// * the fields, we check if a flag
		// * is *NOT* set and conditionally
		// * *REMOVE* the field
		if param.ResultOption&0x1 == 0 {
			rankingResult.MetaInfo.Tags = make([]string, 0)
		}

		if param.ResultOption&0x2 == 0 {
			rankingResult.MetaInfo.Ratings = make([]*datastore_types.DataStoreRatingInfoWithSlot, 0)
		}

		if param.ResultOption&0x4 == 0 {
			rankingResult.MetaInfo.MetaBinary = make([]byte, 0)
		}

		if param.ResultOption&0x20 == 0 {
			rankingResult.Score = 0
		}

		// * Since all errors are thrown away in
		// * datastore_smm_db.GetCustomRankingsByDataIDs
		// * assume all objects returned were a success
		pResults = append(pResults, nex.NewResultSuccess(nex.Errors.Core.Unknown))
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(pRankingResult)
	rmcResponseStream.WriteListResult(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetCustomRankingByDataID, rmcResponseBody)

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
