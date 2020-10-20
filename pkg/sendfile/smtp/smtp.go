package smtp

import (
	"crypto/tls"
	"fmt"
	"github.com/kruzio/exodus/pkg/usageprinter"
	"net/smtp"
	"net/textproto"
	"net/url"
	"strings"

	"github.com/jordan-wright/email"
	log "k8s.io/klog"
)

type Email struct {
	Url            string
	AttachmentName string
}

const Scheme = "smtp"

func (s *Email) UsageInfo() string {
	table, buf := usageprinter.NewUsageTable(Scheme)

	data := [][]string{
		[]string{"Send file via email (smtp)", "smtp://smtpserver?to=<email>&from=<email>&username=myuser&password=mypass"},
		[]string{"", ""},
		[]string{"to=<target>[,<target>]", "the destination email address(es) - required"},
		[]string{"from=<from-email>", "From email address - required"},
		[]string{"", ""},
		[]string{"username=<username>", "The smtp server authentication information - optional"},
		[]string{"password=<password>", "The smtp server authentication information - optional"},
		[]string{"", ""},
		[]string{"subject=<subject>", "The Subject line of the email message"},
		[]string{"skip-verify=true", "Skip SMTP server TLS verification - (not recommended)"},
		[]string{"mime-type=application/json", "MIME type of the mail attachment - default is application/json"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (s *Email) Scheme() string {
	return Scheme
}

func (s *Email) Export(data []byte) error {
	urlInfo, err := url.Parse(s.Url)
	if err != nil {
		return err
	}

	if strings.ToLower(urlInfo.Scheme) != Scheme {
		return fmt.Errorf("Scheme %v not support by this uploader", urlInfo.Scheme)
	}

	sendTo := []string{}
	if toList := urlInfo.Query().Get("to"); toList != "" {
		to := strings.Split(toList, ",")
		sendTo = to
	}
	if len(sendTo) == 0 {
		return fmt.Errorf("failed to find destination emails in '%v'", urlInfo.Query().Get("to"))
	}

	subject := fmt.Sprintf("Kruz IO - Exodus File Notification - %v", s.AttachmentName)
	if customeTitle := urlInfo.Query().Get("subject"); customeTitle != "" {
		subject = customeTitle
	}

	from := "Kruz IO <do-not-reply@kruz.io>"
	if s := urlInfo.Query().Get("from"); s != "" {
		from = s
	}

	mimetype := "application/json"
	if s := urlInfo.Query().Get("mime-type"); s != "" {
		mimetype = s
	}

	username := urlInfo.Query().Get("username")
	password := urlInfo.Query().Get("password")
	var auth smtp.Auth = nil
	if username != "" || password != "" {
		auth = smtp.PlainAuth("", username, password, urlInfo.Host)
	}

	skipVerify := false
	if skip := urlInfo.Query().Get("skip-verify"); skip == "true" {
		skipVerify = true
	}

	e := &email.Email{
		To:      sendTo,
		From:    from,
		Subject: subject,
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}
	_, err = e.Attach(strings.NewReader(string(data)), s.AttachmentName, mimetype)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{}
	port := urlInfo.Port()
	if port == "465" {

		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName = urlInfo.Hostname()
		}

		tlsConfig.InsecureSkipVerify = skipVerify
		err = e.SendWithTLS(urlInfo.Host, auth, tlsConfig)

	} else {
		err = e.Send(urlInfo.Host, auth)
	}

	if err != nil {
		log.Errorf("Failed to email file - %v", urlInfo.Host)
		return err
	}

	log.V(5).Infof("File sent to %v as %v", urlInfo.Host, s.AttachmentName)
	return nil
}

func (s *Email) SetUploadUrl(url string) error {
	s.Url = url

	return nil
}

func (s *Email) SetDestName(name string) error {
	s.AttachmentName = name

	return nil
}
