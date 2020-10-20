package webhook

import (
	"fmt"
	"io/ioutil"
	log "k8s.io/klog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_HTTPS_WebHookPostWithHeaders(t *testing.T) {
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

		fmt.Printf("%s", data)

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
	uploadUrl := fmt.Sprintf("%v://%v?skip-verify=true&x-headers=X-myheader:myval&token-bearer=1234", SchemeHTTPS, tsUrlInfo.Host)
	webhookClient := &WebhookClient{
		Url: uploadUrl,
	}
	_ = webhookClient.SetDestName("stuff")

	err := webhookClient.Export([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
}

func Test_HTTP_WebHookPost(t *testing.T) {
	data := `{"XXX":""}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Incorrect method %v", r.Method)
		}

		msg, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		if string(msg) != data {
			t.Fatalf("\ncontents='%v'\nexpected= '%v'\n%+v", string(msg), data, r)
		}
	}))
	defer ts.Close()

	tsUrlInfo, _ := url.Parse(ts.URL)

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v", SchemeHTTP, tsUrlInfo.Host)
	webhookClient := &WebhookClient{
		Url: uploadUrl,
	}
	_ = webhookClient.SetDestName("stuff")

	err := webhookClient.Export([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
}

func Test_HTTP_Post_Text(t *testing.T) {
	data := `XXX`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Incorrect method %v", r.Method)
		}

		msg, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		if string(msg) != data {
			t.Fatalf("\ncontents='%v'\nexpected= '%v'\n%+v", string(msg), data, r)
		}
	}))
	defer ts.Close()

	tsUrlInfo, _ := url.Parse(ts.URL)

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v?content-type=text", SchemeHTTP, tsUrlInfo.Host)
	webhookClient := &WebhookClient{
		Url: uploadUrl,
	}
	_ = webhookClient.SetDestName("stuff")

	err := webhookClient.Export([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
}
