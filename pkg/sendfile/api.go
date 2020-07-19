package sendfile

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"

	"github.com/kruzio/exodus/pkg/sendfile/azureblob"
	"github.com/kruzio/exodus/pkg/sendfile/file"
	"github.com/kruzio/exodus/pkg/sendfile/gcs"
	"github.com/kruzio/exodus/pkg/sendfile/s3"
	"github.com/kruzio/exodus/pkg/sendfile/slack"
	"github.com/kruzio/exodus/pkg/sendfile/smtp"
	"github.com/kruzio/exodus/pkg/sendfile/webhook"
)

type Target interface {
	Export(data []byte) error
	SetUploadUrl(url string) error
	SetDestName(name string) error

	UsageInfo() string
	Scheme() string
}

type CreateUploader func() Target

var targets = map[string]CreateUploader{
	s3.Scheme:        func() Target { return &s3.S3Uploader{} },
	azureblob.Scheme: func() Target { return &azureblob.AzureBlobUploader{} },
	gcs.Scheme:       func() Target { return &gcs.GCSUploader{} },
	slack.Scheme:     func() Target { return &slack.SlackFileUploader{} },
	file.Scheme:      func() Target { return &file.LocalFileUploader{} },

	webhook.SchemeHTTP:  func() Target { return &webhook.WebhookHTTPClient{} },
	webhook.SchemeHTTPS: func() Target { return &webhook.WebhookHTTPSClient{} },
	smtp.Scheme:         func() Target { return &smtp.Email{} },
}

func ValidateTargetUrl(targetUrl string) error {
	urlInfo, err := url.Parse(targetUrl)

	if err != nil {
		return err
	}

	_, exist := targets[strings.ToLower(urlInfo.Scheme)]
	if !exist {
		return fmt.Errorf("The URL scheme %v is not supported", urlInfo.Scheme)
	}

	return nil
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

func UploadWithRetry(target Target, data []byte, retryInterval time.Duration, retryTimeout time.Duration) error {
	retryCount := 0

	errPolling := wait.PollImmediate(retryInterval, retryTimeout, func() (done bool, failToRunError error) {
		err := target.Export(data)
		if err != nil {
			klog.V(5).Infof("Failed to to upload after %v retries. - %v", retryCount, err)
			retryCount++
			return false, nil
		}

		return true, nil
	})

	if errPolling != nil {
		klog.V(5).Infof("Failed to upload. Polling failed. %v", errPolling)
		return errPolling
	}

	return nil
}
