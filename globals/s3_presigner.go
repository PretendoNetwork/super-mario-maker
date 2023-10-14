package globals

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
)

type S3Presigner struct {
	minio *minio.Client
}

func (p *S3Presigner) GetObject(bucket, key string, lifetime time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)

	return p.minio.PresignedGetObject(context.Background(), bucket, key, lifetime, reqParams)
}

func (p *S3Presigner) PostObject(bucket, key string, lifetime time.Duration) (*url.URL, map[string]string, error) {
	policy := minio.NewPostPolicy()

	policy.SetBucket(bucket)
	policy.SetKey(key)
	policy.SetExpires(time.Now().UTC().Add(lifetime).UTC())

	policy.SetCondition("eq", "$bucket", "pn-amaj-d1")
	policy.SetCondition("eq", "$key", "17179869186.bin")
	policy.SetCondition("eq", "$x-amz-algorithm", "AWS4-HMAC-SHA256")
	policy.SetCondition("starts-with", "$x-amz-credential", os.Getenv("PN_SMM_CONFIG_S3_ACCESS_KEY"))
	policy.SetCondition("starts-with", "$x-amz-date", time.Now().Format("2006-01-02"))

	return p.minio.PresignedPostPolicy(context.Background(), policy)
}

func NewS3Presigner(minioClient *minio.Client) *S3Presigner {
	return &S3Presigner{
		minio: minioClient,
	}
}
