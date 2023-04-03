package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join chat room session",
	Long:  "Use this command to join a live room session created by a host",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
	},
}

func init() {
	// Join Command
	//joinCmd.Flags().BoolP("join", "j", true, joinCmd.Short)
	rootCmd.AddCommand(joinCmd)
}
