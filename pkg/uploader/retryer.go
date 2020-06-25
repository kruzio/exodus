package uploader

import (
	"github.com/kruzio/exodus/pkg/targets"
	"k8s.io/apimachinery/pkg/util/wait"
	log "k8s.io/klog"
	"time"
)

func UploadWithRetry(target targets.Target, data []byte, _retryInterval int64, _retryTimeout int64) error {
	retryCount := 0
	retryInterval := time.Second * time.Duration(_retryInterval)
	retryTimeout := time.Second * time.Duration(_retryTimeout)

	errPolling := wait.PollImmediate(retryInterval, retryTimeout, func() (done bool, failToRunError error) {
		err := target.Export(data)
		if err != nil {
			log.Warningf("Failed to to upload after %v retries. - %v", retryCount, err)
			retryCount++
			return false, nil
		}

		return true, nil
	})

	if errPolling != nil {
		log.Errorf("Failed to upload. Polling failed. %v", errPolling)
		return errPolling
	}

	return nil
}
