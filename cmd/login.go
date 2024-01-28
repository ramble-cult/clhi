package cmd

import (
	"context"

	"github.com/ramble-cult/clhi/services"
	"github.com/spf13/cobra"
)

var (
	host     string
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
		ctx := services.SignalContext(context.Background())
		services.Client(host, password, username).Start(ctx)
	},
}
