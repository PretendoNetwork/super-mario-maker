package grpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/PretendoNetwork/grpc/go/smm/v1"
	"github.com/PretendoNetwork/super-mario-maker/globals"
	"google.golang.org/grpc"
)

type gRPCSMMServer struct {
	pb.UnimplementedSMMServiceServer
}

func StartGRPCServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", globals.GRPCServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(apiKeyInterceptor),
	)

	pb.RegisterSMMServiceServer(server, &gRPCSMMServer{})

	log.Printf("server listening at %v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
