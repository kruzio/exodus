package cmd

import (
	"fmt"
	"github.com/kruzio/exodus/pkg/notifier"
	"github.com/kruzio/exodus/pkg/notifier/reader"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type NotifyOpts struct {
	Targets  []string
	Filename string

	Watchdir string
	Forever  bool

	RetryTimeout  time.Duration
	RetryInterval time.Duration
}

func (o *NotifyOpts) Validate() error {
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
		if err := notifier.ValidateTargetUrl(t); err != nil {
			return err
		}
	}

	return nil
}

func NewCommandNotify() *cobra.Command {

	opts := &NotifyOpts{}

	// Support overrides
	cmd := &cobra.Command{
		Use:     "notify",
		Aliases: []string{"alert", "send-alert"},
		Short:   "Send Notification to one or more destinations",
		Long: fmt.Sprintf(`
Send notification alert to one or more destinations

# Send to Slack
echo "{}" | bin/exodus sendfile -f -  --target="slack://mychannel?apikey=xoxb-myslackapp-oauth-token&title=My File"


Supported Targets
------------------
%s

`, notifier.UsageInfo()),
		Hidden: false,
		RunE: func(c *cobra.Command, args []string) error {

			if err := opts.Validate(); err != nil {
				return err
			}

			if opts.Filename != "" {
				alerts, err := reader.LoadAlerts(opts.Filename)
				if err != nil {
					return err
				}

				return notifier.NotifyTargetsWithRetry(
					opts.Targets,
					opts.RetryInterval,
					opts.RetryTimeout,
					alerts...)
			}

			//We are watching director
			err, doneChan := notifier.NotifyWatcher(opts.Watchdir,
				opts.Targets,
				opts.RetryInterval,
				opts.RetryTimeout, opts.Forever)

			if err != nil {
				return err
			}

			ok := <-doneChan
			close(doneChan)

			if !ok {
				return fmt.Errorf("Failed to watch & notify successfully")
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.Filename, "file", "f", "", "Input File - use '-' to read from stdin")

	flags.StringVarP(&opts.Watchdir, "watch", "w", "", "Watch directory")
	flags.BoolVar(&opts.Forever, "watch-forever", false, "Watch the directory forever - if set to false, the watch will end after sending a file")

	flags.StringArrayVarP(&opts.Targets, "target", "t", strings.Split(os.Getenv("KRUZIO_EXODUS_SENDFILE_TARGETS"), ","), "One or more target URLs")
	flags.DurationVar(&opts.RetryTimeout, "retry-timeout", time.Second*30, "Retry timeout - set to 0 to skip retries")
	flags.DurationVar(&opts.RetryInterval, "retry-interval", time.Second*10, "The retry wait interval")

	return cmd
}
