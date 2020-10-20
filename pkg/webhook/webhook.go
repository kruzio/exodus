package webhook

import (
	"context"
	commoncfg "github.com/prometheus/common/config"
	"io"
	"io/ioutil"
	"net/http"
)

type WebhookConfig struct {
	commoncfg.HTTPClientConfig

	AdditionalHeaders map[string]string `yaml:"additional_headers,omitempty"`
}

// Notifier implements a Notifier for generic webhooks.
type Webhook struct {
	conf   *WebhookConfig
	client *http.Client
}

// New returns a new Webhook.
func New(conf *WebhookConfig) (*Webhook, error) {
	client, err := commoncfg.NewClientFromConfig(conf.HTTPClientConfig, "webhook")
	if err != nil {
		return nil, err
	}

	return &Webhook{
		conf: conf,
		// Webhooks are assumed to respond with 2xx response codes on a successful
		// request and 5xx response codes are assumed to be recoverable.
		client: client,
	}, nil
}

// RedactURL removes the URL part from an error of *url.Error type.
//func RedactURL(err error) error {
//	e, ok := err.(*url.Error)
//	if !ok {
//		return err
//	}
//	e.URL = "<redacted>"
//	return e
//}

// PostJSON sends a POST request with JSON payload to the given URL.
func (w *Webhook) PostJSON(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	return w.post(ctx, w.client, url, "application/json", body)
}

// PostText sends a POST request with text payload to the given URL.
func (w *Webhook) PostText(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	return w.post(ctx, w.client, url, "text/plain", body)
}

func (w *Webhook) post(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	if len(w.conf.AdditionalHeaders) > 0 {
		for k, v := range w.conf.AdditionalHeaders {
			req.Header.Set(k, v)
		}
	}
	return client.Do(req.WithContext(ctx))
}

// Drain consumes and closes the response's body to make sure that the
// HTTP client can reuse existing connections.
func (w *Webhook) Drain(r *http.Response) {
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
}
