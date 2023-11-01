package globals

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func S3StatObject(bucket, key string) (minio.ObjectInfo, error) {
	return MinIOClient.StatObject(context.TODO(), bucket, key, minio.StatObjectOptions{})
}

func S3ObjectSize(bucket, key string) (uint64, error) {
	info, err := S3StatObject(bucket, key)
	if err != nil {
		return 0, err
	}

	return uint64(info.Size), nil
}
