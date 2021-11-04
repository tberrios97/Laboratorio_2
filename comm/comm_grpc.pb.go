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
	UnirseJuegoCalamar(ctx context.Context, in *RequestUnirse, opts ...grpc.CallOption) (*ResponseUnirse, error)
	InicioEtapa(ctx context.Context, in *RequestInicio, opts ...grpc.CallOption) (*ResponseInicio, error)
	JugadaPrimeraEtapa(ctx context.Context, in *RequestPrimeraEtapa, opts ...grpc.CallOption) (*ResponsePrimeraEtapa, error)
	JugadaSegundaEtapa(ctx context.Context, in *RequestSegundaEtapa, opts ...grpc.CallOption) (*ResponseSegundaEtapa, error)
	JugadaTerceraEtapa(ctx context.Context, in *RequestTerceraEtapa, opts ...grpc.CallOption) (*ResponseTerceraEtapa, error)
	RegistrarJugadaJugador(ctx context.Context, in *RequestRJJ, opts ...grpc.CallOption) (*ResponseRJJ, error)
	RegistrarJugadaDN(ctx context.Context, in *RequestRJDN, opts ...grpc.CallOption) (*ResponseRJDN, error)
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

func (c *commClient) UnirseJuegoCalamar(ctx context.Context, in *RequestUnirse, opts ...grpc.CallOption) (*ResponseUnirse, error) {
	out := new(ResponseUnirse)
	err := c.cc.Invoke(ctx, "/comm.Comm/UnirseJuegoCalamar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commClient) InicioEtapa(ctx context.Context, in *RequestInicio, opts ...grpc.CallOption) (*ResponseInicio, error) {
	out := new(ResponseInicio)
	err := c.cc.Invoke(ctx, "/comm.Comm/InicioEtapa", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commClient) JugadaPrimeraEtapa(ctx context.Context, in *RequestPrimeraEtapa, opts ...grpc.CallOption) (*ResponsePrimeraEtapa, error) {
	out := new(ResponsePrimeraEtapa)
	err := c.cc.Invoke(ctx, "/comm.Comm/JugadaPrimeraEtapa", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commClient) JugadaSegundaEtapa(ctx context.Context, in *RequestSegundaEtapa, opts ...grpc.CallOption) (*ResponseSegundaEtapa, error) {
	out := new(ResponseSegundaEtapa)
	err := c.cc.Invoke(ctx, "/comm.Comm/JugadaSegundaEtapa", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commClient) JugadaTerceraEtapa(ctx context.Context, in *RequestTerceraEtapa, opts ...grpc.CallOption) (*ResponseTerceraEtapa, error) {
	out := new(ResponseTerceraEtapa)
	err := c.cc.Invoke(ctx, "/comm.Comm/JugadaTerceraEtapa", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commClient) RegistrarJugadaJugador(ctx context.Context, in *RequestRJJ, opts ...grpc.CallOption) (*ResponseRJJ, error) {
	out := new(ResponseRJJ)
	err := c.cc.Invoke(ctx, "/comm.Comm/RegistrarJugadaJugador", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commClient) RegistrarJugadaDN(ctx context.Context, in *RequestRJDN, opts ...grpc.CallOption) (*ResponseRJDN, error) {
	out := new(ResponseRJDN)
	err := c.cc.Invoke(ctx, "/comm.Comm/RegistrarJugadaDN", in, out, opts...)
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
	UnirseJuegoCalamar(context.Context, *RequestUnirse) (*ResponseUnirse, error)
	InicioEtapa(context.Context, *RequestInicio) (*ResponseInicio, error)
	JugadaPrimeraEtapa(context.Context, *RequestPrimeraEtapa) (*ResponsePrimeraEtapa, error)
	JugadaSegundaEtapa(context.Context, *RequestSegundaEtapa) (*ResponseSegundaEtapa, error)
	JugadaTerceraEtapa(context.Context, *RequestTerceraEtapa) (*ResponseTerceraEtapa, error)
	RegistrarJugadaJugador(context.Context, *RequestRJJ) (*ResponseRJJ, error)
	RegistrarJugadaDN(context.Context, *RequestRJDN) (*ResponseRJDN, error)
	mustEmbedUnimplementedCommServer()
}

// UnimplementedCommServer must be embedded to have forward compatible implementations.
type UnimplementedCommServer struct {
}

func (UnimplementedCommServer) FunTest(context.Context, *RequestTest) (*ResponseTest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FunTest not implemented")
}
func (UnimplementedCommServer) UnirseJuegoCalamar(context.Context, *RequestUnirse) (*ResponseUnirse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnirseJuegoCalamar not implemented")
}
func (UnimplementedCommServer) InicioEtapa(context.Context, *RequestInicio) (*ResponseInicio, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InicioEtapa not implemented")
}
func (UnimplementedCommServer) JugadaPrimeraEtapa(context.Context, *RequestPrimeraEtapa) (*ResponsePrimeraEtapa, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JugadaPrimeraEtapa not implemented")
}
func (UnimplementedCommServer) JugadaSegundaEtapa(context.Context, *RequestSegundaEtapa) (*ResponseSegundaEtapa, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JugadaSegundaEtapa not implemented")
}
func (UnimplementedCommServer) JugadaTerceraEtapa(context.Context, *RequestTerceraEtapa) (*ResponseTerceraEtapa, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JugadaTerceraEtapa not implemented")
}
func (UnimplementedCommServer) RegistrarJugadaJugador(context.Context, *RequestRJJ) (*ResponseRJJ, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegistrarJugadaJugador not implemented")
}
func (UnimplementedCommServer) RegistrarJugadaDN(context.Context, *RequestRJDN) (*ResponseRJDN, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegistrarJugadaDN not implemented")
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

func _Comm_UnirseJuegoCalamar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUnirse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).UnirseJuegoCalamar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/UnirseJuegoCalamar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).UnirseJuegoCalamar(ctx, req.(*RequestUnirse))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comm_InicioEtapa_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInicio)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).InicioEtapa(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/InicioEtapa",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).InicioEtapa(ctx, req.(*RequestInicio))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comm_JugadaPrimeraEtapa_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestPrimeraEtapa)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).JugadaPrimeraEtapa(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/JugadaPrimeraEtapa",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).JugadaPrimeraEtapa(ctx, req.(*RequestPrimeraEtapa))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comm_JugadaSegundaEtapa_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestSegundaEtapa)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).JugadaSegundaEtapa(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/JugadaSegundaEtapa",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).JugadaSegundaEtapa(ctx, req.(*RequestSegundaEtapa))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comm_JugadaTerceraEtapa_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestTerceraEtapa)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).JugadaTerceraEtapa(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/JugadaTerceraEtapa",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).JugadaTerceraEtapa(ctx, req.(*RequestTerceraEtapa))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comm_RegistrarJugadaJugador_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestRJJ)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).RegistrarJugadaJugador(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/RegistrarJugadaJugador",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).RegistrarJugadaJugador(ctx, req.(*RequestRJJ))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comm_RegistrarJugadaDN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestRJDN)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommServer).RegistrarJugadaDN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/comm.Comm/RegistrarJugadaDN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommServer).RegistrarJugadaDN(ctx, req.(*RequestRJDN))
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
		{
			MethodName: "UnirseJuegoCalamar",
			Handler:    _Comm_UnirseJuegoCalamar_Handler,
		},
		{
			MethodName: "InicioEtapa",
			Handler:    _Comm_InicioEtapa_Handler,
		},
		{
			MethodName: "JugadaPrimeraEtapa",
			Handler:    _Comm_JugadaPrimeraEtapa_Handler,
		},
		{
			MethodName: "JugadaSegundaEtapa",
			Handler:    _Comm_JugadaSegundaEtapa_Handler,
		},
		{
			MethodName: "JugadaTerceraEtapa",
			Handler:    _Comm_JugadaTerceraEtapa_Handler,
		},
		{
			MethodName: "RegistrarJugadaJugador",
			Handler:    _Comm_RegistrarJugadaJugador_Handler,
		},
		{
			MethodName: "RegistrarJugadaDN",
			Handler:    _Comm_RegistrarJugadaDN_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comm/comm.proto",
}
