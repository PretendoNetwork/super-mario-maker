package main

import (
	"fmt"
	"log"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol
var s3Client *s3.S3

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("DO_SPACES_KEY")
	secret := os.Getenv("DO_SPACES_SECRET")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("http://b-cdn.net"),
		Region:      aws.String("us-east-1"),
	}

	newSession, _ := session.NewSession(s3Config)
	s3Client = s3.New(newSession)

	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(4)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetAccessKey("9f2b4678")

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==SMM1 - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("=================")
	})

	secureServer = nexproto.NewSecureProtocol(nexServer)
	dataStoreSmmServer := nexproto.NewDataStoreSMMProtocol(nexServer)

	// Handle PRUDP CONNECT packet (not an RMC method)
	nexServer.On("Connect", connect)

	secureServer.Register(register)

	dataStoreSmmServer.GetMeta(getMeta)
	dataStoreSmmServer.PreparePostObject(preparePostObject)
	dataStoreSmmServer.PrepareGetObject(prepareGetObject)
	dataStoreSmmServer.CompletePostObject(completePostObject)
	dataStoreSmmServer.GetMetasMultipleParam(getMetasMultipleParam)
	dataStoreSmmServer.ChangeMeta(changeMeta)
	dataStoreSmmServer.RateObjects(rateObjects)
	dataStoreSmmServer.GetObjectInfos(getObjectInfos)
	dataStoreSmmServer.RateCustomRanking(rateCustomRanking)
	dataStoreSmmServer.GetCustomRankingByDataId(getCustomRankingByDataId)
	dataStoreSmmServer.GetBufferQueue(getBufferQueue)
	dataStoreSmmServer.CompleteAttachFile(completeAttachFile)
	dataStoreSmmServer.PrepareAttachFile(prepareAttachFile)
	dataStoreSmmServer.GetApplicationConfig(getApplicationConfig)
	dataStoreSmmServer.FollowingsLatestCourseSearchObject(followingsLatestCourseSearchObject)
	dataStoreSmmServer.RecommendedCourseSearchObject(recommendedCourseSearchObject)
	dataStoreSmmServer.SuggestedCourseSearchObject(suggestedCourseSearchObject)
	dataStoreSmmServer.GetCourseRecord(getCourseRecord)
	dataStoreSmmServer.GetApplicationConfigString(getApplicationConfigString)
	dataStoreSmmServer.GetMetasWithCourseRecord(getMetasWithCourseRecord)
	dataStoreSmmServer.CTRPickUpCourseSearchObject(ctrPickUpCourseSearchObject)

	nexServer.Listen(":60003")
}
