package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	guuid "github.com/google/uuid"
	chat "github.com/ramble-cult/clhi/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type user struct {
	stream chan *chat.Message
	name   string
	error  chan error
}

type group struct {
	chat.UnimplementedBroadcastServer
	Broadcast chan *chat.Message
	Password  string
	Name      string
	Users     map[string]*user
	Error     chan error
	mu        sync.RWMutex
}

type server struct {
	chat.UnimplementedBroadcastServer
	Host        string
	OnlineUsers map[string]*user
	Groups      map[string]*group
	Password    string
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

	uid := guuid.New().String()

	newUser := &user{
		name:   u.Name,
		stream: make(chan *chat.Message, 100),
		error:  make(chan error),
	}

	s.OnlineUsers[uid] = newUser
	return &chat.LoginRes{Token: uid}, nil
}

func (s *server) ListUsers(context.Context, *chat.ListUsersReq) (*chat.ListUsersRes, error) {
	u := &chat.ListUsersRes{}
	for _, v := range s.OnlineUsers {
		u.Users = append(u.Users, &chat.UserResponse{
			Name: v.name,
		})
	}

	return u, nil
}

func (s *server) CreateGroup(ctx context.Context, req *chat.CreateGroupReq) (*chat.CreateResponse, error) {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	// Validate user token and group existence
	// user, ok := s.OnlineUsers[chat.Token]
	// if !ok {
	// 	return errors.New("user not found")
	// }

	// Check if the group already exists
	_, ok := s.Groups[req.Name]
	if ok {
		return nil, errors.New("group already exists")
	}

	users := map[string]*user{}

	for _, v := range req.Users {
		for k, o := range s.OnlineUsers {
			if o.name == v {
				users[k] = o
			}
		}
	}

	// Create a new group
	convo := &group{
		Broadcast: make(chan *chat.Message),
		Password:  req.Password,
		Name:      req.Name,
		Users:     users,
		Error:     make(chan error),
	}

	s.Groups[req.Name] = convo

	return &chat.CreateResponse{Message: "room created"}, nil
}

func (s *server) JoinGroup(ctx context.Context, req *chat.JoinReq) (*chat.JoinRes, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("invalid token")
	}

	token := md["user-token"][0]

	group, ok := s.Groups[req.Name]
	if !ok {
		return nil, errors.New("group not found")
	}

	user, ok := s.OnlineUsers[token]
	if !ok {
		return nil, errors.New("user not found")
	}

	// Check if the user is already in the group
	for _, u := range group.Users {
		if u.name == req.User {
			return nil, errors.New("user already in the group")
		}
	}

	// Add the user to the group
	group.mu.Lock()
	defer group.mu.Unlock()
	group.Users[token] = user

	go group.broadcast(ctx)

	return &chat.JoinRes{Message: "successfully joined"}, nil
}

func (s *server) Start(ctx context.Context) error {
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()
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

func (s *server) Stream(srv chat.Broadcast_StreamServer) error {
	md, _ := metadata.FromIncomingContext(srv.Context())
	token := md["user-token"][0]
	groupName := md["user-group"][0]

	g, ok := s.Groups[groupName]
	if !ok {
		return status.Error(codes.Unauthenticated, "missing group name")
	}

	if _, ok := g.Users[token]; !ok {
		return status.Error(codes.Unauthenticated, "missing token")
	}

	go g.sendBroadcasts(srv, token)

	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		g.Broadcast <- &chat.Message{
			Username: req.Username,
			Message:  req.Message,
		}
	}

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (g *group) broadcast(_ context.Context) {
	for res := range g.Broadcast {
		g.mu.RLock()
		for _, u := range g.Users {
			select {
			case u.stream <- res:
				// keep stream open
			default:
				fmt.Printf("%v client stream full, dropping message", time.Now())
			}
		}
		g.mu.RUnlock()
	}
}

func (g *group) sendBroadcasts(srv chat.Broadcast_StreamServer, tkn string) {
	stream := g.Users[tkn].stream

	for {
		select {
		case <-srv.Context().Done():
			return
		case res := <-stream:
			if s, ok := status.FromError(srv.Send(res)); ok {
				switch s.Code() {
				case codes.OK:
					// noop
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:

					return
				default:

					return
				}
			}
		}
	}
}
