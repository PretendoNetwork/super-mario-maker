package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/PretendoNetwork/grpc/go/smm/v1"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func (s *gRPCSMMServer) DeleteCourse(ctx context.Context, in *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	dataID := types.NewUInt64(in.DataId)
	err := datastore_db.IsObjectAvailable(dataID)
	if err != nil {
		return &pb.DeleteCourseResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "Course not found")
	}

	metaInfo, err := datastore_db.GetObjectInfoByDataID(dataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return &pb.DeleteCourseResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	err = datastore_db.DeleteObjectByDataID(dataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return &pb.DeleteCourseResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.DeleteCourseResponse{
		CourseName: string(metaInfo.Name),
		Success: true,
	}, nil
}
