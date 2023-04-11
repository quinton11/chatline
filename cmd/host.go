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
		//go console.UpdateChats()
		console.Init()
		/*
			ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlNiMjl0SWpwN0ltNWhiV1VpT2lKblpYSnRZU0lzSW5WMWFXUWlPaUk1WTJFNU56SmtPQzA1T0RFd0xUUmhNV0l0T0RVMFlpMDBZMlJoTVdZNU56Z3lZalVpTENKb2IzTjBJam9pTVRreUxqRTJPQzR4TURBdU1USWlMQ0p3YjNKMElqbzFOVFF3ZlN3aVpYaHdJam94TmpneE1qVXlNRFkwZlEuWFZ2VVV3T1F6cGVYQ2lyZUZlTEVKQTdxcjF3Y0dUTHFQX1Bkc2pSbnN2QQ==

			OWNhOTcyZDgtOTgxMC00YTFiLTg1NGItNGNkYTFmOTc4MmI1
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
