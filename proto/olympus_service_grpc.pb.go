// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.2
// source: olympus_service.proto

package proto

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

// OlympusClient is the client API for Olympus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OlympusClient interface {
	Zone(ctx context.Context, opts ...grpc.CallOption) (Olympus_ZoneClient, error)
	Tracking(ctx context.Context, opts ...grpc.CallOption) (Olympus_TrackingClient, error)
}

type olympusClient struct {
	cc grpc.ClientConnInterface
}

func NewOlympusClient(cc grpc.ClientConnInterface) OlympusClient {
	return &olympusClient{cc}
}

func (c *olympusClient) Zone(ctx context.Context, opts ...grpc.CallOption) (Olympus_ZoneClient, error) {
	stream, err := c.cc.NewStream(ctx, &Olympus_ServiceDesc.Streams[0], "/proto.Olympus/Zone", opts...)
	if err != nil {
		return nil, err
	}
	x := &olympusZoneClient{stream}
	return x, nil
}

type Olympus_ZoneClient interface {
	Send(*ZoneUpStream) error
	Recv() (*ZoneDownStream, error)
	grpc.ClientStream
}

type olympusZoneClient struct {
	grpc.ClientStream
}

func (x *olympusZoneClient) Send(m *ZoneUpStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *olympusZoneClient) Recv() (*ZoneDownStream, error) {
	m := new(ZoneDownStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *olympusClient) Tracking(ctx context.Context, opts ...grpc.CallOption) (Olympus_TrackingClient, error) {
	stream, err := c.cc.NewStream(ctx, &Olympus_ServiceDesc.Streams[1], "/proto.Olympus/Tracking", opts...)
	if err != nil {
		return nil, err
	}
	x := &olympusTrackingClient{stream}
	return x, nil
}

type Olympus_TrackingClient interface {
	Send(*TrackingUpStream) error
	Recv() (*TrackingDownStream, error)
	grpc.ClientStream
}

type olympusTrackingClient struct {
	grpc.ClientStream
}

func (x *olympusTrackingClient) Send(m *TrackingUpStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *olympusTrackingClient) Recv() (*TrackingDownStream, error) {
	m := new(TrackingDownStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OlympusServer is the server API for Olympus service.
// All implementations must embed UnimplementedOlympusServer
// for forward compatibility
type OlympusServer interface {
	Zone(Olympus_ZoneServer) error
	Tracking(Olympus_TrackingServer) error
	mustEmbedUnimplementedOlympusServer()
}

// UnimplementedOlympusServer must be embedded to have forward compatible implementations.
type UnimplementedOlympusServer struct {
}

func (UnimplementedOlympusServer) Zone(Olympus_ZoneServer) error {
	return status.Errorf(codes.Unimplemented, "method Zone not implemented")
}
func (UnimplementedOlympusServer) Tracking(Olympus_TrackingServer) error {
	return status.Errorf(codes.Unimplemented, "method Tracking not implemented")
}
func (UnimplementedOlympusServer) mustEmbedUnimplementedOlympusServer() {}

// UnsafeOlympusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OlympusServer will
// result in compilation errors.
type UnsafeOlympusServer interface {
	mustEmbedUnimplementedOlympusServer()
}

func RegisterOlympusServer(s grpc.ServiceRegistrar, srv OlympusServer) {
	s.RegisterService(&Olympus_ServiceDesc, srv)
}

func _Olympus_Zone_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OlympusServer).Zone(&olympusZoneServer{stream})
}

type Olympus_ZoneServer interface {
	Send(*ZoneDownStream) error
	Recv() (*ZoneUpStream, error)
	grpc.ServerStream
}

type olympusZoneServer struct {
	grpc.ServerStream
}

func (x *olympusZoneServer) Send(m *ZoneDownStream) error {
	return x.ServerStream.SendMsg(m)
}

func (x *olympusZoneServer) Recv() (*ZoneUpStream, error) {
	m := new(ZoneUpStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Olympus_Tracking_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OlympusServer).Tracking(&olympusTrackingServer{stream})
}

type Olympus_TrackingServer interface {
	Send(*TrackingDownStream) error
	Recv() (*TrackingUpStream, error)
	grpc.ServerStream
}

type olympusTrackingServer struct {
	grpc.ServerStream
}

func (x *olympusTrackingServer) Send(m *TrackingDownStream) error {
	return x.ServerStream.SendMsg(m)
}

func (x *olympusTrackingServer) Recv() (*TrackingUpStream, error) {
	m := new(TrackingUpStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Olympus_ServiceDesc is the grpc.ServiceDesc for Olympus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Olympus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Olympus",
	HandlerType: (*OlympusServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Zone",
			Handler:       _Olympus_Zone_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Tracking",
			Handler:       _Olympus_Tracking_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "olympus_service.proto",
}