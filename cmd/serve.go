package cmd

import (
	"context"
	"fmt"

	"github.com/ramble-cult/clhi/services"
	"github.com/spf13/cobra"
)

var (
	serverPassword string
)

func init() {
	serve.Flags().StringVarP(&serverPassword, "password", "p", "", "server password")
	rootCmd.AddCommand(serve)
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "serve -H <host> -p <password>",
	Long:  `serve -H <host> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := services.SignalContext(context.Background())
		services.Server(host, serverPassword).Start(ctx)
		fmt.Println("starting server")
	},
}
