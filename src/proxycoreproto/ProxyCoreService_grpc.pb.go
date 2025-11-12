





package proxycoreproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)




const _ = grpc.SupportPackageIsVersion9

const (
	ProxyCore_StartCore_FullMethodName     = "/ProxyCore.ProxyCore/startCore"
	ProxyCore_StopCore_FullMethodName      = "/ProxyCore.ProxyCore/stopCore"
	ProxyCore_IsCoreRunning_FullMethodName = "/ProxyCore.ProxyCore/isCoreRunning"
	ProxyCore_GetVersion_FullMethodName    = "/ProxyCore.ProxyCore/getVersion"
	ProxyCore_FetchLogs_FullMethodName     = "/ProxyCore.ProxyCore/fetchLogs"
	ProxyCore_ClearLogs_FullMethodName     = "/ProxyCore.ProxyCore/clearLogs"
	ProxyCore_MeasurePing_FullMethodName   = "/ProxyCore.ProxyCore/measurePing"
)






type ProxyCoreClient interface {
	StartCore(ctx context.Context, in *StartCoreRequest, opts ...grpc.CallOption) (*Empty, error)
	StopCore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	IsCoreRunning(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BooleanResponse, error)
	GetVersion(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*VersionResponse, error)
	FetchLogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LogResponse, error)
	ClearLogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	MeasurePing(ctx context.Context, in *MeasurePingRequest, opts ...grpc.CallOption) (*MeasurePingResponse, error)
}

type proxyCoreClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyCoreClient(cc grpc.ClientConnInterface) ProxyCoreClient {
	return &proxyCoreClient{cc}
}

func (c *proxyCoreClient) StartCore(ctx context.Context, in *StartCoreRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ProxyCore_StartCore_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyCoreClient) StopCore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ProxyCore_StopCore_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyCoreClient) IsCoreRunning(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BooleanResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BooleanResponse)
	err := c.cc.Invoke(ctx, ProxyCore_IsCoreRunning_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyCoreClient) GetVersion(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, ProxyCore_GetVersion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyCoreClient) FetchLogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LogResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogResponse)
	err := c.cc.Invoke(ctx, ProxyCore_FetchLogs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyCoreClient) ClearLogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, ProxyCore_ClearLogs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyCoreClient) MeasurePing(ctx context.Context, in *MeasurePingRequest, opts ...grpc.CallOption) (*MeasurePingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MeasurePingResponse)
	err := c.cc.Invoke(ctx, ProxyCore_MeasurePing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}






type ProxyCoreServer interface {
	StartCore(context.Context, *StartCoreRequest) (*Empty, error)
	StopCore(context.Context, *Empty) (*Empty, error)
	IsCoreRunning(context.Context, *Empty) (*BooleanResponse, error)
	GetVersion(context.Context, *Empty) (*VersionResponse, error)
	FetchLogs(context.Context, *Empty) (*LogResponse, error)
	ClearLogs(context.Context, *Empty) (*Empty, error)
	MeasurePing(context.Context, *MeasurePingRequest) (*MeasurePingResponse, error)
	mustEmbedUnimplementedProxyCoreServer()
}






type UnimplementedProxyCoreServer struct{}

func (UnimplementedProxyCoreServer) StartCore(context.Context, *StartCoreRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartCore not implemented")
}
func (UnimplementedProxyCoreServer) StopCore(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopCore not implemented")
}
func (UnimplementedProxyCoreServer) IsCoreRunning(context.Context, *Empty) (*BooleanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsCoreRunning not implemented")
}
func (UnimplementedProxyCoreServer) GetVersion(context.Context, *Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedProxyCoreServer) FetchLogs(context.Context, *Empty) (*LogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchLogs not implemented")
}
func (UnimplementedProxyCoreServer) ClearLogs(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearLogs not implemented")
}
func (UnimplementedProxyCoreServer) MeasurePing(context.Context, *MeasurePingRequest) (*MeasurePingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MeasurePing not implemented")
}
func (UnimplementedProxyCoreServer) mustEmbedUnimplementedProxyCoreServer() {}
func (UnimplementedProxyCoreServer) testEmbeddedByValue()                   {}




type UnsafeProxyCoreServer interface {
	mustEmbedUnimplementedProxyCoreServer()
}

func RegisterProxyCoreServer(s grpc.ServiceRegistrar, srv ProxyCoreServer) {
	
	
	
	
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProxyCore_ServiceDesc, srv)
}

func _ProxyCore_StartCore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartCoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyCoreServer).StartCore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyCore_StartCore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyCoreServer).StartCore(ctx, req.(*StartCoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyCore_StopCore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyCoreServer).StopCore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyCore_StopCore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyCoreServer).StopCore(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyCore_IsCoreRunning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyCoreServer).IsCoreRunning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyCore_IsCoreRunning_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyCoreServer).IsCoreRunning(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyCore_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyCoreServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyCore_GetVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyCoreServer).GetVersion(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyCore_FetchLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyCoreServer).FetchLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyCore_FetchLogs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyCoreServer).FetchLogs(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyCore_ClearLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyCoreServer).ClearLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyCore_ClearLogs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyCoreServer).ClearLogs(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyCore_MeasurePing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MeasurePingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyCoreServer).MeasurePing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyCore_MeasurePing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyCoreServer).MeasurePing(ctx, req.(*MeasurePingRequest))
	}
	return interceptor(ctx, in, info, handler)
}




var ProxyCore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ProxyCore.ProxyCore",
	HandlerType: (*ProxyCoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "startCore",
			Handler:    _ProxyCore_StartCore_Handler,
		},
		{
			MethodName: "stopCore",
			Handler:    _ProxyCore_StopCore_Handler,
		},
		{
			MethodName: "isCoreRunning",
			Handler:    _ProxyCore_IsCoreRunning_Handler,
		},
		{
			MethodName: "getVersion",
			Handler:    _ProxyCore_GetVersion_Handler,
		},
		{
			MethodName: "fetchLogs",
			Handler:    _ProxyCore_FetchLogs_Handler,
		},
		{
			MethodName: "clearLogs",
			Handler:    _ProxyCore_ClearLogs_Handler,
		},
		{
			MethodName: "measurePing",
			Handler:    _ProxyCore_MeasurePing_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ProxyCoreService.proto",
}
