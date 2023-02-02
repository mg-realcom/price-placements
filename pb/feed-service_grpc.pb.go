// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: feed-service.proto

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

// FeedServiceClient is the client API for FeedService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FeedServiceClient interface {
	CheckPhonesAll(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	CheckPhonesRealty(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	CheckPhonesCian(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	CheckPhonesAvito(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	CheckPhonesDomclick(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	ValidateFeed(ctx context.Context, in *ValidateFeedRequest, opts ...grpc.CallOption) (*ValidateFeedResponse, error)
	ValidateFeedAll(ctx context.Context, in *ValidateFeedAllRequest, opts ...grpc.CallOption) (*ValidateFeedAllResponse, error)
}

type feedServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedServiceClient(cc grpc.ClientConnInterface) FeedServiceClient {
	return &feedServiceClient{cc}
}

func (c *feedServiceClient) CheckPhonesAll(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/FeedService/CheckPhonesAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) CheckPhonesRealty(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/FeedService/CheckPhonesRealty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) CheckPhonesCian(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/FeedService/CheckPhonesCian", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) CheckPhonesAvito(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/FeedService/CheckPhonesAvito", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) CheckPhonesDomclick(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/FeedService/CheckPhonesDomclick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) ValidateFeed(ctx context.Context, in *ValidateFeedRequest, opts ...grpc.CallOption) (*ValidateFeedResponse, error) {
	out := new(ValidateFeedResponse)
	err := c.cc.Invoke(ctx, "/FeedService/ValidateFeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) ValidateFeedAll(ctx context.Context, in *ValidateFeedAllRequest, opts ...grpc.CallOption) (*ValidateFeedAllResponse, error) {
	out := new(ValidateFeedAllResponse)
	err := c.cc.Invoke(ctx, "/FeedService/ValidateFeedAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedServiceServer is the server API for FeedService service.
// All implementations must embed UnimplementedFeedServiceServer
// for forward compatibility
type FeedServiceServer interface {
	CheckPhonesAll(context.Context, *CheckRequest) (*CheckResponse, error)
	CheckPhonesRealty(context.Context, *CheckRequest) (*CheckResponse, error)
	CheckPhonesCian(context.Context, *CheckRequest) (*CheckResponse, error)
	CheckPhonesAvito(context.Context, *CheckRequest) (*CheckResponse, error)
	CheckPhonesDomclick(context.Context, *CheckRequest) (*CheckResponse, error)
	ValidateFeed(context.Context, *ValidateFeedRequest) (*ValidateFeedResponse, error)
	ValidateFeedAll(context.Context, *ValidateFeedAllRequest) (*ValidateFeedAllResponse, error)
	mustEmbedUnimplementedFeedServiceServer()
}

// UnimplementedFeedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFeedServiceServer struct {
}

func (UnimplementedFeedServiceServer) CheckPhonesAll(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPhonesAll not implemented")
}
func (UnimplementedFeedServiceServer) CheckPhonesRealty(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPhonesRealty not implemented")
}
func (UnimplementedFeedServiceServer) CheckPhonesCian(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPhonesCian not implemented")
}
func (UnimplementedFeedServiceServer) CheckPhonesAvito(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPhonesAvito not implemented")
}
func (UnimplementedFeedServiceServer) CheckPhonesDomclick(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPhonesDomclick not implemented")
}
func (UnimplementedFeedServiceServer) ValidateFeed(context.Context, *ValidateFeedRequest) (*ValidateFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateFeed not implemented")
}
func (UnimplementedFeedServiceServer) ValidateFeedAll(context.Context, *ValidateFeedAllRequest) (*ValidateFeedAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateFeedAll not implemented")
}
func (UnimplementedFeedServiceServer) mustEmbedUnimplementedFeedServiceServer() {}

// UnsafeFeedServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FeedServiceServer will
// result in compilation errors.
type UnsafeFeedServiceServer interface {
	mustEmbedUnimplementedFeedServiceServer()
}

func RegisterFeedServiceServer(s grpc.ServiceRegistrar, srv FeedServiceServer) {
	s.RegisterService(&FeedService_ServiceDesc, srv)
}

func _FeedService_CheckPhonesAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).CheckPhonesAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FeedService/CheckPhonesAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).CheckPhonesAll(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_CheckPhonesRealty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).CheckPhonesRealty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FeedService/CheckPhonesRealty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).CheckPhonesRealty(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_CheckPhonesCian_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).CheckPhonesCian(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FeedService/CheckPhonesCian",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).CheckPhonesCian(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_CheckPhonesAvito_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).CheckPhonesAvito(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FeedService/CheckPhonesAvito",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).CheckPhonesAvito(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_CheckPhonesDomclick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).CheckPhonesDomclick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FeedService/CheckPhonesDomclick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).CheckPhonesDomclick(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_ValidateFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).ValidateFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FeedService/ValidateFeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).ValidateFeed(ctx, req.(*ValidateFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_ValidateFeedAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateFeedAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).ValidateFeedAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FeedService/ValidateFeedAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).ValidateFeedAll(ctx, req.(*ValidateFeedAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FeedService_ServiceDesc is the grpc.ServiceDesc for FeedService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FeedService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FeedService",
	HandlerType: (*FeedServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckPhonesAll",
			Handler:    _FeedService_CheckPhonesAll_Handler,
		},
		{
			MethodName: "CheckPhonesRealty",
			Handler:    _FeedService_CheckPhonesRealty_Handler,
		},
		{
			MethodName: "CheckPhonesCian",
			Handler:    _FeedService_CheckPhonesCian_Handler,
		},
		{
			MethodName: "CheckPhonesAvito",
			Handler:    _FeedService_CheckPhonesAvito_Handler,
		},
		{
			MethodName: "CheckPhonesDomclick",
			Handler:    _FeedService_CheckPhonesDomclick_Handler,
		},
		{
			MethodName: "ValidateFeed",
			Handler:    _FeedService_ValidateFeed_Handler,
		},
		{
			MethodName: "ValidateFeedAll",
			Handler:    _FeedService_ValidateFeedAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "feed-service.proto",
}
