package gcs

import (
	"github.com/kruzio/exodus/pkg/cloudblob"
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
	return "Upload to GCP Cloud Storage using the following URL scheme: gs://my-bucket. For additional information see https://gocloud.dev/howto/blob/#gcs-ctor"
}
