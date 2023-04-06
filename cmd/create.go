package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/quinton11/chatline/config"
	"github.com/quinton11/chatline/internal/utils"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a room session",
	Long:  "Use [chatline create] to create a room session and receive a room hash to share with peers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
		name, err := CreateRoomPrompt()
		if err != nil {
			panic(err)
		}
		fmt.Println(name)

		//get ip
		ip, err := utils.GetServerIp()
		if err != nil {
			panic(err)
		}
		fmt.Println(ip)

		//create room uuid, get port and ip
		//port
		port := 5540
		uuid := utils.GenerateRoomHash()

		room := utils.Room{Name: name, Uuid: uuid, Host: ip, Port: port}

		//convert to jwt
		token, err := utils.GenerateToken(room)
		if err != nil {
			panic(err)
		}

		//create hash and print it out
		rHash, sHash := utils.GenerateHash(token, room)
		fmt.Println("")
		fmt.Printf("RoomHash: %s \n", rHash)
		fmt.Println("")

		fmt.Printf("RoomKey: %s \n", sHash)
		err = config.WriteConfig(utils.RoomConfig{Room: room, Key: sHash})
		if err != nil {
			panic(err)
		}

		/*
			Give name to room
			Create hash of server details and print to console.
			Show prompt to start session
			Start session to read and write message to and from stdout and socket
		*/
	},
}

func init() {
	//start
	//startCmd.Flags().BoolP("start", "s", true, startCmd.Short)
	rootCmd.AddCommand(createCmd)
}

func CreateRoomPrompt() (string, error) {
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
		Label:       "Room Name: ",
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
