package gcs

import (
	"github.com/kruzio/exodus/pkg/cloudblob"
	"github.com/kruzio/exodus/pkg/usageprinter"
	"gocloud.dev/blob/gcsblob"
)

type GCSUploader struct {
	cloudblob.CloudBlob
}

const Scheme = gcsblob.Scheme

//It will use Application Default Credentials;
//if you have authenticated via gcloud auth login, it will use those credentials.
//See Application Default Credentials to learn about authentication alternatives, including using environment variables.
func (gcs *GCSUploader) UsageInfo() string {
	table, buf := usageprinter.NewUsageTable(Scheme)

	data := [][]string{
		[]string{"Upload to GCP Cloud Storage", "gs://my-bucket\nFor additional information see https://gocloud.dev/howto/blob/#gcs-ctor"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (gcs *GCSUploader) Scheme() string {
	return Scheme
}
