package main

import (
	"bytes"
	goflag "flag"
	"fmt"
	"k8s.io/klog"
	"os"

	"github.com/kruzio/exodus/cmd"
	"github.com/spf13/cobra"
)

func ExodusCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "exodus",
		Short: "exodus",
		Long:  `exodus`,
	}

	var genBashCompletionCmd = &cobra.Command{
		Use:   "bash-completion",
		Short: "Generate bash completion. source < (exodus bash-completion)",
		Long:  "Generate bash completion. source < (exodus bash-completion)",
		Run: func(cmd *cobra.Command, args []string) {
			out := new(bytes.Buffer)
			_ = rootCmd.GenBashCompletion(out)
			println(out.String())
		},
	}

	cmds := []*cobra.Command{
		cmd.NewCommandVersion(),
		cmd.NewCommandSendFile(),
		cmd.NewCommandNotify(),
		genBashCompletionCmd,
	}

	flags := rootCmd.PersistentFlags()

	klog.InitFlags(nil)
	flags.AddGoFlagSet(goflag.CommandLine)

	// Hide all klog flags except for -v
	goflag.CommandLine.VisitAll(func(f *goflag.Flag) {
		if f.Name != "v" {
			flags.Lookup(f.Name).Hidden = true
		}
	})

	rootCmd.AddCommand(cmds...)
	return rootCmd
}

func main() {

	rootCmd := ExodusCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
