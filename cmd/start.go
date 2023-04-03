package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a room session after hosting",
	Long:  "Use after you've created a room session with the host command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
	},
}

func init() {
	//start
	//startCmd.Flags().BoolP("start", "s", true, startCmd.Short)
	rootCmd.AddCommand(startCmd)
}
