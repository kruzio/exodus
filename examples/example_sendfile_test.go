package examples

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/sendfile"
	"github.com/kruzio/exodus/pkg/sendfile/webhook"
	"io/ioutil"
	log "k8s.io/klog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func Test_LocalFile_Example(t *testing.T) {

	dir, err := ioutil.TempDir("/tmp", "somedir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	data := "Programming today is a race between software engineers striving to " +
		"build bigger and better idiot-proof programs, and the Universe trying " +
		"to produce bigger and better idiots. So far, the Universe is winning."

	dstFile := "exodus-example"
	uploadUrl := fmt.Sprintf("file://%v", dir)

	uploader, err := sendfile.NewUploader(uploadUrl)
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

func Test_Webhook_Example(t *testing.T) {
	data := `{"XXX":""}`

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Incorrect method %v", r.Method)
		}

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
	defer ts.Close()

	tsUrlInfo, _ := url.Parse(ts.URL)

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v?skip-verify=true&x-headers=X-myheader:myval&token-bearer=1234", webhook.SchemeHTTPS, tsUrlInfo.Host)

	uploader, err := sendfile.NewUploader(uploadUrl)
	if err != nil {
		t.Fatal(err)
	}

	_ = uploader.SetDestName("helloworld")

	err = uploader.Export([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
}
