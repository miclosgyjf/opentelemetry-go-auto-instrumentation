// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: grpc.proto

package pkgs

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
	HelloGrpc_Hello_FullMethodName = "/HelloGrpc/Hello"
)

// HelloGrpcClient is the client API for HelloGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloGrpcClient interface {
	Hello(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error)
}

type helloGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloGrpcClient(cc grpc.ClientConnInterface) HelloGrpcClient {
	return &helloGrpcClient{cc}
}

func (c *helloGrpcClient) Hello(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error) {
	out := new(Resp)
	err := c.cc.Invoke(ctx, HelloGrpc_Hello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HelloGrpcServer is the server API for HelloGrpc service.
// All implementations must embed UnimplementedHelloGrpcServer
// for forward compatibility
type HelloGrpcServer interface {
	Hello(context.Context, *Req) (*Resp, error)
	mustEmbedUnimplementedHelloGrpcServer()
}

// UnimplementedHelloGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedHelloGrpcServer struct {
}

func (UnimplementedHelloGrpcServer) Hello(context.Context, *Req) (*Resp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedHelloGrpcServer) mustEmbedUnimplementedHelloGrpcServer() {}

// UnsafeHelloGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloGrpcServer will
// result in compilation errors.
type UnsafeHelloGrpcServer interface {
	mustEmbedUnimplementedHelloGrpcServer()
}

func RegisterHelloGrpcServer(s grpc.ServiceRegistrar, srv HelloGrpcServer) {
	s.RegisterService(&HelloGrpc_ServiceDesc, srv)
}

func _HelloGrpc_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloGrpcServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HelloGrpc_Hello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloGrpcServer).Hello(ctx, req.(*Req))
	}
	return interceptor(ctx, in, info, handler)
}

// HelloGrpc_ServiceDesc is the grpc.ServiceDesc for HelloGrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HelloGrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "HelloGrpc",
	HandlerType: (*HelloGrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _HelloGrpc_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc.proto",
}