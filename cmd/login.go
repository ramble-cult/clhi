package cmd

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	chat "github.com/ramble-cult/clhi/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	password string
	username string
)

func init() {
	login.Flags().StringVarP(&host, "host", "H", "0.0.0.0:50051", "Chat server host")
	login.Flags().StringVarP(&password, "password", "p", "", "User password")
	login.Flags().StringVarP(&username, "username", "u", "", "Username")
	rootCmd.AddCommand(login)

}

var login = &cobra.Command{
	Use:   "login",
	Short: "log into server",
	Long:  `log into server -H <host> -u <username> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		hostPrompt := promptui.Prompt{
			Label: "host",
		}
		h, err := hostPrompt.Run()
		if err != nil {
			fmt.Println(err)
		}

		namePrompt := promptui.Prompt{
			Label: "username",
		}

		u, err := namePrompt.Run()
		if err != nil {
			fmt.Println(err)
		}

		passPrompt := promptui.Prompt{
			Label: "password",
		}
		p, err := passPrompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		ctx := context.Background()

		conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Println(err)
		}

		client := chat.NewBroadcastClient(conn)

		viper.SetDefault("Connection", conn)
		viper.SetDefault("Client", client)

		res, err := client.Login(ctx, &chat.User{
			Host:     h,
			Password: p,
			Name:     u,
		})

		t := res.Token

		if err != nil {
			fmt.Println("incorrect credentials")
			return
		}

		viper.Set("user-token", t)
		viper.Set("User", u)
		viper.Set("Password", p)

		viper.WriteConfig()
	},
}
