package utility

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/PretendoNetwork/nex-go"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/datastore/super-mario-maker/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/super-mario-maker-secure/database"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
	"github.com/PretendoNetwork/super-mario-maker-secure/types"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CourseMetadataToDataStoreCustomRankingResult(courseMetadata *types.CourseMetadata) *datastore_super_mario_maker_types.DataStoreCustomRankingResult {
	rankingResult := datastore_super_mario_maker_types.NewDataStoreCustomRankingResult()

	rankingResult.Order = 0 // unknown
	rankingResult.Score = courseMetadata.Stars
	rankingResult.MetaInfo = CourseMetadataToDataStoreMetaInfo(courseMetadata)

	return rankingResult
}

func CourseMetadataToDataStoreMetaInfo(courseMetadata *types.CourseMetadata) *datastore_types.DataStoreMetaInfo {
	metaInfo := datastore_types.NewDataStoreMetaInfo()

	metaInfo.DataID = courseMetadata.DataID
	metaInfo.OwnerID = courseMetadata.OwnerPID
	metaInfo.Size = courseMetadata.Size
	metaInfo.Name = courseMetadata.Name
	metaInfo.DataType = courseMetadata.DataType
	metaInfo.MetaBinary = courseMetadata.MetaBinary
	metaInfo.Permission = datastore_types.NewDataStorePermission()
	metaInfo.Permission.Permission = 0 // unknown
	metaInfo.Permission.RecipientIDs = []uint32{}
	metaInfo.DelPermission = datastore_types.NewDataStorePermission()
	metaInfo.DelPermission.Permission = 3 // unknown
	metaInfo.DelPermission.RecipientIDs = []uint32{}
	metaInfo.CreatedTime = courseMetadata.CreatedTime
	metaInfo.UpdatedTime = courseMetadata.UpdatedTime
	metaInfo.Period = courseMetadata.Period
	metaInfo.Status = 0      // unknown
	metaInfo.ReferredCnt = 0 // unknown
	metaInfo.ReferDataID = 0 // unknown
	metaInfo.Flag = courseMetadata.Flag
	metaInfo.ReferredTime = courseMetadata.CreatedTime
	metaInfo.ExpireTime = nex.NewDateTime(671075926016) // December 31st, year 9999
	metaInfo.Tags = []string{""}                        // unknown
	metaInfo.Ratings = []*datastore_types.DataStoreRatingInfoWithSlot{
		datastore_types.NewDataStoreRatingInfoWithSlot(), // attempts
		datastore_types.NewDataStoreRatingInfoWithSlot(), // unknown
		datastore_types.NewDataStoreRatingInfoWithSlot(), // completions
		datastore_types.NewDataStoreRatingInfoWithSlot(), // failures
		datastore_types.NewDataStoreRatingInfoWithSlot(), // unknown
		datastore_types.NewDataStoreRatingInfoWithSlot(), // unknown
		datastore_types.NewDataStoreRatingInfoWithSlot(), // unknown
	}

	// attempts
	metaInfo.Ratings[0].Slot = 0
	metaInfo.Ratings[0].Rating = datastore_types.NewDataStoreRatingInfo()
	metaInfo.Ratings[0].Rating.TotalValue = int64(courseMetadata.Attempts)
	metaInfo.Ratings[0].Rating.Count = courseMetadata.Attempts
	metaInfo.Ratings[0].Rating.InitialValue = 0

	// unknown
	metaInfo.Ratings[1].Slot = 1
	metaInfo.Ratings[1].Rating = datastore_types.NewDataStoreRatingInfo()
	metaInfo.Ratings[1].Rating.TotalValue = 2
	metaInfo.Ratings[1].Rating.Count = 2
	metaInfo.Ratings[1].Rating.InitialValue = 0

	// completions
	metaInfo.Ratings[2].Slot = 2
	metaInfo.Ratings[2].Rating = datastore_types.NewDataStoreRatingInfo()
	metaInfo.Ratings[2].Rating.TotalValue = int64(courseMetadata.Completions)
	metaInfo.Ratings[2].Rating.Count = courseMetadata.Completions
	metaInfo.Ratings[2].Rating.InitialValue = 0

	// failures
	metaInfo.Ratings[3].Slot = 3
	metaInfo.Ratings[3].Rating = datastore_types.NewDataStoreRatingInfo()
	metaInfo.Ratings[3].Rating.TotalValue = int64(courseMetadata.Failures)
	metaInfo.Ratings[3].Rating.Count = courseMetadata.Failures
	metaInfo.Ratings[3].Rating.InitialValue = 0

	// unknown
	metaInfo.Ratings[4].Slot = 4
	metaInfo.Ratings[4].Rating = datastore_types.NewDataStoreRatingInfo()
	metaInfo.Ratings[4].Rating.TotalValue = 5
	metaInfo.Ratings[4].Rating.Count = 5
	metaInfo.Ratings[4].Rating.InitialValue = 0

	// unknown
	metaInfo.Ratings[5].Slot = 5
	metaInfo.Ratings[5].Rating = datastore_types.NewDataStoreRatingInfo()
	metaInfo.Ratings[5].Rating.TotalValue = 6
	metaInfo.Ratings[5].Rating.Count = 6
	metaInfo.Ratings[5].Rating.InitialValue = 0

	// Number of new Miiverse comments
	metaInfo.Ratings[6].Slot = 6
	metaInfo.Ratings[6].Rating = datastore_types.NewDataStoreRatingInfo()
	metaInfo.Ratings[6].Rating.TotalValue = 0
	metaInfo.Ratings[6].Rating.Count = 0
	metaInfo.Ratings[6].Rating.InitialValue = 0

	return metaInfo
}

