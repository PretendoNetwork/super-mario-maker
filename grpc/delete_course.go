package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/PretendoNetwork/grpc/go/smm"
	datastore_db "github.com/PretendoNetwork/super-mario-maker/database/datastore"
	"github.com/PretendoNetwork/super-mario-maker/globals"
)

func (s *gRPCSMMServer) DeleteCourse(ctx context.Context, in *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	err := datastore_db.IsObjectAvailable(in.DataID)
	if err != nil {
		return &pb.DeleteCourseResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "Course not found")
	}

	metaInfo, err := datastore_db.GetObjectInfoByDataID(in.DataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return &pb.DeleteCourseResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	err = datastore_db.DeleteObjectByDataID(in.DataID)
	if err != nil {
		globals.Logger.Critical(err.Error())
		return &pb.DeleteCourseResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.SendUserFriendRequestResponse{
		CourseName: metaInfo.Name,
		Success: true,
	}, nil
}
