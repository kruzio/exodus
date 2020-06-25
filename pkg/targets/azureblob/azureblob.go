package azureblob

import (
	"github.com/kruzio/exodus/pkg/cloudblob"
)

type AzureBlobUploader struct {
	cloudblob.CloudBlob
}

const Scheme = "azblob"

// This URL will open the container "my-container" using default
// credentials found in the environment variables
// AZURE_STORAGE_ACCOUNT plus at least one of AZURE_STORAGE_KEY
// and AZURE_STORAGE_SAS_TOKEN.
func (azblob *AzureBlobUploader) UsageInfo() string {
	return "Upload to AWS S3 bucket using the following URL scheme: azblob://my-container. For additional information see https://gocloud.dev/howto/blob/#azure"
}
