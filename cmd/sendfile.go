package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kruzio/exodus/pkg/sendfile"
	"github.com/spf13/cobra"
)

type SendFileOpts struct {
	Targets     []string
	Filename    string
	DstFilename string

	Watchdir string
	Forever  bool

	RetryTimeout  time.Duration
	RetryInterval time.Duration
}

func (o *SendFileOpts) Validate() error {
	if o.Filename != "" && o.Watchdir != "" {
		return fmt.Errorf("Either use --file <filename> or --watch-dir <dir>")
	}

	if len(o.Targets) == 0 {
		return fmt.Errorf("At least one target must be configured")
	}

	if o.Filename != "-" && o.Filename != "" {
		info, err := os.Stat(o.Filename)
		if os.IsNotExist(err) {
			return fmt.Errorf("Failed to stat '%v' - %v", o.Filename, err)
		}

		if info.IsDir() {
			return fmt.Errorf("File is a directoty")
		}
	}

	if o.Watchdir != "" {
		info, err := os.Stat(o.Watchdir)
		if os.IsNotExist(err) {
			return fmt.Errorf("Failed to stat '%v' - %v", o.Watchdir, err)
		}

		if !info.IsDir() {
			return fmt.Errorf("File %v is not a directoty", o.Watchdir)
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

# Send to Slack
echo myfilecontent | bin/exodus sendfile -f -  --target="slack://mychannel?apikey=xoxb-myslackapp-oauth-token&title=My File"

# Send files from a watch directory (one shot)
exodus sendfile  --target=webhook+http://localhost:8080/stuff?content-type=text --watch /tmp/exodus

# Send files from a watch directory (forever)
exodus sendfile  --target=webhook+http://localhost:8080/stuff?content-type=text --watch /tmp/exodus --watch-forever

Supported Targets
------------------
%s

`, sendfile.UsageInfo()),
		Hidden: false,
		RunE: func(c *cobra.Command, args []string) error {

			if err := opts.Validate(); err != nil {
				return err
			}

			if opts.Filename != "" {
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

				return sendfile.UploadTargetsWithRetry(
					opts.Targets,
					data,
					opts.DstFilename,
					opts.RetryInterval,
					opts.RetryTimeout)
			}

			//We are watching director
			err, doneChan := sendfile.UploadWatcher(opts.Watchdir,
				opts.Targets,
				opts.RetryInterval,
				opts.RetryTimeout, opts.Forever)

			if err != nil {
				return err
			}

			ok := <-doneChan
			close(doneChan)

			if !ok {
				return fmt.Errorf("Failed to watch & upload successfully")
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.Filename, "file", "f", "", "Input File - use '-' to read from stdin")
	flags.StringVar(&opts.DstFilename, "dest-filename", "", "Destination file name. Defaults to the filename provided by the '--file' CLI option")

	flags.StringVarP(&opts.Watchdir, "watch", "w", "", "Watch directory")
	flags.BoolVar(&opts.Forever, "watch-forever", false, "Watch the directory forever - if set to false, the watch will end after sending a file")

	flags.StringArrayVarP(&opts.Targets, "target", "t", strings.Split(os.Getenv("KRUZIO_EXODUS_SENDFILE_TARGETS"), ","), "One or more target URLs")
	flags.DurationVar(&opts.RetryTimeout, "retry-timeout", time.Second*30, "Retry timeout - set to 0 to skip retries")
	flags.DurationVar(&opts.RetryInterval, "retry-interval", time.Second*10, "The retry wait interval")

	return cmd
}
