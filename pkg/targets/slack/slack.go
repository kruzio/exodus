package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	log "k8s.io/klog"
	"net/url"
	"strings"
)

type SlackUploader struct {
	BucketUrl string
	DestName  string
}

const Scheme = "slack"

// This URL will open the container "my-container" using default
// credentials found in the environment variables
// AZURE_STORAGE_ACCOUNT plus at least one of AZURE_STORAGE_KEY
// and AZURE_STORAGE_SAS_TOKEN.
func (s *SlackUploader) UsageInfo() string {
	return "Upload to Slack Channel using the following URL scheme: slack://mychannel?apikey=<mykey>[&filetype=json&title=mymsgtitle] . For additional information see https://api.slack.com/apps"
}

func (s *SlackUploader) Export(data []byte) error {
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
	if ftype := urlInfo.Query().Get("filetype"); ftype != "" {
		fileType = ftype
	}

	title := fmt.Sprintf("Kruz IO - Exodus Upload - %v", s.DestName)
	if customeTitle := urlInfo.Query().Get("title"); customeTitle != "" {
		title = customeTitle
	}

	api := slack.New(apitoken)
	params := slack.FileUploadParameters{
		Title:          title,
		Filetype:       fileType,
		Content:        string(data),
		Channels:       []string{channel},
		InitialComment: "Kruz IO - Exodus Upload",
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

func (s *SlackUploader) SetUploadUrl(url string) error {
	s.BucketUrl = url

	return nil
}

func (s *SlackUploader) SetDestName(name string) error {
	s.DestName = name

	return nil
}
