package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createGroup)
}

var createGroup = &cobra.Command{
	Use:   "create",
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
		fmt.Println(Ctx)
		Client.CreateGroup(Ctx, result)
	},
}
