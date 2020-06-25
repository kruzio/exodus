package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version = ""
	Commit  = ""
)

func NewCommandVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print exodus version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version: " + Version + "\nCommit: " + Commit)
		},
	}
}
