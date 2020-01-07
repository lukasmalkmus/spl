package main

import (
	"runtime"

	"github.com/spf13/cobra"

	"github.com/lukasmalkmus/spl/pkg/version"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version and build details",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("spl v%s built for %s/%s\n\n",
			version.Release(), runtime.GOOS, runtime.GOARCH)
		cmd.Printf("Build Time: %s\nCommit: %s\nGo Version: %s\nUser: %s\n",
			version.BuildTime(), version.Commit(), version.GoVersion(), version.User())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
