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

// NestedClient is the client API for Nested service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NestedClient interface {
	LoneNest(ctx context.Context, in *LoneNestRequest, opts ...grpc.CallOption) (*LoneNestResponse, error)
	SeveralParams(ctx context.Context, in *SeveralParamsRequest, opts ...grpc.CallOption) (*SeveralParamsResponse, error)
	MultiDepth(ctx context.Context, in *MultiDepthRequest, opts ...grpc.CallOption) (*MultiDepthResponse, error)
}

type nestedClient struct {
	cc grpc.ClientConnInterface
}

func NewNestedClient(cc grpc.ClientConnInterface) NestedClient {
	return &nestedClient{cc}
}

func (c *nestedClient) LoneNest(ctx context.Context, in *LoneNestRequest, opts ...grpc.CallOption) (*LoneNestResponse, error) {
	out := new(LoneNestResponse)
	err := c.cc.Invoke(ctx, "/pb.Nested/LoneNest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nestedClient) SeveralParams(ctx context.Context, in *SeveralParamsRequest, opts ...grpc.CallOption) (*SeveralParamsResponse, error) {
	out := new(SeveralParamsResponse)
	err := c.cc.Invoke(ctx, "/pb.Nested/SeveralParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nestedClient) MultiDepth(ctx context.Context, in *MultiDepthRequest, opts ...grpc.CallOption) (*MultiDepthResponse, error) {
	out := new(MultiDepthResponse)
	err := c.cc.Invoke(ctx, "/pb.Nested/MultiDepth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NestedServer is the server API for Nested service.
// All implementations must embed UnimplementedNestedServer
// for forward compatibility
type NestedServer interface {
	LoneNest(context.Context, *LoneNestRequest) (*LoneNestResponse, error)
	SeveralParams(context.Context, *SeveralParamsRequest) (*SeveralParamsResponse, error)
	MultiDepth(context.Context, *MultiDepthRequest) (*MultiDepthResponse, error)
	mustEmbedUnimplementedNestedServer()
}

// UnimplementedNestedServer must be embedded to have forward compatible implementations.
type UnimplementedNestedServer struct {
}

func (UnimplementedNestedServer) LoneNest(context.Context, *LoneNestRequest) (*LoneNestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoneNest not implemented")
}
func (UnimplementedNestedServer) SeveralParams(context.Context, *SeveralParamsRequest) (*SeveralParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SeveralParams not implemented")
}
func (UnimplementedNestedServer) MultiDepth(context.Context, *MultiDepthRequest) (*MultiDepthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiDepth not implemented")
}
func (UnimplementedNestedServer) mustEmbedUnimplementedNestedServer() {}

// UnsafeNestedServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NestedServer will
// result in compilation errors.
type UnsafeNestedServer interface {
	mustEmbedUnimplementedNestedServer()
}

func RegisterNestedServer(s grpc.ServiceRegistrar, srv NestedServer) {
	s.RegisterService(&Nested_ServiceDesc, srv)
}

func _Nested_LoneNest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoneNestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NestedServer).LoneNest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Nested/LoneNest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NestedServer).LoneNest(ctx, req.(*LoneNestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nested_SeveralParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SeveralParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NestedServer).SeveralParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Nested/SeveralParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NestedServer).SeveralParams(ctx, req.(*SeveralParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nested_MultiDepth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiDepthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NestedServer).MultiDepth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Nested/MultiDepth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NestedServer).MultiDepth(ctx, req.(*MultiDepthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Nested_ServiceDesc is the grpc.ServiceDesc for Nested service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Nested_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Nested",
	HandlerType: (*NestedServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoneNest",
			Handler:    _Nested_LoneNest_Handler,
		},
		{
			MethodName: "SeveralParams",
			Handler:    _Nested_SeveralParams_Handler,
		},
		{
			MethodName: "MultiDepth",
			Handler:    _Nested_MultiDepth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}