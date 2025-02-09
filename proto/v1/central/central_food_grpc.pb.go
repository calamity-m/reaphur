// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: proto/v1/central/central_food.proto

package centralproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	CentralFoodService_CreateFoodRecord_FullMethodName = "/centralproto.v1.CentralFoodService/CreateFoodRecord"
	CentralFoodService_GetFoodRecords_FullMethodName   = "/centralproto.v1.CentralFoodService/GetFoodRecords"
)

// CentralFoodServiceClient is the client API for CentralFoodService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CentralFoodServiceClient interface {
	// Simple RPC
	//
	// Create some food record in the food diary/journal
	CreateFoodRecord(ctx context.Context, in *CreateFoodRecordRequest, opts ...grpc.CallOption) (*CreateFoodRecordResponse, error)
	// Simple RPC
	//
	// Fetch some food records from the food diary/journal
	GetFoodRecords(ctx context.Context, in *GetFoodRecordsRequest, opts ...grpc.CallOption) (*GetFoodRecordsResponse, error)
}

type centralFoodServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCentralFoodServiceClient(cc grpc.ClientConnInterface) CentralFoodServiceClient {
	return &centralFoodServiceClient{cc}
}

func (c *centralFoodServiceClient) CreateFoodRecord(ctx context.Context, in *CreateFoodRecordRequest, opts ...grpc.CallOption) (*CreateFoodRecordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFoodRecordResponse)
	err := c.cc.Invoke(ctx, CentralFoodService_CreateFoodRecord_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralFoodServiceClient) GetFoodRecords(ctx context.Context, in *GetFoodRecordsRequest, opts ...grpc.CallOption) (*GetFoodRecordsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFoodRecordsResponse)
	err := c.cc.Invoke(ctx, CentralFoodService_GetFoodRecords_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CentralFoodServiceServer is the server API for CentralFoodService service.
// All implementations must embed UnimplementedCentralFoodServiceServer
// for forward compatibility.
type CentralFoodServiceServer interface {
	// Simple RPC
	//
	// Create some food record in the food diary/journal
	CreateFoodRecord(context.Context, *CreateFoodRecordRequest) (*CreateFoodRecordResponse, error)
	// Simple RPC
	//
	// Fetch some food records from the food diary/journal
	GetFoodRecords(context.Context, *GetFoodRecordsRequest) (*GetFoodRecordsResponse, error)
	mustEmbedUnimplementedCentralFoodServiceServer()
}

// UnimplementedCentralFoodServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCentralFoodServiceServer struct{}

func (UnimplementedCentralFoodServiceServer) CreateFoodRecord(context.Context, *CreateFoodRecordRequest) (*CreateFoodRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFoodRecord not implemented")
}
func (UnimplementedCentralFoodServiceServer) GetFoodRecords(context.Context, *GetFoodRecordsRequest) (*GetFoodRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFoodRecords not implemented")
}
func (UnimplementedCentralFoodServiceServer) mustEmbedUnimplementedCentralFoodServiceServer() {}
func (UnimplementedCentralFoodServiceServer) testEmbeddedByValue()                            {}

// UnsafeCentralFoodServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CentralFoodServiceServer will
// result in compilation errors.
type UnsafeCentralFoodServiceServer interface {
	mustEmbedUnimplementedCentralFoodServiceServer()
}

func RegisterCentralFoodServiceServer(s grpc.ServiceRegistrar, srv CentralFoodServiceServer) {
	// If the following call pancis, it indicates UnimplementedCentralFoodServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CentralFoodService_ServiceDesc, srv)
}

func _CentralFoodService_CreateFoodRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFoodRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralFoodServiceServer).CreateFoodRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentralFoodService_CreateFoodRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralFoodServiceServer).CreateFoodRecord(ctx, req.(*CreateFoodRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentralFoodService_GetFoodRecords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFoodRecordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralFoodServiceServer).GetFoodRecords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CentralFoodService_GetFoodRecords_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralFoodServiceServer).GetFoodRecords(ctx, req.(*GetFoodRecordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CentralFoodService_ServiceDesc is the grpc.ServiceDesc for CentralFoodService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CentralFoodService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "centralproto.v1.CentralFoodService",
	HandlerType: (*CentralFoodServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFoodRecord",
			Handler:    _CentralFoodService_CreateFoodRecord_Handler,
		},
		{
			MethodName: "GetFoodRecords",
			Handler:    _CentralFoodService_GetFoodRecords_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/central/central_food.proto",
}
