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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
	rootCmd.AddCommand(join)
}

var user string

var join = &cobra.Command{
	Use:   "join",
	Short: "log into server",
	Long:  `log into server -H <host> -u <username> -p <password>`,
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
			fmt.Println(err)
		}

		defer conn.Close()

		c := chat.NewBroadcastClient(conn)

		res, err := c.ListGroups(ctx, &emptypb.Empty{})
		if err != nil {
			fmt.Println("error listing group", err)
		}

		groupPrompt := promptui.Select{
			Label: "Available Groups",
			Items: res.Groups,
		}

		_, group, err := groupPrompt.Run()
		if err != nil {
			fmt.Println(err)
		}

		c.JoinGroup(ctx, &chat.JoinReq{
			Name: group,
			User: user,
		})

		md = metadata.New(map[string]string{"user-token": t, "user-group": group})
		ctx = metadata.NewOutgoingContext(ctx, md)

		err = stream(ctx, c)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v connected to stream \n", time.Now().Format("15:04"))
	},
}

func stream(ctx context.Context, c chat.BroadcastClient) error {
	client, err := c.Stream(ctx)
	if err != nil {
		return fmt.Errorf("error connecting to stream: %v", err)
	}
	defer client.CloseSend()

	fmt.Printf("%v connected to stream\n", time.Now().Format("15:04"))
	go send(client)
	return receive(client)
}

func send(client chat.Broadcast_StreamClient) {
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)
	u := viper.GetString("user")

	for {
		select {
		case <-client.Context().Done():
			// DebugLogf("client send loop disconnected")
		default:
			if sc.Scan() {
				if err := client.Send(&chat.Message{Username: u, Message: sc.Text()}); err != nil {
					// ClientLogf(time.Now(), "failed to send message: %v", err)
					return
				}
			} else {
				// ClientLogf(time.Now(), "input scanner failure: %v", sc.Err())
				return
			}
		}
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

		fmt.Printf("%s: %s\n", res.Username, res.Message)
	}
}
