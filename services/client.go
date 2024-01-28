package services

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	chat "github.com/ramble-cult/clhi/proto"
	"google.golang.org/grpc"
)

type client struct {
	chat.BroadcastClient
	Host, Password, Name, Token, GroupName string
	Shutdown                               bool
}

func Client(host, pass, name string) *client {
	return &client{
		Host:     host,
		Password: pass,
		Name:     name,
	}
}

func (c *client) Start(ctx context.Context) error {
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

	chatting := false
	for !chatting {

		var cmd string
		fmt.Scanln(&cmd)
		switch cmd {
		case "ls":
			u, err := c.listUsers(ctx)
			if err != nil {
				fmt.Println(err)
			}
			for _, v := range u {
				fmt.Println(v)
			}
		case "c":

			_, err = c.BroadcastClient.CreateGroupStream(ctx, &chat.CreateGroup{Name: "hera&caleb", Password: "test", Users: []string{"hera", "caleb"}})
			if err != nil {
				fmt.Println("Error creating group:", err)
				return err
			}
			fmt.Println("Group created successfully")

			// Start a goroutine to continuously send messages
			go c.sendMessageLoop(ctx)

			// Wait for the user to exit the application
			<-ctx.Done()
		case "j":
			fmt.Print("Enter group name to join: ")
			fmt.Scanln(&c.GroupName)
			_, err := c.BroadcastClient.JoinGroup(ctx, &chat.JoinGroupReq{Name: c.GroupName, User: c.Name})
			if err != nil {
				fmt.Println("Error joining group:", err)
			} else {
				fmt.Println("Joined group successfully")
				go c.sendMessageLoop(ctx)
			}
		}
	}
	return nil
}

func (c *client) sendMessageLoop(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Print("Enter message: ")
			scanner.Scan()
			text := scanner.Text()

			msg := &chat.Message{Username: c.Name, Message: text}
			_, err := c.BroadcastClient.DirectMessage(ctx, msg)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	}
}

func (c *client) listUsers(ctx context.Context) ([]string, error) {
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

func (c *client) login(ctx context.Context) (string, error) {
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
