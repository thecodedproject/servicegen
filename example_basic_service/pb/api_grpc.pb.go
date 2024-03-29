// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.3
// source: api.proto

package pb

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

// BasicClient is the client API for Basic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BasicClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	UseManyTypes(ctx context.Context, in *UseManyTypesRequest, opts ...grpc.CallOption) (*UseManyTypesResponse, error)
}

type basicClient struct {
	cc grpc.ClientConnInterface
}

func NewBasicClient(cc grpc.ClientConnInterface) BasicClient {
	return &basicClient{cc}
}

func (c *basicClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/basicpb.Basic/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, "/basicpb.Basic/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicClient) UseManyTypes(ctx context.Context, in *UseManyTypesRequest, opts ...grpc.CallOption) (*UseManyTypesResponse, error) {
	out := new(UseManyTypesResponse)
	err := c.cc.Invoke(ctx, "/basicpb.Basic/UseManyTypes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BasicServer is the server API for Basic service.
// All implementations must embed UnimplementedBasicServer
// for forward compatibility
type BasicServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	UseManyTypes(context.Context, *UseManyTypesRequest) (*UseManyTypesResponse, error)
	mustEmbedUnimplementedBasicServer()
}

// UnimplementedBasicServer must be embedded to have forward compatible implementations.
type UnimplementedBasicServer struct {
}

func (UnimplementedBasicServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedBasicServer) Echo(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedBasicServer) UseManyTypes(context.Context, *UseManyTypesRequest) (*UseManyTypesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseManyTypes not implemented")
}
func (UnimplementedBasicServer) mustEmbedUnimplementedBasicServer() {}

// UnsafeBasicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BasicServer will
// result in compilation errors.
type UnsafeBasicServer interface {
	mustEmbedUnimplementedBasicServer()
}

func RegisterBasicServer(s grpc.ServiceRegistrar, srv BasicServer) {
	s.RegisterService(&Basic_ServiceDesc, srv)
}

func _Basic_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/basicpb.Basic/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Basic_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/basicpb.Basic/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Basic_UseManyTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseManyTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicServer).UseManyTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/basicpb.Basic/UseManyTypes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicServer).UseManyTypes(ctx, req.(*UseManyTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Basic_ServiceDesc is the grpc.ServiceDesc for Basic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Basic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "basicpb.Basic",
	HandlerType: (*BasicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Basic_Ping_Handler,
		},
		{
			MethodName: "Echo",
			Handler:    _Basic_Echo_Handler,
		},
		{
			MethodName: "UseManyTypes",
			Handler:    _Basic_UseManyTypes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
