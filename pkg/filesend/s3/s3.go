package s3

import (
	"github.com/kruzio/exodus/pkg/cloudblob"
	"github.com/kruzio/exodus/pkg/usageprinter"
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
	table, buf := usageprinter.NewUsageTable(Scheme)

	data := [][]string{
		[]string{"Upload to AWS S3 bucket", "s3://bucket-name/subdir?region=us-west-1\nFor additional information see https://gocloud.dev/howto/blob/#s3"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (s3 *S3Uploader) Scheme() string {
	return Scheme
}
