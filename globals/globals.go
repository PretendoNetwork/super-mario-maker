package globals

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/PretendoNetwork/grpc-go/account"
)

var Logger = plogger.NewLogger()
var NEXServer *nex.Server
var MinIOClient *minio.Client
var Presigner *S3Presigner
var DataStoreIDGenerators []*DataStoreIDGenerator
var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb.AccountClient
var GRPCAccountCommonMetadata metadata.MD
