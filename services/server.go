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
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type user struct {
	stream chan *chat.Message
	name   string
	error  chan error
}

type group struct {
	Broadcast chan *chat.Message
	Password  string
	Name      string
	Users     map[string]*user
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

func (s *server) ListUsers(context.Context, *emptypb.Empty) (*chat.ListUsersRes, error) {
	u := &chat.ListUsersRes{}
	for _, v := range s.OnlineUsers {
		u.Users = append(u.Users, &chat.UserResponse{
			Name: v.name,
		})
	}

	return u, nil
}

func (s *server) ListGroups(context.Context, *emptypb.Empty) (*chat.ListGroupsRes, error) {
	res := &chat.ListGroupsRes{}
	for _, v := range s.Groups {
		res.Groups = append(res.Groups, v.Name)
	}

	return res, nil
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
	}

	s.Groups[req.Name] = convo

	go convo.broadcast(ctx)

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
	u, ok := g.Users[token]
	if !ok {
		return status.Error(codes.Unauthenticated, "missing token")
	}

	errs, ctx := errgroup.WithContext(srv.Context())
	stream := u.stream

	errs.Go(func() error {
		for {
			select {
			case <-ctx.Done():
			// none
			case res := <-stream:
				if s, ok := status.FromError(srv.Send(res)); ok {
					switch s.Code() {
					case codes.OK:
					case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					default:
					}
				}
			}
		}
	})

	errs.Go(func() error {
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
		return nil
	})

	// <-srv.Context().Done()
	// return srv.Context().Err()
	return errs.Wait()
}

func (g *group) broadcast(_ context.Context) {
	for res := range g.Broadcast {
		g.mu.RLock()
		for _, u := range g.Users {
			if res.Username != u.name {

				select {
				case u.stream <- res:
					// keep stream open
				default:
					fmt.Printf("%v client stream full, dropping message", time.Now())
				}
			}
		}
		g.mu.RUnlock()
	}
}
