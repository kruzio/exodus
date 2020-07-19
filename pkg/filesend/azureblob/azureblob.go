package azureblob

import (
	"github.com/kruzio/exodus/pkg/cloudblob"
	"github.com/kruzio/exodus/pkg/usageprinter"
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
	table, buf := usageprinter.NewUsageTable(Scheme)

	data := [][]string{
		[]string{"Upload to Azure Blob storage", "azblob://my-container\nFor additional information see https://gocloud.dev/howto/blob/#azure"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (azblob *AzureBlobUploader) Scheme() string {
	return Scheme
}
