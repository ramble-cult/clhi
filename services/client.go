package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	chat "github.com/ramble-cult/clhi/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	chat.BroadcastClient
	Host, Password, Name, Token, GroupName string
	Shutdown                               bool
}

func NewClient(host, pass, name string) *Client {
	return &Client{
		Host:     host,
		Password: pass,
		Name:     name,
	}
}

func (c *Client) Start(ctx context.Context) error {
	connCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(connCtx, c.Host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c.BroadcastClient = chat.NewBroadcastClient(conn)

	if c.Token, err = c.login(ctx); err != nil {
		return err
	}

	md := metadata.New(map[string]string{"user-token": c.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	chatting := false
	for !chatting {

		var cmd string
		fmt.Scanln(&cmd)
		switch cmd {
		case "ls":
			u, err := c.ListUsers(ctx)
			if err != nil {
				fmt.Println(err)
			}
			for _, v := range u {
				fmt.Println(v)
			}
		case "c":

			_, err = c.BroadcastClient.CreateGroup(ctx, &chat.CreateGroupReq{Name: "test", Password: "test", Users: []string{}})
			if err != nil {
				fmt.Println("Error creating group:", err)
				return err
			}
			fmt.Println("Group created successfully")
		case "j":
			fmt.Println("Enter group name to join: ")
			fmt.Scanln(&c.GroupName)
			_, err := c.BroadcastClient.JoinGroup(ctx, &chat.JoinReq{Name: "test", User: c.Name})
			if err != nil {
				fmt.Println("Error joining group:", err)
			} else {
				fmt.Println("Joined group successfully")
				chatting = true
			}
		}
	}

	err = c.stream(ctx)

	// _, err = c.BroadcastClient.CreateGroup(ctx, &chat.CreateGroupReq{Name: "test", Password: "test", Users: []string{}})
	// c.GroupName = "test"

	if err != nil {
		fmt.Println(err)
	}

	return errors.WithMessage(err, "stream error")
}

func (c *Client) JoinGroup(ctx context.Context, group string) error {
	c.BroadcastClient.JoinGroup(ctx, &chat.JoinReq{Name: group, User: c.Name})
	err := c.stream(ctx)

	return err
}

func (c *Client) CreateGroup(ctx context.Context, group string) error {
	_, err := c.BroadcastClient.CreateGroup(ctx, &chat.CreateGroupReq{Name: group, Password: "test", Users: []string{}})
	if err != nil {
		fmt.Println("Error creating group:", err)
		return err
	}
	fmt.Println("Group created successfully")

	return err
}

func (c *Client) ListUsers(ctx context.Context) ([]string, error) {
	u, err := c.BroadcastClient.ListUsers(ctx, &chat.ListUsersReq{Token: c.Token})
	if err != nil {
		return nil, err
	}

	var users []string
	for _, v := range u.Users {
		users = append(users, v.Name)
	}

	return users, nil
}

func (c *Client) login(ctx context.Context) (string, error) {
	res, err := c.BroadcastClient.Login(ctx, &chat.User{
		Name:     c.Name,
		Host:     c.Host,
		Password: c.Password,
	})
	if err != nil {
		return "", err
	}

	return res.Token, nil
}

func (c *Client) stream(ctx context.Context) error {
	md := metadata.New(map[string]string{"user-token": c.Token, "user-group": c.GroupName})
	ctx = metadata.NewOutgoingContext(ctx, md)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := c.BroadcastClient.Stream(ctx)
	if err != nil {
		return err
	}
	defer client.CloseSend()

	fmt.Printf("%v connected to stream", time.Now())
	go c.send(client)
	return c.receive(client)
}

func (c *Client) receive(sc chat.Broadcast_StreamClient) error {
	for {
		res, err := sc.Recv()

		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			return nil
		} else if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		log.Printf("%s:%s", res.Username, res.Message)
	}
}

func (c *Client) send(client chat.Broadcast_StreamClient) {
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)

	for {
		select {
		case <-client.Context().Done():
			// DebugLogf("client send loop disconnected")
		default:
			if sc.Scan() {
				if err := client.Send(&chat.Message{Username: c.Name, Message: sc.Text()}); err != nil {
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
