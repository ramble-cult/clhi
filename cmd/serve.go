package cmd

import (
	"context"
	"fmt"

	"github.com/ramble-cult/clhi/chat"
	"github.com/spf13/cobra"
)

var (
	serverHost     string
	serverPassword string
)

func init() {
	serve.Flags().StringVarP(&serverHost, "serverhost", "H", "0.0.0.0:50051", "server host")
	serve.Flags().StringVarP(&serverPassword, "password", "p", "", "server password")
	rootCmd.AddCommand(serve)
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "serve -H <host> -p <password>",
	Long:  `serve -H <host> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := chat.SignalContext(context.Background())
		chat.Server(serverHost, serverPassword).Run(ctx)
		fmt.Println("starting server")
	},
}
