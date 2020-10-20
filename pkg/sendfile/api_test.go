package sendfile

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/sendfile/webhook"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	log "k8s.io/klog"
)

func Test_LocalFileUploader(t *testing.T) {

	dir, err := ioutil.TempDir("/tmp", "somedir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	data := "Programming today is a race between software engineers striving to " +
		"build bigger and better idiot-proof programs, and the Universe trying " +
		"to produce bigger and better idiots. So far, the Universe is winning."

	dstFile := "exodus-test"
	uploadUrl := fmt.Sprintf("file://%v", dir)

	uploader, err := NewUploader(uploadUrl)
	if err != nil {
		t.Fatal(err)
	}

	_ = uploader.SetDestName(dstFile)

	err = uploader.Export([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadFile(filepath.Join(dir, dstFile))
	if err != nil {
		t.Fatalf("ReadFile %s: %v", dstFile, err)
	}

	if string(contents) != data {
		t.Fatalf("contents = %q\nexpected = %q", string(contents), data)
	}

	// cleanup
	_ = os.Remove(filepath.Join(dir, dstFile)) // ignore error
}

func Test_WebhookRetryWithFailures(t *testing.T) {
	data := `{"XXX":""}`
	block := true
	lock := sync.Mutex{}
	blockPtr := &block

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		apiBlock := *blockPtr
		lock.Unlock()

		if apiBlock == true {
			r.Body.Close()
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Println("Webhook called - blocking")
			return
		}

		fmt.Println("Webhook called - processing")

		msg, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		if string(msg) != data {
			t.Fatalf("contents = %q\nexpected = %q", string(msg), data)
		}

		if r.Header.Get("X-myheader") != "myval" {
			t.Fatalf("headers mismatch - %v", r.Header)
		}

		if r.Header.Get("Authorization") != "Bearer 1234" {
			t.Fatalf("auth token bearer - mismatch - %v", r.Header)
		}
	}))
	completed := make(chan bool, 1)
	defer ts.Close()
	defer close(completed)

	tsUrlInfo, _ := url.Parse(ts.URL)

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v?skip-verify=true&x-headers=X-myheader:myval&token-bearer=1234", webhook.SchemeHTTPS, tsUrlInfo.Host)

	go func(completionChan chan<- bool) {
		fmt.Println("target ", uploadUrl)
		err := UploadTargetsWithRetry(
			[]string{uploadUrl},
			[]byte(data),
			"testfile",
			3*time.Second,
			10*time.Second)

		if err == nil {
			completionChan <- true
		} else {
			completionChan <- false
		}

	}(completed)

	go func() {
		fmt.Printf("Suspend Webhook Server availability\n")
		time.Sleep(4 * time.Second)
		fmt.Printf("Make Webhook Server available\n")
		lock.Lock()
		*blockPtr = false
		lock.Unlock()
	}()

	sentOK := <-completed

	if !sentOK {
		t.Fatalf("failed to send")
	}
}
