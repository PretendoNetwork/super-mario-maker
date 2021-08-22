package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol

func main() {
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
	dataStoreSmmServer.PrepareGetObject(prepareGetObject)
	dataStoreSmmServer.GetMetasMultipleParam(getMetasMultipleParam)
	dataStoreSmmServer.ChangeMeta(changeMeta)
	dataStoreSmmServer.RateCustomRanking(rateCustomRanking)
	dataStoreSmmServer.GetCustomRankingByDataId(getCustomRankingByDataId)
	dataStoreSmmServer.GetBufferQueue(getBufferQueue)
	dataStoreSmmServer.GetApplicationConfig(getApplicationConfig)
	dataStoreSmmServer.FollowingsLatestCourseSearchObject(followingsLatestCourseSearchObject)
	dataStoreSmmServer.RecommendedCourseSearchObject(recommendedCourseSearchObject)
	dataStoreSmmServer.GetApplicationConfigString(getApplicationConfigString)
	dataStoreSmmServer.GetMetasWithCourseRecord(getMetasWithCourseRecord)

	nexServer.Listen(":60003")
}
