// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// RawDataClient is the client API for RawData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RawDataClient interface {
	Transform(ctx context.Context, in *TransformRequest, opts ...grpc.CallOption) (*TransformResponse, error)
	TrackPersons(ctx context.Context, in *TrackRequest, opts ...grpc.CallOption) (*TrackResponse, error)
}

type rawDataClient struct {
	cc grpc.ClientConnInterface
}

func NewRawDataClient(cc grpc.ClientConnInterface) RawDataClient {
	return &rawDataClient{cc}
}

func (c *rawDataClient) Transform(ctx context.Context, in *TransformRequest, opts ...grpc.CallOption) (*TransformResponse, error) {
	out := new(TransformResponse)
	err := c.cc.Invoke(ctx, "/pb.RawData/Transform", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rawDataClient) TrackPersons(ctx context.Context, in *TrackRequest, opts ...grpc.CallOption) (*TrackResponse, error) {
	out := new(TrackResponse)
	err := c.cc.Invoke(ctx, "/pb.RawData/TrackPersons", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RawDataServer is the server API for RawData service.
// All implementations must embed UnimplementedRawDataServer
// for forward compatibility
type RawDataServer interface {
	Transform(context.Context, *TransformRequest) (*TransformResponse, error)
	TrackPersons(context.Context, *TrackRequest) (*TrackResponse, error)
	mustEmbedUnimplementedRawDataServer()
}

// UnimplementedRawDataServer must be embedded to have forward compatible implementations.
type UnimplementedRawDataServer struct {
}

func (UnimplementedRawDataServer) Transform(context.Context, *TransformRequest) (*TransformResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transform not implemented")
}
func (UnimplementedRawDataServer) TrackPersons(context.Context, *TrackRequest) (*TrackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TrackPersons not implemented")
}
func (UnimplementedRawDataServer) mustEmbedUnimplementedRawDataServer() {}

// UnsafeRawDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RawDataServer will
// result in compilation errors.
type UnsafeRawDataServer interface {
	mustEmbedUnimplementedRawDataServer()
}

func RegisterRawDataServer(s grpc.ServiceRegistrar, srv RawDataServer) {
	s.RegisterService(&RawData_ServiceDesc, srv)
}

func _RawData_Transform_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransformRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RawDataServer).Transform(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RawData/Transform",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RawDataServer).Transform(ctx, req.(*TransformRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RawData_TrackPersons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RawDataServer).TrackPersons(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RawData/TrackPersons",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RawDataServer).TrackPersons(ctx, req.(*TrackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RawData_ServiceDesc is the grpc.ServiceDesc for RawData service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RawData_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.RawData",
	HandlerType: (*RawDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Transform",
			Handler:    _RawData_Transform_Handler,
		},
		{
			MethodName: "TrackPersons",
			Handler:    _RawData_TrackPersons_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rawdata.proto",
}
