// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: edge.proto

package proto

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
	SensorService_SendData_FullMethodName      = "/edge.SensorService/SendData"
	SensorService_Configuration_FullMethodName = "/edge.SensorService/Configuration"
)

// SensorServiceClient is the client API for SensorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SensorServiceClient interface {
	SendData(ctx context.Context, in *SensorData, opts ...grpc.CallOption) (*Response, error)
	Configuration(ctx context.Context, in *SensorIdentification, opts ...grpc.CallOption) (*CommunicationConfiguration, error)
}

type sensorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSensorServiceClient(cc grpc.ClientConnInterface) SensorServiceClient {
	return &sensorServiceClient{cc}
}

func (c *sensorServiceClient) SendData(ctx context.Context, in *SensorData, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, SensorService_SendData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sensorServiceClient) Configuration(ctx context.Context, in *SensorIdentification, opts ...grpc.CallOption) (*CommunicationConfiguration, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommunicationConfiguration)
	err := c.cc.Invoke(ctx, SensorService_Configuration_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SensorServiceServer is the server API for SensorService service.
// All implementations must embed UnimplementedSensorServiceServer
// for forward compatibility.
type SensorServiceServer interface {
	SendData(context.Context, *SensorData) (*Response, error)
	Configuration(context.Context, *SensorIdentification) (*CommunicationConfiguration, error)
	mustEmbedUnimplementedSensorServiceServer()
}

// UnimplementedSensorServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSensorServiceServer struct{}

func (UnimplementedSensorServiceServer) SendData(context.Context, *SensorData) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendData not implemented")
}
func (UnimplementedSensorServiceServer) Configuration(context.Context, *SensorIdentification) (*CommunicationConfiguration, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configuration not implemented")
}
func (UnimplementedSensorServiceServer) mustEmbedUnimplementedSensorServiceServer() {}
func (UnimplementedSensorServiceServer) testEmbeddedByValue()                       {}

// UnsafeSensorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SensorServiceServer will
// result in compilation errors.
type UnsafeSensorServiceServer interface {
	mustEmbedUnimplementedSensorServiceServer()
}

func RegisterSensorServiceServer(s grpc.ServiceRegistrar, srv SensorServiceServer) {
	// If the following call pancis, it indicates UnimplementedSensorServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SensorService_ServiceDesc, srv)
}

func _SensorService_SendData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SensorData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SensorServiceServer).SendData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SensorService_SendData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SensorServiceServer).SendData(ctx, req.(*SensorData))
	}
	return interceptor(ctx, in, info, handler)
}

func _SensorService_Configuration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SensorIdentification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SensorServiceServer).Configuration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SensorService_Configuration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SensorServiceServer).Configuration(ctx, req.(*SensorIdentification))
	}
	return interceptor(ctx, in, info, handler)
}

// SensorService_ServiceDesc is the grpc.ServiceDesc for SensorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SensorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "edge.SensorService",
	HandlerType: (*SensorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendData",
			Handler:    _SensorService_SendData_Handler,
		},
		{
			MethodName: "Configuration",
			Handler:    _SensorService_Configuration_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "edge.proto",
}
