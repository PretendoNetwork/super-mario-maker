package database

import (
	"context"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"google.golang.org/grpc/metadata"
)

func GetUserMiiInfoByPID(pid uint32) *pb.Mii {
	ctx := metadata.NewOutgoingContext(context.Background(), globals.GRPCAccountCommonMetadata)

	response, err := globals.GRPCAccountClient.GetUserData(ctx, &pb.GetUserDataRequest{Pid: pid})
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil
	}

	return response.Mii
}
