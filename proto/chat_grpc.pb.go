// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: proto/chat.proto

package chat

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BroadcastClient is the client API for Broadcast service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BroadcastClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (Broadcast_StreamClient, error)
	CreateGroup(ctx context.Context, in *CreateGroupReq, opts ...grpc.CallOption) (*CreateResponse, error)
	JoinGroup(ctx context.Context, in *JoinReq, opts ...grpc.CallOption) (*JoinRes, error)
	Login(ctx context.Context, in *User, opts ...grpc.CallOption) (*LoginRes, error)
	ListUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersRes, error)
	ListGroups(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListGroupsRes, error)
}

type broadcastClient struct {
	cc grpc.ClientConnInterface
}

func NewBroadcastClient(cc grpc.ClientConnInterface) BroadcastClient {
	return &broadcastClient{cc}
}

func (c *broadcastClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Broadcast_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broadcast_ServiceDesc.Streams[0], "/chat.Broadcast/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &broadcastStreamClient{stream}
	return x, nil
}

type Broadcast_StreamClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type broadcastStreamClient struct {
	grpc.ClientStream
}

func (x *broadcastStreamClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *broadcastStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *broadcastClient) CreateGroup(ctx context.Context, in *CreateGroupReq, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/chat.Broadcast/CreateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) JoinGroup(ctx context.Context, in *JoinReq, opts ...grpc.CallOption) (*JoinRes, error) {
	out := new(JoinRes)
	err := c.cc.Invoke(ctx, "/chat.Broadcast/JoinGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) Login(ctx context.Context, in *User, opts ...grpc.CallOption) (*LoginRes, error) {
	out := new(LoginRes)
	err := c.cc.Invoke(ctx, "/chat.Broadcast/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) ListUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersRes, error) {
	out := new(ListUsersRes)
	err := c.cc.Invoke(ctx, "/chat.Broadcast/ListUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) ListGroups(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListGroupsRes, error) {
	out := new(ListGroupsRes)
	err := c.cc.Invoke(ctx, "/chat.Broadcast/ListGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BroadcastServer is the server API for Broadcast service.
// All implementations must embed UnimplementedBroadcastServer
// for forward compatibility
type BroadcastServer interface {
	Stream(Broadcast_StreamServer) error
	CreateGroup(context.Context, *CreateGroupReq) (*CreateResponse, error)
	JoinGroup(context.Context, *JoinReq) (*JoinRes, error)
	Login(context.Context, *User) (*LoginRes, error)
	ListUsers(context.Context, *emptypb.Empty) (*ListUsersRes, error)
	ListGroups(context.Context, *emptypb.Empty) (*ListGroupsRes, error)
	mustEmbedUnimplementedBroadcastServer()
}

// UnimplementedBroadcastServer must be embedded to have forward compatible implementations.
type UnimplementedBroadcastServer struct {
}

func (UnimplementedBroadcastServer) Stream(Broadcast_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedBroadcastServer) CreateGroup(context.Context, *CreateGroupReq) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (UnimplementedBroadcastServer) JoinGroup(context.Context, *JoinReq) (*JoinRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinGroup not implemented")
}
func (UnimplementedBroadcastServer) Login(context.Context, *User) (*LoginRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedBroadcastServer) ListUsers(context.Context, *emptypb.Empty) (*ListUsersRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedBroadcastServer) ListGroups(context.Context, *emptypb.Empty) (*ListGroupsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGroups not implemented")
}
func (UnimplementedBroadcastServer) mustEmbedUnimplementedBroadcastServer() {}

// UnsafeBroadcastServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BroadcastServer will
// result in compilation errors.
type UnsafeBroadcastServer interface {
	mustEmbedUnimplementedBroadcastServer()
}

func RegisterBroadcastServer(s grpc.ServiceRegistrar, srv BroadcastServer) {
	s.RegisterService(&Broadcast_ServiceDesc, srv)
}

func _Broadcast_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BroadcastServer).Stream(&broadcastStreamServer{stream})
}

type Broadcast_StreamServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type broadcastStreamServer struct {
	grpc.ServerStream
}

func (x *broadcastStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *broadcastStreamServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Broadcast_CreateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).CreateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Broadcast/CreateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).CreateGroup(ctx, req.(*CreateGroupReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_JoinGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).JoinGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Broadcast/JoinGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).JoinGroup(ctx, req.(*JoinReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Broadcast/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).Login(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Broadcast/ListUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).ListUsers(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_ListGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).ListGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.Broadcast/ListGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).ListGroups(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Broadcast_ServiceDesc is the grpc.ServiceDesc for Broadcast service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broadcast_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Broadcast",
	HandlerType: (*BroadcastServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGroup",
			Handler:    _Broadcast_CreateGroup_Handler,
		},
		{
			MethodName: "JoinGroup",
			Handler:    _Broadcast_JoinGroup_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Broadcast_Login_Handler,
		},
		{
			MethodName: "ListUsers",
			Handler:    _Broadcast_ListUsers_Handler,
		},
		{
			MethodName: "ListGroups",
			Handler:    _Broadcast_ListGroups_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Broadcast_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/chat.proto",
}
