package targets

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/targets/azureblob"
	"github.com/kruzio/exodus/pkg/targets/file"
	"github.com/kruzio/exodus/pkg/targets/gcs"
	"github.com/kruzio/exodus/pkg/targets/s3"
	"github.com/kruzio/exodus/pkg/targets/slack"
	"net/url"
	"strings"
)

type Target interface {
	Export(data []byte) error
	SetUploadUrl(url string) error
	SetDestName(name string) error

	UsageInfo() string
}

type CreateUploader func() Target

var targets = map[string]CreateUploader{
	s3.Scheme:        func() Target { return &s3.S3Uploader{} },
	azureblob.Scheme: func() Target { return &azureblob.AzureBlobUploader{} },
	gcs.Scheme:       func() Target { return &gcs.GCSUploader{} },
	slack.Scheme:     func() Target { return &slack.SlackUploader{} },

	file.Scheme: func() Target { return &file.LocalFileUploader{} },
}

func NewUploader(targetUrl string) (Target, error) {
	urlInfo, err := url.Parse(targetUrl)

	if err != nil {
		return nil, err
	}

	createUploader, exist := targets[strings.ToLower(urlInfo.Scheme)]
	if !exist {
		return nil, fmt.Errorf("The URL scheme %v is not supported", urlInfo.Scheme)
	}

	uploader := createUploader()
	if err := uploader.SetUploadUrl(targetUrl); err != nil {
		return nil, err
	}

	return uploader, nil
}
