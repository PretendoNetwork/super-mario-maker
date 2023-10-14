package globals

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/minio/minio-go/v7"
)

var Logger = plogger.NewLogger()
var NEXServer *nex.Server
var MinIOClient *minio.Client
var Presigner *S3Presigner
var DataStoreIDGenerators []*database.DataStoreIDGenerator
