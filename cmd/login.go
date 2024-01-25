package cmd

import (
	"context"
	"fmt"

	"github.com/ramble-cult/clhi/chat"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(login)
}

var login = &cobra.Command{
	Use:   "login",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := chat.SignalContext(context.Background())
		chat.Client("localhost", "password", "username").Run(ctx)
		fmt.Println("you have logged in")
	},
}