func UserMiiDataToDataStoreCustomRankingResult(ownerID uint32, miiInfo primitive.M) *datastore_super_mario_maker_types.DataStoreCustomRankingResult {
	rankingResult := datastore_super_mario_maker_types.NewDataStoreCustomRankingResult()

	rankingResult.Order = 0
	rankingResult.Score = 0
	rankingResult.MetaInfo = UserMiiDataToDataStoreMetaInfo(ownerID, miiInfo)

	return rankingResult
}

func UserMiiDataToDataStoreMetaInfo(ownerID uint32, miiInfo primitive.M) *datastore_types.DataStoreMetaInfo {
	encodedMiiData := miiInfo["data"].(string)
	decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

	metaBinaryStream := nex.NewStreamOut(globals.NEXServer)
	metaBinaryStream.Grow(140)
	metaBinaryStream.WriteBytesNext([]byte{
		0x42, 0x50, 0x46, 0x43, // BPFC magic
		0x00, 0x00, 0x00, 0x01, // Unknown
		0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x01, 0x00, 0x00, // Unknown
	})
	metaBinaryStream.WriteBytesNext(decodedMiiData) // Actual Mii data
	metaBinaryStream.WriteBytesNext([]byte{
		0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x00, 0x00, 0x00, // Unknown
		0x00, 0x00, 0x00, 0x01, // Unknown
	})

	now := uint64(time.Now().Unix())

	metaInfo := datastore_types.NewDataStoreMetaInfo()
	metaInfo.DataID = uint64(ownerID) // This isn;t actually a user PID in Nintendo's servers, but it makes it much easier for us to do it this way
	metaInfo.OwnerID = ownerID
	metaInfo.Size = 0
	metaInfo.Name = miiInfo["name"].(string)
	metaInfo.DataType = 1 // Mii data type?
	metaInfo.MetaBinary = metaBinaryStream.Bytes()
	metaInfo.Permission = datastore_types.NewDataStorePermission()
	metaInfo.Permission.Permission = 0 // idk?
	metaInfo.Permission.RecipientIDs = []uint32{}
	metaInfo.DelPermission = datastore_types.NewDataStorePermission()
	metaInfo.DelPermission.Permission = 3 // idk?
	metaInfo.DelPermission.RecipientIDs = []uint32{}
	metaInfo.CreatedTime = nex.NewDateTime(now)
	metaInfo.UpdatedTime = nex.NewDateTime(now)
	metaInfo.Period = 90 // idk?
	metaInfo.Status = 0
	metaInfo.ReferredCnt = 0
	metaInfo.ReferDataID = 0
	metaInfo.Flag = 256 // idk?
	metaInfo.ReferredTime = nex.NewDateTime(now)
	metaInfo.ExpireTime = nex.NewDateTime(now)
	metaInfo.Tags = []string{"49"} // idk?
	metaInfo.Ratings = []*datastore_types.DataStoreRatingInfoWithSlot{}

	return metaInfo
}

func UserNotOwnCourse(courseID uint64, pid uint32) bool {
	courseMetadata := database.GetCourseMetadataByDataID(courseID)

	return courseMetadata.OwnerPID != pid
}

func S3StatObject(bucket, key string) (minio.ObjectInfo, error) {
	return globals.MinIOClient.StatObject(context.TODO(), bucket, key, minio.StatObjectOptions{})
}

func S3ObjectSize(bucket, key string) (uint64, error) {
	info, err := S3StatObject(bucket, key)
	if err != nil {
		return 0, err
	}

	return uint64(info.Size), nil
}
