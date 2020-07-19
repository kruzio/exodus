package slack

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/usageprinter"
	log "k8s.io/klog"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
)

type SlackFileUploader struct {
	BucketUrl string
	DestName  string
}

const Scheme = "slack"

func (s *SlackFileUploader) UsageInfo() string {
	table, buf := usageprinter.NewUsageTable(Scheme)

	data := [][]string{
		[]string{"Post file to a slack channel", "slack://mychannel?apikey=<mykey>[&file-type=json&title=mymsgtitle]"},
		[]string{"", ""},
		[]string{"apikey=<mykey>", "Slack API token - xoxo-YOURTOKEN\nFor additional information see https://api.slack.com/apps\nNote that your app must join the destintation channel"},
		[]string{"file-type=json", "The content type"},
		[]string{"title=mymsgtitle", "The notification title"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (s *SlackFileUploader) Scheme() string {
	return Scheme
}

func (s *SlackFileUploader) Export(data []byte) error {
	urlInfo, err := url.Parse(s.BucketUrl)
	if err != nil {
		return err
	}

	if strings.ToLower(urlInfo.Scheme) != Scheme {
		return fmt.Errorf("Scheme %v not support by this uploader", urlInfo.Scheme)
	}
	apitoken := urlInfo.Query().Get("apikey")
	if apitoken == "" {
		return fmt.Errorf("Missing Slack API Key")
	}

	channel := urlInfo.Host
	fileType := "text"
	if ftype := urlInfo.Query().Get("file-type"); ftype != "" {
		fileType = ftype
	}

	title := fmt.Sprintf("Kruz IO - Exodus Upload - %v", s.DestName)
	if customeTitle := urlInfo.Query().Get("title"); customeTitle != "" {
		title = customeTitle
	}

	api := slack.New(apitoken)
	params := slack.FileUploadParameters{
		Title:    title,
		Filetype: fileType,
		Content:  string(data),
		Channels: []string{channel},
		//InitialComment: "Kruz IO - Exodus Upload",
	}

	file, err := api.UploadFile(params)
	if err != nil {
		log.Errorf("Failed to upload file to slack yaml output - %v", err)
		return err
	} else {
		log.V(5).Infof("File sent to %v as %v with ID - '%v'", channel, s.DestName, file.ID)
	}

	return nil
}

func (s *SlackFileUploader) SetUploadUrl(url string) error {
	s.BucketUrl = url

	return nil
}

func (s *SlackFileUploader) SetDestName(name string) error {
	s.DestName = name

	return nil
}
