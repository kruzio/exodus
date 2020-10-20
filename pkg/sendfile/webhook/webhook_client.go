package webhook

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kruzio/exodus/pkg/usageprinter"
	"github.com/olekukonko/tablewriter"
	"github.com/parnurzeal/gorequest"
	"k8s.io/apimachinery/pkg/util/errors"
	log "k8s.io/klog"
)

type WebhookClient struct {
	Url         string
	Method      string
	Header      http.Header
	ContentType string

	AdditionalHeaders map[string]string

	client *gorequest.SuperAgent
}

type WebhookHTTPClient struct {
	WebhookClient
}

type WebhookHTTPSClient struct {
	WebhookClient
}

const SchemeHTTPS = "webhook" //HTTPS
const SchemeHTTP = "webhook+http"

func usageInfo(scheme string) (*tablewriter.Table, *bytes.Buffer) {
	table, buf := usageprinter.NewUsageTable(scheme)

	data := [][]string{
		[]string{"Post to a webhook", fmt.Sprintf("For example: %v://myserver?x-headers=X-myheader:myval&token-bearer=1234", scheme)},
		[]string{"", ""},
		[]string{">> Authentication Options <<", ""},
		[]string{"token-bearer=<token>", "Support Authorization Bearer token based authentication"},
		[]string{"username=<username>", "Basic HTTP Authentication scheme"},
		[]string{"password=<password>", "Basic HTTP Authentication scheme"},
		[]string{"", ""},
		[]string{">> Additional Options <<", ""},
		[]string{"proxy-url=<proxy>", "The proxy URL the webhook client should connect to"},
		[]string{"content-type=<contentType>", "defaults to json and can be one of: json | text | xml | html | multipart"},
		[]string{"x-headers=k1:v1,k2:v2", "additional custom request headers "},
	}

	table.AppendBulk(data)

	return table, buf
}

func (w *WebhookHTTPSClient) UsageInfo() string {
	table, buf := usageInfo(SchemeHTTPS)
	data := [][]string{
		[]string{"", ""},
		[]string{">> TLS Options <<", ""},
		[]string{"skip-verify=true", "If one wished to allow connection to untrusted server"},
		[]string{"ca-file=<path-to-file>", "CA PEM file"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}
func (w *WebhookHTTPSClient) Scheme() string {
	return SchemeHTTPS
}

func (w *WebhookHTTPClient) UsageInfo() string {
	table, buf := usageInfo(SchemeHTTP)
	table.Render()

	return buf.String()
}
func (w *WebhookHTTPClient) Scheme() string {
	return SchemeHTTP
}

func (w *WebhookClient) Export(data []byte) error {
	urlInfo, err := url.Parse(w.Url)
	if err != nil {
		return err
	}

	if strings.ToLower(urlInfo.Scheme) != SchemeHTTP && strings.ToLower(urlInfo.Scheme) != SchemeHTTPS {
		return fmt.Errorf("Scheme %v not support by this uploader", urlInfo.Scheme)
	}

	webhookClient := gorequest.New()

	proxyUrl := urlInfo.Query().Get("proxy-url")
	if proxyUrl != "" {
		webhookClient = webhookClient.Proxy(proxyUrl)
		if len(webhookClient.Errors) > 0 {
			return errors.NewAggregate(webhookClient.Errors)
		}
	}

	//Authentication
	username := urlInfo.Query().Get("username")
	password := urlInfo.Query().Get("password")
	tokenBearer := urlInfo.Query().Get("token-bearer")
	if tokenBearer != "" && (username != "" || password != "") {
		return fmt.Errorf("Only one authentication scheme is allowed")
	}

	if username != "" || password != "" {
		webhookClient.SetBasicAuth(username, password)
	}

	//Custom Headers
	additionalHeaders := urlInfo.Query().Get("x-headers")
	appendHeaders := map[string]string{}
	if additionalHeaders != "" {
		headers := strings.Split(additionalHeaders, ",")
		for _, header := range headers {
			parts := strings.Split(header, ":")
			if len(parts) != 2 {
				return fmt.Errorf("malformed additonal-headers - %v at %v", header, additionalHeaders)
			}
			appendHeaders[parts[0]] = parts[1]
		}
	}

	//TLS Stuff
	if strings.ToLower(urlInfo.Scheme) == SchemeHTTPS {
		skipVerify := urlInfo.Query().Get("skip-verify")
		caFile := urlInfo.Query().Get("ca-file")
		webhookClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}

		// If a CA cert is provided then let's read it in so we can validate the
		// scrape target's certificate properly.
		if len(caFile) > 0 {
			data, err := ioutil.ReadFile(caFile)
			if err != nil {
				return fmt.Errorf("unable to load specified CA cert %s: %s", caFile, err)
			}

			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(data) {
				return fmt.Errorf("unable to use specified CA cert %s - %v", caFile, string(data))
			}
			webhookClient.Transport.TLSClientConfig.RootCAs = caCertPool
		}

		if skipVerify == "true" {
			webhookClient.Transport.TLSClientConfig.InsecureSkipVerify = true
		}
	}

	contentType := urlInfo.Query().Get("content-type")
	if contentType == "" {
		contentType = gorequest.TypeJSON
	}

	method := gorequest.POST
	scheme := "https"
	if strings.ToLower(urlInfo.Scheme) == SchemeHTTP {
		scheme = "http"
	}

	targetUrl := fmt.Sprintf("%v://%v%v", scheme, urlInfo.Host, urlInfo.Path)

	log.V(5).Infof("%v to %v", method, targetUrl)
	postReq := webhookClient.Post(targetUrl).Type(contentType)
	if len(appendHeaders) > 0 {
		for k, v := range appendHeaders {
			postReq = postReq.AppendHeader(k, v)
		}
	}
	if tokenBearer != "" {
		postReq = webhookClient.AppendHeader("Authorization", fmt.Sprintf("Bearer %s", tokenBearer))
	}

	res, msg, errs := postReq.Send(string(data)).End()
	if len(errs) > 0 {
		return fmt.Errorf("%v", errors.NewAggregate(errs))
	}

	if res.StatusCode/100 == http.StatusOK/100 {
		return nil
	}

	return fmt.Errorf("Status %v - %v", res.Status, msg)
}

func (w *WebhookClient) SetUploadUrl(_url string) error {
	urlInfo, err := url.Parse(_url)
	if err != nil {
		return err
	}

	if strings.ToLower(urlInfo.Scheme) != SchemeHTTP && strings.ToLower(urlInfo.Scheme) != SchemeHTTPS {
		return fmt.Errorf("Scheme %v not support by this uploader", urlInfo.Scheme)
	}

	w.Url = _url

	return nil
}

func (w *WebhookClient) SetDestName(name string) error {
	return nil
}
