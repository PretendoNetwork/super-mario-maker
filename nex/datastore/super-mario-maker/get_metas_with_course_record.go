package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetMetasWithCourseRecord(err error, client *nex.Client, callID uint32, dataStoreGetCourseRecordParams []*datastore_super_mario_maker.DataStoreGetCourseRecordParam, dataStoreGetMetaParam *datastore.DataStoreGetMetaParam) {
	// TODO: complete this

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteUInt32LE(0x00000000) // pMetaInfo List length 0
	rmcResponseStream.WriteUInt32LE(0x00000000) // pCourseResults List length 0
	rmcResponseStream.WriteUInt32LE(0x00000000) // pResults List length 0

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetMetasWithCourseRecord, rmcResponseBody)

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
