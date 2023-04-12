package cmd

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/quinton11/chatline/internal/cui"
	"github.com/quinton11/chatline/internal/socket"
	"github.com/quinton11/chatline/internal/utils"
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join chat room session",
	Long:  "Use this command to join a live room session created by a host",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
		roomCrypted, uname, err := AskRoomCredentials()
		if err != nil {
			log.Fatal(err)
		}

		tokenDecoded, err := utils.Decode64(roomCrypted.Hash)
		if err != nil {
			log.Fatal(err)
		}

		//validate room
		room, err := utils.ValidateToken(tokenDecoded, roomCrypted.Key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(uname)
		client := socket.NewClient(room, uname)

		go func() {
			err = client.Connect()
			if err != nil {
				//
				log.Fatal(err)
			}
		}()

		console := cui.NewChatUi(client)

		console.Init()
	},
}

func init() {
	// Join Command
	rootCmd.AddCommand(joinCmd)
}

func AskRoomCredentials() (utils.RoomCrypted, string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	username := promptui.Prompt{
		Label:       "UserName: ",
		HideEntered: true,
		Templates:   templates,
	}

	uname, err := username.Run()
	if err != nil {
		return utils.RoomCrypted{}, "", err
	}

	promptHash := promptui.Prompt{
		Label:       "Room Hash: ",
		HideEntered: true,
		Mask:        '*',
		Templates:   templates,
	}

	rhash, err := promptHash.Run()
	if err != nil {
		return utils.RoomCrypted{}, "", err
	}

	promptKey := promptui.Prompt{
		Label:       "Room Key: ",
		HideEntered: true,
		Mask:        '*',
		Templates:   templates,
	}

	khash, err := promptKey.Run()
	if err != nil {
		return utils.RoomCrypted{}, "", err
	}

	return utils.RoomCrypted{Hash: rhash, Key: khash}, uname, nil
}
