package cmd

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	chat "github.com/ramble-cult/clhi/proto"
	"github.com/ramble-cult/clhi/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	password string
	username string
	Client   *services.Client
	Ctx      context.Context
)

func init() {
	login.Flags().StringVarP(&password, "password", "p", "", "User password")
	login.Flags().StringVarP(&username, "username", "u", "", "Username")
	rootCmd.AddCommand(login)
}

var login = &cobra.Command{
	Use:   "login",
	Short: "log into server",
	Long:  `log into server -H <host> -u <username> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		namePrompt := promptui.Prompt{
			Label: "username",
		}
		u, err := namePrompt.Run()
		if err != nil {
			fmt.Println(err)
		}
		username = u
		passPrompt := promptui.Prompt{
			Label: "password",
		}
		p, err := passPrompt.Run()
		if err != nil {
			fmt.Println(err)
		}

		c := viper.Get("client").(chat.BroadcastClient)
		ctx := viper.Get("context").(context.Context)
		c.Login(ctx, &chat.User{
			Host:     host,
			Name:     u,
			Password: p,
		})

	},
}
