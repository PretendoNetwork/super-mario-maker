package main

import (
	"sync"

	"github.com/PretendoNetwork/super-mario-maker/grpc"
	"github.com/PretendoNetwork/super-mario-maker/nex"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	go grpc.StartGRPCServer()
	go nex.StartAuthenticationServer()
	go nex.StartSecureServer()

	wg.Wait()
}
