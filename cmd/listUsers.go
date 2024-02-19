/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	chat "github.com/ramble-cult/clhi/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "lu",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := viper.Get("Client").(chat.BroadcastClient)
		t := viper.GetString("Token")
		u, err := c.ListUsers(context.Background(), &chat.ListUsersReq{Token: t})
		if err != nil {
			return
		}

		for _, v := range u.Users {
			fmt.Println(v.Name)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(listUsersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listUsersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listUsersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
