package notifier

import (
	"encoding/json"
	"fmt"
	"github.com/kruzio/exodus/pkg/notifier/types"

	"io/ioutil"
	log "k8s.io/klog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

var alert_1 = types.Alert{
	Severity:     "high",
	Message:      "alert-1",
	StartsAt:     time.Now(),
	GeneratorURL: "unittest",
}

var alert_2 = types.Alert{
	Severity:     "low",
	Message:      "alert-2",
	StartsAt:     time.Now(),
	GeneratorURL: "unittest",
}

func makeAlerts() []*types.Alert {
	alerts := []*types.Alert{
		&alert_1,
		&alert_2,
	}

	return alerts
}

func Test_WebhookRetryWithFailures(t *testing.T) {
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

		if r.Header.Get("X-myheader") != "myval" {
			t.Fatalf("headers mismatch - %v", r.Header)
		}

		if r.Header.Get("Authorization") != "Bearer 1234" {
			t.Fatalf("auth token bearer - mismatch - %v", r.Header)
		}

		fmt.Println(string(msg))
		alerts := []*types.Alert{}
		err = json.Unmarshal(msg, &alerts)
		if err != nil {
			log.Fatal(err)
		}
	}))
	completed := make(chan bool, 1)
	defer ts.Close()
	defer close(completed)

	tsUrlInfo, _ := url.Parse(ts.URL)

	//Let's create our client
	uploadUrl := fmt.Sprintf("webhook://%v?skip-verify=true&x-headers=X-myheader:myval&token-bearer=1234", tsUrlInfo.Host)

	go func(completionChan chan<- bool) {
		fmt.Println("target ", uploadUrl)
		err := NotifyTargetsWithRetry(
			[]string{uploadUrl},
			3*time.Second,
			10*time.Second,
			makeAlerts()...)

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

func Test_Slack(t *testing.T) {

	//Let's create our client
	uploadUrl := "slack://channel?apikey=xoxo-somekey"

	fmt.Println("target ", uploadUrl)
	err := NotifyTargetsWithRetry(
		[]string{uploadUrl},
		0,
		0,
		makeAlerts()...)

	if err == nil {
		t.Fatalf("expected to fail")
	}

}
