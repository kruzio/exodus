package s3

import (
	"github.com/kruzio/exodus/pkg/cloudblob"
	"gocloud.dev/blob/s3blob"
)

type S3Uploader struct {
	cloudblob.CloudBlob
}

const Scheme = s3blob.Scheme

// Establish an AWS session.
// See https://docs.aws.amazon.com/sdk-for-go/api/aws/session/ for more info.
// The region must match the region for "my-bucket".
func (s3 *S3Uploader) UsageInfo() string {
	return "Upload to AWS S3 bucket using the following URL scheme: s3://bucket-name/subdir?region=us-west-1. For additional information see https://gocloud.dev/howto/blob/#s3"
}
