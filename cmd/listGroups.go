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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// listGroupsCmd represents the listGroups command
var listGroupsCmd = &cobra.Command{
	Use:   "lg",
	Short: "List Available Groups",
	Long:  `List Available Groups`,
	Run: func(cmd *cobra.Command, args []string) {

		t := viper.GetString("user-token")
		user = viper.GetString("user")

		ctx, cancel := context.WithCancel(context.Background())
		defer ctx.Done()
		defer cancel()
		md := metadata.New(map[string]string{"user-token": t})
		ctx = metadata.NewOutgoingContext(ctx, md)

		conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Println("error listing group", err)
		}

		defer conn.Close()

		c := chat.NewBroadcastClient(conn)
		res, err := c.ListGroups(ctx, &emptypb.Empty{})
		if err != nil {
			fmt.Println("error listing group", err)
		}

		for _, v := range res.Groups {
			fmt.Println(v)
		}
	},
}

func init() {
	rootCmd.AddCommand(listGroupsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listGroupsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listGroupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
