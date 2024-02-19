package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	chat "github.com/ramble-cult/clhi/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		g, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
		}

		c := viper.Get("Client").(chat.BroadcastClient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = c.CreateGroup(ctx, &chat.CreateGroupReq{Name: g, Password: "test", Users: []string{}})
		if err != nil {
			fmt.Println("Error creating group:", err)
		}
		fmt.Println("Group created successfully")
	},
}
