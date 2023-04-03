package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "create room sessions",
	Long:  "Create room sessions. You'll be provided with a hash permitting connections to room to share with peers.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
		name, err := RoomNamePrompt()
		if err != nil {
			panic(err)
		}

		fmt.Println(name)
		//create room uuid, get port and ip
		//create hash and print it out
		/* Prompt for room name

		Create room hash and print to console
		*/

		/*
			Give name to room
			Create hash of server details and print to console.
			Show prompt to start session
			Start session to read and write message to and from stdout and socket
		*/
	},
}

func init() {
	//hostCmd.Flags().BoolP("host", "o", true, hostCmd.Short)
	rootCmd.AddCommand(hostCmd)
}

func RoomNamePrompt() (string, error) {
	validate := func(input string) error {
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}
	fmt.Println("Before prompt")
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
	fmt.Println("After prompt")

	return res, nil
}
