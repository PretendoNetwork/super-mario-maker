package nex_datastore

import (
	nex "github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker-secure/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func RateObjects(err error, client *nex.Client, callID uint32, targets []*datastore_types.DataStoreRatingTarget, params []*datastore_types.DataStoreRateObjectParam, transactional bool, fetchRatings bool) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.DataStore.Unknown
	}

	pRatings := make([]*datastore_types.DataStoreRatingInfo, 0)
	pResults := make([]*nex.Result, 0)

	// * Real DataStore does not actually check this.
	// * I just didn't feel like working out the
	// * logic for differing sized lists. So force
	// * them to always be the same
	if len(targets) != len(params) {
		return nex.Errors.DataStore.InvalidArgument
	}

	for i := 0; i < len(targets); i++ {
		target := targets[i]
		param := params[i]

		rating, errCode := datastore_db.RateObjectWithPassword(target.DataID, target.Slot, param.RatingValue, param.AccessPassword)
		if errCode != 0 {
			return errCode
		}

		if fetchRatings {
			pRatings = append(pRatings, rating)
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListStructure(pRatings)
	rmcResponseStream.WriteListResult(pResults) // * pResults is always empty

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodRateObjects, rmcResponseBody)

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
