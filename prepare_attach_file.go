package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func prepareAttachFile(err error, client *nex.Client, callID uint32, dataStoreAttachFileParam *nexproto.DataStoreAttachFileParam) {
	key := "image/1.jpg"
	bucket := "pds-amaj-d1"
	date := strconv.Itoa(int(time.Now().Unix()))
	pid := strconv.Itoa(int(client.PID()))

	data := pid + bucket + key + date

	hmac := hmac.New(sha256.New, hmacSecret)
	hmac.Write([]byte(data))

	signature := hex.EncodeToString(hmac.Sum(nil))

	fieldBucket := nexproto.NewDataStoreKeyValue()
	fieldBucket.Key = "bucket"
	fieldBucket.Value = bucket

	fieldKey := nexproto.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldACL := nexproto.NewDataStoreKeyValue()
	fieldACL.Key = "acl"
	fieldACL.Value = "public-read"

	fieldContentType := nexproto.NewDataStoreKeyValue()
	fieldContentType.Key = "content-type"
	fieldContentType.Value = "image/jpeg"

	fieldPID := nexproto.NewDataStoreKeyValue()
	fieldPID.Key = "pid"
	fieldPID.Value = pid

	fieldDate := nexproto.NewDataStoreKeyValue()
	fieldDate.Key = "date"
	fieldDate.Value = date

	fieldSignature := nexproto.NewDataStoreKeyValue()
	fieldSignature.Key = "signature"
	fieldSignature.Value = signature

	rmcResponseStream := nex.NewStreamOut(nexServer)

	pReqPostInfo := nexproto.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = 1
	pReqPostInfo.URL = "http://datastore.pretendo.cc/upload"
	pReqPostInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*nexproto.DataStoreKeyValue{fieldBucket, fieldKey, fieldACL, fieldContentType, fieldPID, fieldDate, fieldSignature}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodPrepareAttachFile, rmcResponseBody)

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
