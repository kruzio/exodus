package sendfile

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/sendfile/webhook"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_UploadWatcher(t *testing.T) {

	completed := make(chan bool, 1)
	defer close(completed)

	data := `{"XXX":""}`

	dir, err := ioutil.TempDir("/tmp", "export-box")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		completed <- true
	}))

	defer ts.Close()

	tsUrlInfo, _ := url.Parse(ts.URL)

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v?skip-verify=true&x-headers=X-myheader:myval&token-bearer=1234", webhook.SchemeHTTPS, tsUrlInfo.Host)

	err, doneChan := UploadWatcher(dir,
		[]string{uploadUrl},
		3*time.Second,
		10*time.Second, false)

	if err != nil {
		t.Fatal(err)
	}

	//time.Sleep(2*time.Second)

	fmt.Println("Write file to:", filepath.Join(dir, "testfile.json"))
	err = ioutil.WriteFile(filepath.Join(dir, "testfile.json"), []byte(data), 0600)
	if err != nil {
		t.Fatal(err)
	}

	ok := <-doneChan
	close(doneChan)

	if !ok {
		t.Fatalf("failed to send")
	}

	ok = <-completed

	if !ok {
		t.Fatalf("failed to send")
	}
}

func Test_UploadWatcher_Failure(t *testing.T) {
	data := `{"XXX":""}`

	dir, err := ioutil.TempDir("/tmp", "export-box")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://localhost:66660?skip-verify=true", webhook.SchemeHTTPS)

	err, doneChan := UploadWatcher(dir,
		[]string{uploadUrl},
		3*time.Second,
		5*time.Second, false)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Write file to:", filepath.Join(dir, "testfile.json"))
	err = ioutil.WriteFile(filepath.Join(dir, "testfile.json"), []byte(data), 0600)
	if err != nil {
		t.Fatal(err)
	}

	ok := <-doneChan
	close(doneChan)

	if ok {
		t.Fatalf("expected send to fail")
	}
}
