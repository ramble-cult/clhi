package cmd

import (
	"context"
	"fmt"

	"github.com/ramble-cult/clhi/chat"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serve)
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := chat.SignalContext(context.Background())
		chat.Server("localhost", "password").Run(ctx)
		fmt.Println("starting server")
	},
}
