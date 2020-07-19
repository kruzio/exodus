package webhook

import (
	"context"
	"fmt"
	commoncfg "github.com/prometheus/common/config"
	"io/ioutil"
	log "k8s.io/klog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_Webhook_Https(t *testing.T) {
	data := `{"XXX":""}`
	AdditionalHeaders := map[string]string{
		"X-MyHeader": "SomeVal",
	}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Incorrect method %v", r.Method)
		}

		msg, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", data)

		if r.Header.Get("X-MyHeader") != "SomeVal" {
			t.Fatalf("missing header %+v", r.Header)
		}

		if string(msg) != data {
			t.Fatalf("contents = %q\nexpected = %q", string(msg), data)
		}
	}))
	defer ts.Close()

	tsUrlInfo, _ := url.Parse(ts.URL)

	//Let's create our client
	conf := &WebhookConfig{
		AdditionalHeaders: AdditionalHeaders,
		HTTPClientConfig: commoncfg.HTTPClientConfig{
			TLSConfig: commoncfg.TLSConfig{
				InsecureSkipVerify: true,
			},
		},
	}

	webhook, err := New(conf)
	if err != nil {
		log.Fatal(err)
	}

	res, err := webhook.PostJSON(context.Background(), fmt.Sprintf("https://%v", tsUrlInfo.Host), strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	webhook.Drain(res)
}
