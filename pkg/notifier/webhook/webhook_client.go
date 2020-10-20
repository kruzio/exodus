package webhook

import (
	"encoding/json"
	"github.com/kruzio/exodus/pkg/notifier/types"
	"github.com/kruzio/exodus/pkg/sendfile"
	"github.com/kruzio/exodus/pkg/sendfile/webhook"
)

type WebhookHTTPClient struct {
	webhook.WebhookHTTPClient
}

type WebhookHTTPSClient struct {
	webhook.WebhookHTTPSClient
}

const SchemeHTTPS = "webhook" //HTTPS
const SchemeHTTP = "webhook+http"

func (w *WebhookHTTPClient) Notifier(targetUrl string, alerts ...*types.Alert) (bool, error) {
	return notifier(w, targetUrl, alerts...)
}

func (w *WebhookHTTPSClient) Notifier(targetUrl string, alerts ...*types.Alert) (bool, error) {
	return notifier(w, targetUrl, alerts...)
}

func notifier(target sendfile.Target, targetUrl string, alerts ...*types.Alert) (bool, error) {
	data, err := json.Marshal(&alerts)
	if err != nil {
		return false, err
	}

	err = target.SetUploadUrl(targetUrl)
	if err != nil {
		return false, err
	}

	err = target.Export(data)
	if err != nil {
		return false, err
	}

	return false, nil
}
