package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(joinRoom)
}

var joinRoom = &cobra.Command{
	Use:   "join",
	Short: "log into server",
	Long:  `log into server -H <host> -u <username> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Prompt{
			Label: "room name",
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
		}
		Client.JoinGroup(Ctx, result)
	},
}
