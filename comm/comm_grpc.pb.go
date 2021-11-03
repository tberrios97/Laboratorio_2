// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package go_comm_grpc

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

// CommClient is the client API for Comm service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommClient interface {
	FunTest(ctx context.Context, in *RequestTest, opts ...grpc.CallOption) (*ResponseTest, error)
}

type commClient struct {
	cc grpc.ClientConnInterface
}

func NewCommClient(cc grpc.ClientConnInterface) CommClient {
	return &commClient{cc}
}

func (c *commClient) FunTest(ctx context.Context, in *RequestTest, opts ...grpc.CallOption) (*ResponseTest, error) {
	out := new(ResponseTest)
	err := c.cc.Invoke(ctx, "/comm.Comm/FunTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommServer is the server API for Comm service.
// All implementations must embed UnimplementedCommServer
// for forward compatibility
type CommServer interface {
	FunTest(context.Context, *RequestTest) (*ResponseTest, error)
	mustEmbedUnimplementedCommServer()
}

// UnimplementedCommServer must be embedded to have forward compatible implementations.
type UnimplementedCommServer struct {
}

func (UnimplementedCommServer) FunTest(context.Context, *RequestTest) (*ResponseTest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FunTest not implemented")
}
func (UnimplementedCommServer) mustEmbedUnimplementedCommServer() {}

// UnsafeCommServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommServer will
// result in compilation errors.
type UnsafeCommServer interface {
	mustEmbedUnimplementedCommServer()
}

func RegisterCommServer(s grpc.ServiceRegistrar, srv CommServer) {
	s.RegisterService(&Comm_ServiceDesc, srv)
}

func _Comm_FunTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestTest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).FunTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/FunTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).FunTest(ctx, req.(*RequestTest))
	}
	return interceptor(ctx, in, info, handler)
}

// Comm_ServiceDesc is the grpc.ServiceDesc for Comm service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Comm_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "comm.Comm",
	HandlerType: (*CommServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FunTest",
			Handler:    _Comm_FunTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comm/comm.proto",
}
