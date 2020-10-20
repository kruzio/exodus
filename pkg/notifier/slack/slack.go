package slack

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/kruzio/exodus/pkg/notifier/types"
	"github.com/kruzio/exodus/pkg/usageprinter"
	"github.com/nlopes/slack"
	"k8s.io/apimachinery/pkg/util/errors"
	log "k8s.io/klog"
	"net/url"
	"strings"
	"text/template"
	"time"
)

var slackTmpl string = `
`

var tmpl *template.Template = nil

func init() {
	funcs := sprig.TxtFuncMap()
	t, err := template.New("slack").Funcs(funcs).Parse(slackTmpl)

	if err != nil {
		panic(err)
	}
	tmpl = t
}

type SlackNotifier struct {
	targetURL string
}

const Scheme = "slack"

func (s *SlackNotifier) UsageInfo() string {
	table, buf := usageprinter.NewUsageTable(Scheme)

	data := [][]string{
		[]string{"Post notification to a slack channel", "slack://mychannel?apikey=<mykey>"},
		[]string{"", ""},
		[]string{"apikey=<mykey>", "Slack API token - xoxo-YOURTOKEN\nFor additional information see https://api.slack.com/apps\nNote that your app must join the destintation channel"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

var (
	colorMapper = map[string]string{

		"critical": "#b82d17",
		"crit":     "#b82d17",
		"danger":   "#b82d17",

		"high":  "#e64e36",
		"error": "#e64e36",

		"warn":    "#ffbf00",
		"warning": "#ffbf00",

		"medium": "#17a2b8",
		"med":    "#17a2b8",

		"info": "#002db3",
	}
)

func color(severity string) string {
	if color, exist := colorMapper[strings.ToLower(severity)]; exist {
		return color
	}

	return ""
}

func (s *SlackNotifier) create(title string, alert *types.Alert) *slack.Attachment {

	msg := slack.Attachment{
		AuthorLink: alert.GeneratorURL,
		AuthorIcon: "",
		Title:      fmt.Sprintf("%v [%v]", title, alert.Severity),
		TitleLink:  "",
		Pretext:    "",
		Text:       alert.Message,
		Footer:     fmt.Sprint("brought to you by kruz.io // exodus // started: ", alert.StartsAt.Format(time.RFC3339)),
		Color:      color(alert.Severity),
	}

	return &msg
}

func (s *SlackNotifier) notifierDo(apitoken string, channel string, title string, alert *types.Alert) (bool, error) {

	api := slack.New(apitoken)
	msg := s.create(title, alert)

	_, _, err := api.PostMessage(channel, slack.MsgOptionAttachments(*msg))
	if err != nil {
		log.Errorf("Failed to post notification to slack - %v", err)
		return true, err
	} else {
		log.V(5).Infof("Posted message to %v - '%v'", channel, msg)
	}

	return false, nil
}

func (s *SlackNotifier) Notifier(targetUrl string, alerts ...*types.Alert) (bool, error) {
	urlInfo, err := url.Parse(targetUrl)
	if err != nil {
		return false, err
	}

	if strings.ToLower(urlInfo.Scheme) != Scheme {
		return false, fmt.Errorf("Scheme %v not support by this uploader", urlInfo.Scheme)
	}
	apitoken := urlInfo.Query().Get("apikey")
	if apitoken == "" {
		return false, fmt.Errorf("Missing Slack API Key")
	}
	channel := urlInfo.Host

	title := fmt.Sprintf("Alert Notification")
	if customeTitle := urlInfo.Query().Get("title"); customeTitle != "" {
		title = customeTitle
	}

	log.V(5).Infof("Processing '%v' alert notifications to '%v'", len(alerts), channel)

	msgs := []slack.MsgOption{}
	for _, alert := range alerts {
		msg := s.create(title, alert)
		msgs = append(msgs, slack.MsgOptionAttachments(*msg))
	}

	api := slack.New(apitoken)

	errs := []error{}
	for _, msg := range msgs {
		_, _, err = api.PostMessage(channel, msg)
		if err != nil {
			//log.Errorf("Failed to post notification to slack - %v", err)
			errs = append(errs, err)
		}
	}

	log.V(5).Infof("Posted '%v' messages to %v - with '%v' errors", len(msgs), channel, len(errs))

	return false, errors.NewAggregate(errs)
}
