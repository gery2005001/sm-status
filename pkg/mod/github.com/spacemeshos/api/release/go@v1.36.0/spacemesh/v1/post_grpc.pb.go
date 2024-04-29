// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: spacemesh/v1/post.proto

package spacemeshv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PostService_Register_FullMethodName = "/spacemesh.v1.PostService/Register"
)

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	// Register is a bi-directional stream that allows a dedicated PoST node to connect to the spacemesh node.
	// The node will send NodeRequets to PoST and the service will respond with ServiceResponses.
	Register(ctx context.Context, opts ...grpc.CallOption) (PostService_RegisterClient, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) Register(ctx context.Context, opts ...grpc.CallOption) (PostService_RegisterClient, error) {
	stream, err := c.cc.NewStream(ctx, &PostService_ServiceDesc.Streams[0], PostService_Register_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &postServiceRegisterClient{stream}
	return x, nil
}

type PostService_RegisterClient interface {
	Send(*ServiceResponse) error
	Recv() (*NodeRequest, error)
	grpc.ClientStream
}

type postServiceRegisterClient struct {
	grpc.ClientStream
}

func (x *postServiceRegisterClient) Send(m *ServiceResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *postServiceRegisterClient) Recv() (*NodeRequest, error) {
	m := new(NodeRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations should embed UnimplementedPostServiceServer
// for forward compatibility
type PostServiceServer interface {
	// Register is a bi-directional stream that allows a dedicated PoST node to connect to the spacemesh node.
	// The node will send NodeRequets to PoST and the service will respond with ServiceResponses.
	Register(PostService_RegisterServer) error
}

// UnimplementedPostServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPostServiceServer struct {
}

func (UnimplementedPostServiceServer) Register(PostService_RegisterServer) error {
	return status.Errorf(codes.Unimplemented, "method Register not implemented")
}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_Register_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PostServiceServer).Register(&postServiceRegisterServer{stream})
}

type PostService_RegisterServer interface {
	Send(*NodeRequest) error
	Recv() (*ServiceResponse, error)
	grpc.ServerStream
}

type postServiceRegisterServer struct {
	grpc.ServerStream
}

func (x *postServiceRegisterServer) Send(m *NodeRequest) error {
	return x.ServerStream.SendMsg(m)
}

func (x *postServiceRegisterServer) Recv() (*ServiceResponse, error) {
	m := new(ServiceResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spacemesh.v1.PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Register",
			Handler:       _PostService_Register_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "spacemesh/v1/post.proto",
}

const (
	PostInfoService_PostStates_FullMethodName = "/spacemesh.v1.PostInfoService/PostStates"
)

// PostInfoServiceClient is the client API for PostInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostInfoServiceClient interface {
	// PostStates returns information about the state of the PoST for all known IDs.
	PostStates(ctx context.Context, in *PostStatesRequest, opts ...grpc.CallOption) (*PostStatesResponse, error)
}

type postInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostInfoServiceClient(cc grpc.ClientConnInterface) PostInfoServiceClient {
	return &postInfoServiceClient{cc}
}

func (c *postInfoServiceClient) PostStates(ctx context.Context, in *PostStatesRequest, opts ...grpc.CallOption) (*PostStatesResponse, error) {
	out := new(PostStatesResponse)
	err := c.cc.Invoke(ctx, PostInfoService_PostStates_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostInfoServiceServer is the server API for PostInfoService service.
// All implementations should embed UnimplementedPostInfoServiceServer
// for forward compatibility
type PostInfoServiceServer interface {
	// PostStates returns information about the state of the PoST for all known IDs.
	PostStates(context.Context, *PostStatesRequest) (*PostStatesResponse, error)
}

// UnimplementedPostInfoServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPostInfoServiceServer struct {
}

func (UnimplementedPostInfoServiceServer) PostStates(context.Context, *PostStatesRequest) (*PostStatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostStates not implemented")
}

// UnsafePostInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostInfoServiceServer will
// result in compilation errors.
type UnsafePostInfoServiceServer interface {
	mustEmbedUnimplementedPostInfoServiceServer()
}

func RegisterPostInfoServiceServer(s grpc.ServiceRegistrar, srv PostInfoServiceServer) {
	s.RegisterService(&PostInfoService_ServiceDesc, srv)
}

func _PostInfoService_PostStates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostStatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostInfoServiceServer).PostStates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostInfoService_PostStates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostInfoServiceServer).PostStates(ctx, req.(*PostStatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostInfoService_ServiceDesc is the grpc.ServiceDesc for PostInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spacemesh.v1.PostInfoService",
	HandlerType: (*PostInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostStates",
			Handler:    _PostInfoService_PostStates_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spacemesh/v1/post.proto",
}