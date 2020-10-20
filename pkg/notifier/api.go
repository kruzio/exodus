package notifier

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/notifier/webhook"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/url"
	"strings"
	"sync"
	"text/template"
	"time"

	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog"

	"github.com/kruzio/exodus/pkg/notifier/slack"
	"github.com/kruzio/exodus/pkg/notifier/types"
)

// Notifier notifies about alerts under constraints of the given context. It
// returns an error if unsuccessful and a flag whether the error is
// recoverable. This information is useful for a retry logic.
type Target interface {
	Notifier(targetUrl string, alerts ...*types.Alert) (bool, error)
	UsageInfo() string
}

type CreateUploader func(tmpl *template.Template) Target

var targets = map[string]CreateUploader{
	slack.Scheme:        func(tmpl *template.Template) Target { return &slack.SlackNotifier{} },
	webhook.SchemeHTTP:  func(tmpl *template.Template) Target { return &webhook.WebhookHTTPClient{} },
	webhook.SchemeHTTPS: func(tmpl *template.Template) Target { return &webhook.WebhookHTTPSClient{} },
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

func NewNotifier(targetUrl string) (Target, error) {
	urlInfo, err := url.Parse(targetUrl)

	if err != nil {
		klog.V(7).Infof("Failed to to parse - %v - %v", targetUrl, err)
		return nil, err
	}

	createUploader, exist := targets[strings.ToLower(urlInfo.Scheme)]
	if !exist {
		return nil, fmt.Errorf("The URL scheme %v is not supported", urlInfo.Scheme)
	}

	uploader := createUploader(nil)

	return uploader, nil
}

func NotifyWithRetry(target Target, targetUrl string, retryInterval time.Duration, retryTimeout time.Duration, alert ...*types.Alert) error {
	retryCount := 0

	errPolling := wait.PollImmediate(retryInterval, retryTimeout, func() (done bool, failToRunError error) {
		_, err := target.Notifier(targetUrl, alert...)
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

func NotifyTargetsWithRetry(targets []string, retryInterval time.Duration, retryTimeout time.Duration, alerts ...*types.Alert) error {
	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(targets))
	for i, target := range targets {
		wg.Add(1)
		go func(target string, errChan chan<- error, alerts ...*types.Alert) {
			defer wg.Done()

			notifier, err := NewNotifier(target)
			if err != nil {
				klog.Errorf("%v", err)
				errChan <- err
				return
			}

			klog.V(3).Infof("Processing target #%v", i+1)

			if retryTimeout > 0 && retryInterval > 0 {
				err = NotifyWithRetry(notifier, target, retryInterval, retryTimeout, alerts...)
			} else {
				_, err = notifier.Notifier(target, alerts...)
			}

			if err != nil {
				errChan <- err
				klog.Errorf("%v", err)
				return
			}
		}(target, errChan, alerts...)
	}

	klog.V(3).Infof("Waiting for all targets to complete")
	wg.Wait()
	klog.V(3).Infof("Targets Completed")
	close(errChan)

	errs := []error{}
	for err := range errChan {
		errs = append(errs, err)
	}

	return errors.NewAggregate(errs)
}
