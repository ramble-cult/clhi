package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/manifoldco/promptui"
	chat "github.com/ramble-cult/clhi/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

func init() {
	rootCmd.AddCommand(joinRoom)
}

var joinRoom = &cobra.Command{
	Use:   "join",
	Short: "log into server",
	Long:  `log into server -H <host> -u <username> -p <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Prompt{
			Label: "room name",
		}
		group, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
		}

		c := viper.Get("Client").(chat.BroadcastClient)
		viper.ReadInConfig()
		t := viper.GetViper().GetString("user-token")
		u := viper.GetString("User")

		ctx := viper.Get("context").(context.Context)
		md := metadata.New(map[string]string{"user-token": t})
		ctx = metadata.NewOutgoingContext(ctx, md)

		c.JoinGroup(ctx, &chat.JoinReq{
			Name: group,
			User: u,
		})

		md = metadata.New(map[string]string{"user-token": t, "user-group": group})
		ctx = metadata.NewOutgoingContext(ctx, md)

		client, err := c.Stream(ctx)
		if err != nil {
			fmt.Println("error connecting to stream")
			return
		}
		defer client.CloseSend()

		err = stream(ctx, c)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v connected to stream", time.Now())
	},
}

func stream(ctx context.Context, c chat.BroadcastClient) error {
	client, err := c.Stream(ctx)
	if err != nil {
		return fmt.Errorf("error connecting to stream: %v", err)
	}
	defer client.CloseSend()

	fmt.Printf("%v connected to stream", time.Now())
	return sendAndReceive(ctx, client)
}

func sendAndReceive(ctx context.Context, client chat.Broadcast_StreamClient) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic:", r)
		}
	}()

	go send(client)
	return receive(client)
}

func send(client chat.Broadcast_StreamClient) {
	user := viper.GetString("User")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		msg := sc.Text()
		if err := client.Send(&chat.Message{Username: user, Message: msg}); err != nil {
			log.Printf("failed to send message: %v", err)
			break
		}
	}
	if err := sc.Err(); err != nil {
		log.Printf("input scanner failure: %v", err)
	}
}

func receive(sc chat.Broadcast_StreamClient) error {
	for {
		res, err := sc.Recv()
		if err == io.EOF {
			// Stream has been closed by the server
			log.Println("Server closed the stream")
			return nil
		}
		if err != nil {
			// Other errors occurred
			return fmt.Errorf("error receiving message from server: %v", err)
		}

		log.Printf("%s: %s", res.Username, res.Message)
	}
}
