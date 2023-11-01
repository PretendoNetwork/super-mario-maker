package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_smm_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func CTRPickUpCourseSearchObject(err error, packet nex.PacketInterface, callID uint32, param *datastore_types.DataStoreSearchParam, extraData []string) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	client := packet.Sender()

	// * This method is only used by the 3DS version
	// * of Super Mario Maker and is functionally
	// * identical to the WiiU versions
	// * DataStoreSMM::RecommendedCourseSearchObject
	// * method. The 3DS version LIKELY uses this as
	// * a way to filter out courses which use the
	// * Mystery Mushroom powerup, since it's not
	// * officially supported by the 3DS version. That
	// * said, the Mystery Mushroom IS still in the
	// * game and can be used somewhat normally, so
	// * no filtering is done here to prevent that.
	// * I'm not even sure how we would detect that

	// TODO - Research extraData
	// TODO - Use the offet? Real client never uses it, but might be nice for completeness sake?
	pRankingResults, errCode := datastore_smm_db.GetRandomCoursesWithLimit(int(param.ResultRange.Length))
	if errCode != 0 {
		return errCode
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodCTRPickUpCourseSearchObject, rmcResponseBody)

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
