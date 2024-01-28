package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	guuid "github.com/google/uuid"
	chat "github.com/ramble-cult/clhi/proto"
	"google.golang.org/grpc"
)

type user struct {
	chat.UnimplementedBroadcastServer
	stream chat.Broadcast_CreateGroupStreamServer
	id     uuid.UUID
	name   string
	error  chan error
}

type group struct {
	chat.UnimplementedBroadcastServer
	Stream   chat.Broadcast_CreateGroupStreamServer
	Password string
	Name     string
	Id       int
	Users    []*user
	Error    chan error
	mu       sync.RWMutex
}

type server struct {
	chat.UnimplementedBroadcastServer
	Host        string
	OnlineUsers map[string]*user
	Groups      map[string]*group
	Password    string
	mu          sync.RWMutex
}

func Server(host, pass string) *server {
	return &server{
		Host:        host,
		OnlineUsers: make(map[string]*user),
		Groups:      map[string]*group{},
		Password:    pass,
	}
}

func (s *server) Login(ctx context.Context, u *chat.User) (*chat.LoginRes, error) {

	if u.Password != s.Password {
		return nil, errors.New("incorrect password")
	}

	if _, ok := s.OnlineUsers[u.Name]; ok {
		return nil, errors.New("username taken, try a new name")
	}

	newUser := &user{
		id:     guuid.New(),
		name:   u.Name,
		stream: nil,
		error:  make(chan error),
	}

	s.OnlineUsers[u.Name] = newUser
	return &chat.LoginRes{Token: newUser.id.String()}, nil
}

func (s *server) ListUsers(context.Context, *chat.ListUsersReq) (*chat.ListUsersRes, error) {
	u := &chat.ListUsersRes{} // Initialize u
	for k := range s.OnlineUsers {
		u.Users = append(u.Users, &chat.UserResponse{
			Name: k,
		})
	}

	return u, nil
}

func (s *server) CreateGroupStream(chat *chat.CreateGroup, stream chat.Broadcast_CreateGroupStreamServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate user token and group existence
	// user, ok := s.OnlineUsers[chat.Token]
	// if !ok {
	// 	return errors.New("user not found")
	// }

	// Check if the group already exists
	_, ok := s.Groups[chat.Name]
	if ok {
		// Add user to the existing group
		// existingGroup.mu.Lock()
		// existingGroup.Users = append(existingGroup.Users, user.name)
		// existingGroup.mu.Unlock()
		return nil
	}

	users := []*user{}

	for _, v := range chat.Users {
		if _, ok := s.OnlineUsers[v]; ok {
			users = append(users, s.OnlineUsers[v])
		}
	}

	// Create a new group
	convo := &group{
		Stream:   stream,
		Password: chat.Password,
		Name:     chat.Name,
		Users:    users,
		Error:    make(chan error),
	}

	// Lock the mutex before writing to the map
	s.mu.Lock()
	s.Groups[chat.Name] = convo
	s.mu.Unlock()

	return <-convo.Error
}

func (s *server) JoinGroup(req *chat.JoinGroupReq, _ chat.Broadcast_JoinGroupServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	group, ok := s.Groups[req.Name]
	if !ok {
		return errors.New("Group not found")
	}

	user, ok := s.OnlineUsers[req.User]
	if !ok {
		return errors.New("user not found")
	}

	group.mu.Lock()
	defer group.mu.Unlock()

	// Check if the user is already in the group
	for _, u := range group.Users {
		if u.name == req.User {
			return errors.New("User already in the group")
		}
	}

	// Add the user to the group
	group.Users = append(group.Users, user)

	return nil
}

func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	rpc := grpc.NewServer()
	chat.RegisterBroadcastServer(rpc, s)

	l, err := net.Listen("tcp", s.Host)
	if err != nil {
		return err
	}

	fmt.Printf("Server started at %v", s.Host)

	if err := rpc.Serve(l); err != nil {
		return err
	}
	return nil
}

func SignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		<-sigs

		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}

func (g *group) BroadcastMessage(msg *chat.Message) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, user := range g.Users {
		if user.name != msg.Username {
			user.stream.Send(msg)
		}
	}
}

func (g *group) DirectMessage(ctx context.Context, msg *chat.Message) (*chat.CreateResponse, error) {
	g.BroadcastMessage(msg)
	return &chat.CreateResponse{Message: "Message sent successfully"}, nil
}
