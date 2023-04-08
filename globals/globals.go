package globals

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Logger = plogger.NewLogger()
var NEXServer *nex.Server
var S3Client *s3.Client
var HMACSecret []byte
var DataStoreIDGenerators []*database.DataStoreIDGenerator
