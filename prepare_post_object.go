package main

import (
	"encoding/hex"
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func preparePostObject(err error, client *nex.Client, callID uint32, dataStorePreparePostParam *nexproto.DataStorePreparePostParam) {
	fmt.Println(hex.EncodeToString(dataStorePreparePostParam.MetaBinary))

	fieldBucket := nexproto.NewDataStoreKeyValue()
	fieldBucket.Key = "bucket"
	fieldBucket.Value = "pds-amaj-d1"

	fieldKey := nexproto.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = "course/1.bin"

	fieldACL := nexproto.NewDataStoreKeyValue()
	fieldACL.Key = "acl"
	fieldACL.Value = "private"

	rmcResponseStream := nex.NewStreamOut(nexServer)

	pReqPostInfo := nexproto.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = 1
	pReqPostInfo.URL = "http://datastore.pretendo.cc/upload"
	pReqPostInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*nexproto.DataStoreKeyValue{fieldBucket, fieldKey, fieldACL}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodPreparePostObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
