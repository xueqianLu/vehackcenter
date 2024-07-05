// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.1
// source: hackcenter/center.proto

package hackcenter

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

// CenterServiceClient is the client API for CenterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CenterServiceClient interface {
	SubmitBlock(ctx context.Context, in *Block, opts ...grpc.CallOption) (*SubmitBlockResponse, error)
	SubscribeBlock(ctx context.Context, in *SubscribeBlockRequest, opts ...grpc.CallOption) (CenterService_SubscribeBlockClient, error)
	SubBroadcastTask(ctx context.Context, in *SubBroadcastTaskRequest, opts ...grpc.CallOption) (CenterService_SubBroadcastTaskClient, error)
	BeginToHack(ctx context.Context, in *BeginToHackRequest, opts ...grpc.CallOption) (*BeginToHackResponse, error)
}

type centerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCenterServiceClient(cc grpc.ClientConnInterface) CenterServiceClient {
	return &centerServiceClient{cc}
}

func (c *centerServiceClient) SubmitBlock(ctx context.Context, in *Block, opts ...grpc.CallOption) (*SubmitBlockResponse, error) {
	out := new(SubmitBlockResponse)
	err := c.cc.Invoke(ctx, "/hackcenter.CenterService/SubmitBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centerServiceClient) SubscribeBlock(ctx context.Context, in *SubscribeBlockRequest, opts ...grpc.CallOption) (CenterService_SubscribeBlockClient, error) {
	stream, err := c.cc.NewStream(ctx, &CenterService_ServiceDesc.Streams[0], "/hackcenter.CenterService/SubscribeBlock", opts...)
	if err != nil {
		return nil, err
	}
	x := &centerServiceSubscribeBlockClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CenterService_SubscribeBlockClient interface {
	Recv() (*Block, error)
	grpc.ClientStream
}

type centerServiceSubscribeBlockClient struct {
	grpc.ClientStream
}

func (x *centerServiceSubscribeBlockClient) Recv() (*Block, error) {
	m := new(Block)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *centerServiceClient) SubBroadcastTask(ctx context.Context, in *SubBroadcastTaskRequest, opts ...grpc.CallOption) (CenterService_SubBroadcastTaskClient, error) {
	stream, err := c.cc.NewStream(ctx, &CenterService_ServiceDesc.Streams[1], "/hackcenter.CenterService/SubBroadcastTask", opts...)
	if err != nil {
		return nil, err
	}
	x := &centerServiceSubBroadcastTaskClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CenterService_SubBroadcastTaskClient interface {
	Recv() (*Block, error)
	grpc.ClientStream
}

type centerServiceSubBroadcastTaskClient struct {
	grpc.ClientStream
}

func (x *centerServiceSubBroadcastTaskClient) Recv() (*Block, error) {
	m := new(Block)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *centerServiceClient) BeginToHack(ctx context.Context, in *BeginToHackRequest, opts ...grpc.CallOption) (*BeginToHackResponse, error) {
	out := new(BeginToHackResponse)
	err := c.cc.Invoke(ctx, "/hackcenter.CenterService/BeginToHack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CenterServiceServer is the server API for CenterService service.
// All implementations must embed UnimplementedCenterServiceServer
// for forward compatibility
type CenterServiceServer interface {
	SubmitBlock(context.Context, *Block) (*SubmitBlockResponse, error)
	SubscribeBlock(*SubscribeBlockRequest, CenterService_SubscribeBlockServer) error
	SubBroadcastTask(*SubBroadcastTaskRequest, CenterService_SubBroadcastTaskServer) error
	BeginToHack(context.Context, *BeginToHackRequest) (*BeginToHackResponse, error)
	mustEmbedUnimplementedCenterServiceServer()
}

// UnimplementedCenterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCenterServiceServer struct {
}

func (UnimplementedCenterServiceServer) SubmitBlock(context.Context, *Block) (*SubmitBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitBlock not implemented")
}
func (UnimplementedCenterServiceServer) SubscribeBlock(*SubscribeBlockRequest, CenterService_SubscribeBlockServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeBlock not implemented")
}
func (UnimplementedCenterServiceServer) SubBroadcastTask(*SubBroadcastTaskRequest, CenterService_SubBroadcastTaskServer) error {
	return status.Errorf(codes.Unimplemented, "method SubBroadcastTask not implemented")
}
func (UnimplementedCenterServiceServer) BeginToHack(context.Context, *BeginToHackRequest) (*BeginToHackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginToHack not implemented")
}
func (UnimplementedCenterServiceServer) mustEmbedUnimplementedCenterServiceServer() {}

// UnsafeCenterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CenterServiceServer will
// result in compilation errors.
type UnsafeCenterServiceServer interface {
	mustEmbedUnimplementedCenterServiceServer()
}

func RegisterCenterServiceServer(s grpc.ServiceRegistrar, srv CenterServiceServer) {
	s.RegisterService(&CenterService_ServiceDesc, srv)
}

func _CenterService_SubmitBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Block)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CenterServiceServer).SubmitBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hackcenter.CenterService/SubmitBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CenterServiceServer).SubmitBlock(ctx, req.(*Block))
	}
	return interceptor(ctx, in, info, handler)
}

func _CenterService_SubscribeBlock_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeBlockRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CenterServiceServer).SubscribeBlock(m, &centerServiceSubscribeBlockServer{stream})
}

type CenterService_SubscribeBlockServer interface {
	Send(*Block) error
	grpc.ServerStream
}

type centerServiceSubscribeBlockServer struct {
	grpc.ServerStream
}

func (x *centerServiceSubscribeBlockServer) Send(m *Block) error {
	return x.ServerStream.SendMsg(m)
}

func _CenterService_SubBroadcastTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubBroadcastTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CenterServiceServer).SubBroadcastTask(m, &centerServiceSubBroadcastTaskServer{stream})
}

type CenterService_SubBroadcastTaskServer interface {
	Send(*Block) error
	grpc.ServerStream
}

type centerServiceSubBroadcastTaskServer struct {
	grpc.ServerStream
}

func (x *centerServiceSubBroadcastTaskServer) Send(m *Block) error {
	return x.ServerStream.SendMsg(m)
}

func _CenterService_BeginToHack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeginToHackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CenterServiceServer).BeginToHack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hackcenter.CenterService/BeginToHack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CenterServiceServer).BeginToHack(ctx, req.(*BeginToHackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CenterService_ServiceDesc is the grpc.ServiceDesc for CenterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CenterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hackcenter.CenterService",
	HandlerType: (*CenterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitBlock",
			Handler:    _CenterService_SubmitBlock_Handler,
		},
		{
			MethodName: "BeginToHack",
			Handler:    _CenterService_BeginToHack_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeBlock",
			Handler:       _CenterService_SubscribeBlock_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubBroadcastTask",
			Handler:       _CenterService_SubBroadcastTask_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "hackcenter/center.proto",
}
