package pb

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CloudIdeServiceClient interface {
	// CreateSpace 创建云 IDE 空间并等待pod状态变为Running ，第一次创建需要挂载存储卷
	CreateSpace(ctx context.Context, in *WorkspaceInfo, opts ...grpc.CallOption) (*WorkspaceRunningInfo, error)
	// 启动（创建） 云 IDE 空间 非第一次创建，无需挂载储存卷，使用之前的储存卷
	StartSpace(ctx context.Context, in *WorkspaceInfo, opts ...grpc.CallOption) (*WorkspaceRunningInfo, error)
	// 删除云 IDE 空间，需要删除存储卷
	DeleteSpace(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*Response, error)
	// 暂停(删除) 云空间 无需删除储存卷
	StopSpace(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*Response, error)
	// 获取Pod 运行状态
	GetPodSpaceStatus(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*WorkspaceStatus, error)
	// 获取云IDE空间的信息
	GetPodSpaceInfo(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*WorkspaceRunningInfo, error)
}

type cloudIdeServiceClient struct {
	cc grpc.ClientConnInterface
}

func (c *cloudIdeServiceClient) CreateSpace(ctx context.Context, in *WorkspaceInfo, opts ...grpc.CallOption) (*WorkspaceRunningInfo, error) {
	out := &WorkspaceRunningInfo{}
	err := c.cc.Invoke(ctx, "/pb.CloudIdeService/createSpace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudIdeServiceClient) StartSpace(ctx context.Context, in *WorkspaceInfo, opts ...grpc.CallOption) (*WorkspaceRunningInfo, error) {
	out := &WorkspaceRunningInfo{}
	err := c.cc.Invoke(ctx, "/pb.CloudIdeService/startSpace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudIdeServiceClient) DeleteSpace(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/pb.CloudIdeService/deleteSpace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudIdeServiceClient) StopSpace(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/pb.CloudIdeService/stopSpace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudIdeServiceClient) GetPodSpaceStatus(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*WorkspaceStatus, error) {
	out := new(WorkspaceStatus)
	err := c.cc.Invoke(ctx, "/pb.CloudIdeService/getPodSpaceStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudIdeServiceClient) GetPodSpaceInfo(ctx context.Context, in *QueryOption, opts ...grpc.CallOption) (*WorkspaceRunningInfo, error) {
	out := new(WorkspaceRunningInfo)
	err := c.cc.Invoke(ctx, "/pb.CloudIdeService/getPodSpaceInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewCloudIdeServiceClient(cc grpc.ClientConnInterface) CloudIdeServiceClient {
	return &cloudIdeServiceClient{cc}
}

// CloudIdeServiceServer is the server API for CloudIdeService service.
type CloudIdeServiceServer interface {
	// 创建云IDE空间并等待Pod状态变为Running,第一次创建,需要挂载存储卷
	CreateSpace(context.Context, *WorkspaceInfo) (*WorkspaceRunningInfo, error)
	// 启动(创建)云IDE空间,非第一次创建,无需挂载存储卷,使用之前的存储卷
	StartSpace(context.Context, *WorkspaceInfo) (*WorkspaceRunningInfo, error)
	// 删除云IDE空间,需要删除存储卷
	DeleteSpace(context.Context, *QueryOption) (*Response, error)
	// 停止(删除)云工作空间,无需删除存储卷
	StopSpace(context.Context, *QueryOption) (*Response, error)
	// 获取Pod运行状态
	GetPodSpaceStatus(context.Context, *QueryOption) (*WorkspaceStatus, error)
	// 获取云IDE空间Pod的信息
	GetPodSpaceInfo(context.Context, *QueryOption) (*WorkspaceRunningInfo, error)
}

// UnimplementedCloudIdeServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCloudIdeServiceServer struct {
}

func (*UnimplementedCloudIdeServiceServer) CreateSpace(context.Context, *WorkspaceInfo) (*WorkspaceRunningInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSpace not implemented")
}
func (*UnimplementedCloudIdeServiceServer) StartSpace(context.Context, *WorkspaceInfo) (*WorkspaceRunningInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartSpace not implemented")
}
func (*UnimplementedCloudIdeServiceServer) DeleteSpace(context.Context, *QueryOption) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSpace not implemented")
}
func (*UnimplementedCloudIdeServiceServer) StopSpace(context.Context, *QueryOption) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopSpace not implemented")
}
func (*UnimplementedCloudIdeServiceServer) GetPodSpaceStatus(context.Context, *QueryOption) (*WorkspaceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPodSpaceStatus not implemented")
}
func (*UnimplementedCloudIdeServiceServer) GetPodSpaceInfo(context.Context, *QueryOption) (*WorkspaceRunningInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPodSpaceInfo not implemented")
}

func RegisterCloudIdeServiceServer(s *grpc.Server, srv CloudIdeServiceServer) {
	s.RegisterService(&_CloudIdeService_serviceDesc, srv)
}

func _CloudIdeService_CreateSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkspaceInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudIdeServiceServer).CreateSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CloudIdeService/CreateSpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudIdeServiceServer).CreateSpace(ctx, req.(*WorkspaceInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudIdeService_StartSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkspaceInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudIdeServiceServer).StartSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CloudIdeService/StartSpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudIdeServiceServer).StartSpace(ctx, req.(*WorkspaceInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudIdeService_DeleteSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudIdeServiceServer).DeleteSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CloudIdeService/DeleteSpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudIdeServiceServer).DeleteSpace(ctx, req.(*QueryOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudIdeService_StopSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudIdeServiceServer).StopSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CloudIdeService/StopSpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudIdeServiceServer).StopSpace(ctx, req.(*QueryOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudIdeService_GetPodSpaceStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudIdeServiceServer).GetPodSpaceStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CloudIdeService/GetPodSpaceStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudIdeServiceServer).GetPodSpaceStatus(ctx, req.(*QueryOption))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudIdeService_GetPodSpaceInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOption)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudIdeServiceServer).GetPodSpaceInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CloudIdeService/GetPodSpaceInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudIdeServiceServer).GetPodSpaceInfo(ctx, req.(*QueryOption))
	}
	return interceptor(ctx, in, info, handler)
}

var _CloudIdeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.CloudIdeService",
	HandlerType: (*CloudIdeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createSpace",
			Handler:    _CloudIdeService_CreateSpace_Handler,
		},
		{
			MethodName: "startSpace",
			Handler:    _CloudIdeService_StartSpace_Handler,
		},
		{
			MethodName: "deleteSpace",
			Handler:    _CloudIdeService_DeleteSpace_Handler,
		},
		{
			MethodName: "stopSpace",
			Handler:    _CloudIdeService_StopSpace_Handler,
		},
		{
			MethodName: "getPodSpaceStatus",
			Handler:    _CloudIdeService_GetPodSpaceStatus_Handler,
		},
		{
			MethodName: "getPodSpaceInfo",
			Handler:    _CloudIdeService_GetPodSpaceInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/proto/service.proto",
}
