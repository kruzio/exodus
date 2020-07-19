package cmd

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/sendfile"
	"github.com/spf13/cobra"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/errors"
	klog "k8s.io/klog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SendFileOpts struct {
	Targets     []string
	Filename    string
	DstFilename string

	RetryTimeout  time.Duration
	RetryInterval time.Duration
}

func (o *SendFileOpts) Validate() error {
	if len(o.Targets) == 0 {
		return fmt.Errorf("At least one target must be configured")
	}

	if o.Filename != "-" {
		info, err := os.Stat(o.Filename)
		if os.IsNotExist(err) {

			return err
		}

		if info.IsDir() {
			return fmt.Errorf("File is a directoty")
		}
	}

	for _, t := range o.Targets {
		if err := sendfile.ValidateTargetUrl(t); err != nil {
			return err
		}
	}

	return nil
}

func NewCommandSendFile() *cobra.Command {

	opts := &SendFileOpts{}

	// Support overrides
	cmd := &cobra.Command{
		Use:     "sendfile",
		Aliases: []string{"filesender", "file-send", "send-file"},
		Short:   "Send File to one or more destinations",
		Long: fmt.Sprintf(`
Send File to one or more destinations

#Render Locally
echo myfilecontent | bin/exodus sendfile -f -  --target="slack://mychannel?apikey=xoxb-myslackapp-oauth-token&title=My File"

Supported Targets
------------------
%s

`, sendfile.UsageInfo()),
		Hidden: false,
		RunE: func(c *cobra.Command, args []string) error {

			if err := opts.Validate(); err != nil {
				return err
			}

			data := []byte{}
			var err error
			if opts.Filename == "-" {
				data, err = ioutil.ReadAll(os.Stdin)
			} else {
				data, err = ioutil.ReadFile(opts.Filename)
			}
			if err != nil {
				return err
			}

			if opts.DstFilename == "" {
				_, opts.DstFilename = filepath.Split(opts.Filename)
				if opts.DstFilename == "-" {
					opts.DstFilename = "exodus-stdin"
				}
			}

			wg := &sync.WaitGroup{}
			errChan := make(chan error, len(opts.Targets))
			for i, target := range opts.Targets {
				wg.Add(1)
				go func(target string, data []byte, errChan chan<- error) {
					defer wg.Done()

					uploader, err := sendfile.NewUploader(target)
					if err != nil {
						klog.Errorf("%v", err)
						errChan <- err
						return
					}

					klog.V(3).Infof("Processing target #%v - %v", i+1, uploader.Scheme())

					_ = uploader.SetDestName(opts.DstFilename)

					if opts.RetryTimeout > 0 && opts.RetryInterval > 0 {
						err = sendfile.UploadWithRetry(uploader, data, opts.RetryInterval, opts.RetryTimeout)
					} else {
						err = uploader.Export(data)
					}

					if err != nil {
						errChan <- err
						klog.Errorf("%v", err)
						return
					}
				}(target, data, errChan)
			}

			klog.V(3).Infof("Waiting for all targets to complete")
			wg.Wait()
			klog.V(3).Infof("Targets Completed")
			close(errChan)
			errs := []error{}
			for err := range errChan {
				errs = append(errs, err)
			}

			return errors.NewAggregate(errs)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.Filename, "file", "f", "", "Input File - use '-' to read from stdin")
	flags.StringVar(&opts.DstFilename, "dest-filename", "", "Destination file name. Defaults to the filename provided by the '--file' CLI option")
	flags.StringArrayVarP(&opts.Targets, "target", "t", []string{}, "One or more target URLs")
	flags.DurationVar(&opts.RetryTimeout, "retry-timeout", time.Second*30, "Retry timeout - set to 0 to skip retries")
	flags.DurationVar(&opts.RetryInterval, "retry-interval", time.Second*10, "The retry wait interval")

	return cmd
}
