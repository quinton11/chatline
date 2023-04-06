package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/quinton11/chatline/internal/utils"
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join chat room session",
	Long:  "Use this command to join a live room session created by a host",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
		roomCrypted, err := AskRoomCredentials()
		if err != nil {
			panic(err)
		}

		tokenDecoded, err := utils.Decode64(roomCrypted.Hash)
		if err != nil {
			panic(err)
		}

		//validate room
		room, err := utils.ValidateToken(tokenDecoded, roomCrypted.Key)
		if err != nil {
			panic(err)
		}

		fmt.Println(room)

		//start socket client and attempt connection
		//using server ip && port
		//create temp name
		//start session
	},
}

func init() {
	// Join Command
	//joinCmd.Flags().BoolP("join", "j", true, joinCmd.Short)
	rootCmd.AddCommand(joinCmd)
}

func AskRoomCredentials() (utils.RoomCrypted, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	promptHash := promptui.Prompt{
		Label:       "Room Hash: ",
		HideEntered: true,
		Mask:        '*',
		Templates:   templates,
	}

	rhash, err := promptHash.Run()
	if err != nil {
		return utils.RoomCrypted{}, err
	}

	promptKey := promptui.Prompt{
		Label:       "Room Key: ",
		HideEntered: true,
		Mask:        '*',
		Templates:   templates,
	}

	khash, err := promptKey.Run()
	if err != nil {
		return utils.RoomCrypted{}, err
	}

	return utils.RoomCrypted{Hash: rhash, Key: khash}, nil
}
