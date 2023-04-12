package cmd

import (
	"fmt"
	"log"

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
			log.Fatal(err)
		}

		server := socket.NewServer(roomConfig.Room)
		go server.Init()

		console := cui.NewChatUi(server)

		/*
			ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlNiMjl0SWpwN0ltNWhiV1VpT2lKbmIzQm9aWEp6SWl3aWRYVnBaQ0k2SW1ObVl6SXlNVGd5TFRnME5ESXROR00zTWkxaE5qUTVMVGc1TURVeU1XWTJaV00yWmlJc0ltaHZjM1FpT2lJeE9USXVNVFk0TGpFd01DNHhNaUlzSW5CdmNuUWlPalUxTkRCOUxDSmxlSEFpT2pFMk9ERXlPVGczTkROOS55RWFkbG1fcXZBWnpBYktqYVIyYUFaeHRUMVRuek9FZ1RLMEVqWXR4OGZn

			Y2ZjMjIxODItODQ0Mi00YzcyLWE2NDktODkwNTIxZjZlYzZm
		*/

		//console ui
		console.Init()
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}

func RoomHashPrompt() (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}
	prompt := promptui.Prompt{
		Label:       "Room Hash: ",
		HideEntered: true,
		Templates:   templates,
	}

	res, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return res, nil
}
