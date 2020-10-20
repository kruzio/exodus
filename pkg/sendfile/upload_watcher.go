package sendfile

import (
	"io/ioutil"
	"time"

	"github.com/fsnotify/fsnotify"
	"k8s.io/klog"
)

func UploadWatcher(watchdir string, targets []string, retryInterval time.Duration, retryTimeout time.Duration, forever bool) (error, chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err, nil
	}

	err = watcher.Add(watchdir)
	if err != nil {
		return err, nil
	}
	doneChan := make(chan bool, 1)

	go func(watcher *fsnotify.Watcher) {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					klog.Warningf("event channel closed: %v", event)
					doneChan <- false
					return
				}

				klog.V(5).Infof("event: %v", event)

				if event.Op&fsnotify.Create != fsnotify.Create {
					klog.V(5).Infof("none create event - %v", event.Name)
					continue
				}

				data, err := ioutil.ReadFile(event.Name)
				if err != nil {
					klog.Errorf("Failed to read %v", event.Name)
					continue
				}

				err = UploadTargetsWithRetry(targets, data, event.Name, retryInterval, retryTimeout)

				if !forever {
					klog.V(5).Infof("Processing Single File event - %v - err=%v", event.Name, err)
					if err == nil {
						doneChan <- true
					} else {
						doneChan <- false
					}

					return
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					klog.Warningf("error channel closed")
				} else {
					klog.Errorf("error: %v", err)
				}

				if !forever {
					doneChan <- false
					return
				}
			}
		}
	}(watcher)

	return nil, doneChan
}
