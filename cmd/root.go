package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chatline",
	Short: "chat with remote peers via the command line",
	Long:  "Host and join rooms through which multiple peers can connect and chat",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
