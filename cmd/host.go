package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/quinton11/chatline/config"
	"github.com/quinton11/chatline/internal/cui"
	"github.com/quinton11/chatline/internal/socket"
	"github.com/quinton11/chatline/internal/utils"

	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "create room sessions",
	Long:  "Create room sessions. You'll be provided with a hash permitting connections to room to share with peers.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)

		//read from config file
		//start room session
		var roomConfig utils.RoomConfig
		err := config.ReadConfig(&roomConfig)
		if err != nil {
			panic(err)
		}

		fmt.Println(roomConfig)
		fmt.Printf("Hosting Room: %s", roomConfig.Room.Name)

		//verify server ip and port
		//start socket

		server := socket.NewServer(roomConfig.Room)
		go server.Init()

		console := cui.NewChatUi(server, true)

		//console ui
		console.Init()
		/*
			ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlNiMjl0SWpwN0ltNWhiV1VpT2lKM1lXNXdhWE4xUUd4MVptWjVJaXdpZFhWcFpDSTZJamMwWWpneFpUTm1MV1kyTWpBdE5HVmxaQzFoTjJNNExUTTFNbVE0T1RJd01XWmhOQ0lzSW1odmMzUWlPaUl4T1RJdU1UWTRMakV3TUM0eE1pSXNJbkJ2Y25RaU9qVTFOREI5TENKbGVIQWlPakUyT0RFeU1qTXhORGQ5LmFIaFdjMDVaNmlDTlJoWEhPaFJram1EbWstaUo5N1lGRHJjcnFGc0hEMHM=

			NzRiODFlM2YtZjYyMC00ZWVkLWE3YzgtMzUyZDg5MjAxZmE0
		*/

	},
}

func init() {
	//hostCmd.Flags().BoolP("host", "o", true, hostCmd.Short)
	rootCmd.AddCommand(hostCmd)
}

func RoomHashPrompt() (string, error) {
	validate := func(input string) error {
		/* validate if signed obj is a room type */
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}
	prompt := promptui.Prompt{
		Label:       "Room Hash: ",
		HideEntered: true,
		Validate:    validate,
		Templates:   templates,
	}

	res, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return res, nil
}
