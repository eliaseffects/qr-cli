package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion(cmd)
		},
	}
)

func printVersion(cmd *cobra.Command) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "qr-cli %s\ncommit: %s\nbuilt:  %s\n", Version, Commit, Date)
}
