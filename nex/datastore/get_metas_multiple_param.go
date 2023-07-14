package nex_datastore

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/utility"
)

func GetMetasMultipleParam(err error, client *nex.Client, callID uint32, params []*datastore_types.DataStoreGetMetaParam) {
	pMetaInfo := make([]*datastore_types.DataStoreMetaInfo, 0)
	pResults := make([]uint32, 0)

	for _, param := range params {
		if param.DataID == 0 {
			pMetaInfo = append(pMetaInfo, getMetasMultipleParamMiiData(param))
		} else {
			fmt.Println("Unknown meta multiple data ID")
		}

		pResults = append(pResults, 0x690001)
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(pMetaInfo)
	rmcResponseStream.WriteListUInt32LE(pResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodGetMetasMultipleParam, rmcResponseBody)

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

func getMetasMultipleParamMiiData(param *datastore_types.DataStoreGetMetaParam) *datastore_types.DataStoreMetaInfo {
	miiInfo := database.GetUserMiiInfoByPID(param.PersistenceTarget.OwnerID)

	return utility.UserMiiDataToDataStoreMetaInfo(param.PersistenceTarget.OwnerID, miiInfo)
}
