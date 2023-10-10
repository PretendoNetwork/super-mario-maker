package nex_datastore

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
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func PreparePostObject(err error, client *nex.Client, callID uint32, param *datastore_types.DataStorePreparePostParam) uint32 {
	rand.Seed(time.Now().UnixNano())
	nodeID := rand.Intn(len(globals.DataStoreIDGenerators))

	dataStoreIDGenerator := globals.DataStoreIDGenerators[nodeID]

	dataID := dataStoreIDGenerator.Next()
	database.SetDataStoreIDGeneratorLastID(nodeID, dataStoreIDGenerator.Value)
	database.InitializeCourseData(dataID, client.PID(), param.Size, param.Name, param.Flag, param.ExtraData, param.DataType, param.Period)

	if param.DataType != 1 { // 1 is Mii data, assume other values are course meta data
		database.UpdateCourseMetaBinary(dataID, param.MetaBinary)
	}

	key := fmt.Sprintf("course/%d.bin", dataID)
	bucket := os.Getenv("S3_BUCKET_NAME")
	date := strconv.Itoa(int(time.Now().Unix()))
	pid := strconv.Itoa(int(client.PID()))

	data := pid + bucket + key + date

	hmac := hmac.New(sha256.New, []byte{})
	hmac.Write([]byte(data))

	signature := hex.EncodeToString(hmac.Sum(nil))

	fieldBucket := datastore_types.NewDataStoreKeyValue()
	fieldBucket.Key = "bucket"
	fieldBucket.Value = bucket

	fieldKey := datastore_types.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldACL := datastore_types.NewDataStoreKeyValue()
	fieldACL.Key = "acl"
	fieldACL.Value = "private"

	fieldPID := datastore_types.NewDataStoreKeyValue()
	fieldPID.Key = "pid"
	fieldPID.Value = pid

	fieldDate := datastore_types.NewDataStoreKeyValue()
	fieldDate.Key = "date"
	fieldDate.Value = date

	fieldSignature := datastore_types.NewDataStoreKeyValue()
	fieldSignature.Key = "signature"
	fieldSignature.Value = signature

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	pReqPostInfo := datastore_types.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = dataID
	pReqPostInfo.URL = os.Getenv("DATASTORE_UPLOAD_URL")
	pReqPostInfo.RequestHeaders = []*datastore_types.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*datastore_types.DataStoreKeyValue{fieldBucket, fieldKey, fieldACL, fieldPID, fieldDate, fieldSignature}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPreparePostObject, rmcResponseBody)

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
