package cmd

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	chat "github.com/ramble-cult/clhi/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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

		ctx := context.Background()

		conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			fmt.Println(err)
		}

		c := chat.NewBroadcastClient(conn)

		_, err = c.CreateGroup(ctx, &chat.CreateGroupReq{Name: g, Password: "test", Users: []string{}})
		if err != nil {
			fmt.Println("Error creating group:", err)
		}
		fmt.Println("Group created successfully")
	},
}
