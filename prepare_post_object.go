package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func preparePostObject(err error, client *nex.Client, callID uint32, param *nexproto.DataStorePreparePostParam) {
	rand.Seed(time.Now().UnixNano())
	nodeID := rand.Intn(len(dataStoreIDGenerators))

	dataStoreIDGenerator := dataStoreIDGenerators[nodeID]

	dataID := dataStoreIDGenerator.Next()
	setDataStoreIDGeneratorLastID(nodeID, dataStoreIDGenerator.Value)
	initializeCourseData(dataID, client.PID(), param.Size, param.Name, param.Flag, param.ExtraData, param.DataType, param.Period)

	if param.DataType != 1 { // 1 is Mii data, assume other values are course meta data
		updateCourseMetaBinary(dataID, param.MetaBinary)
	}

	key := fmt.Sprintf("course/%d.bin", dataID)
	bucket := os.Getenv("DO_SPACES_NAME")
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
	fieldACL.Value = "private"

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

	pReqPostInfo.DataID = dataID
	pReqPostInfo.URL = os.Getenv("DATASTORE_UPLOAD_URL")
	pReqPostInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*nexproto.DataStoreKeyValue{fieldBucket, fieldKey, fieldACL, fieldPID, fieldDate, fieldSignature}
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
