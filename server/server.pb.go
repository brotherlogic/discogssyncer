// Code generated by protoc-gen-go.
// source: server.proto
// DO NOT EDIT!

/*
Package server is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
	Empty
	ReleaseList
*/
package server

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import godiscogs "github.com/brotherlogic/godiscogs"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ReleaseList struct {
	Releases []*godiscogs.Release `protobuf:"bytes,1,rep,name=releases" json:"releases,omitempty"`
}

func (m *ReleaseList) Reset()                    { *m = ReleaseList{} }
func (m *ReleaseList) String() string            { return proto.CompactTextString(m) }
func (*ReleaseList) ProtoMessage()               {}
func (*ReleaseList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ReleaseList) GetReleases() []*godiscogs.Release {
	if m != nil {
		return m.Releases
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*ReleaseList)(nil), "ReleaseList")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for DiscogsService service

type DiscogsServiceClient interface {
	GetCollection(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReleaseList, error)
}

type discogsServiceClient struct {
	cc *grpc.ClientConn
}

func NewDiscogsServiceClient(cc *grpc.ClientConn) DiscogsServiceClient {
	return &discogsServiceClient{cc}
}

func (c *discogsServiceClient) GetCollection(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReleaseList, error) {
	out := new(ReleaseList)
	err := grpc.Invoke(ctx, "/DiscogsService/GetCollection", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DiscogsService service

type DiscogsServiceServer interface {
	GetCollection(context.Context, *Empty) (*ReleaseList, error)
}

func RegisterDiscogsServiceServer(s *grpc.Server, srv DiscogsServiceServer) {
	s.RegisterService(&_DiscogsService_serviceDesc, srv)
}

func _DiscogsService_GetCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).GetCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DiscogsService/GetCollection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).GetCollection(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _DiscogsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "DiscogsService",
	HandlerType: (*DiscogsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCollection",
			Handler:    _DiscogsService_GetCollection_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 169 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0x32, 0x4c, 0xcf, 0x2c, 0xc9, 0x28, 0x4d,
	0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0x02, 0x8a, 0x64, 0xa4, 0x16, 0xe5, 0xe4, 0xa7, 0x67, 0x26,
	0xeb, 0xa7, 0xe7, 0xa7, 0x64, 0x16, 0x27, 0xe7, 0xa7, 0x17, 0x23, 0x58, 0x10, 0x2d, 0x4a, 0xec,
	0x5c, 0xac, 0xae, 0xb9, 0x05, 0x25, 0x95, 0x4a, 0xc6, 0x5c, 0xdc, 0x41, 0xa9, 0x39, 0xa9, 0x89,
	0xc5, 0xa9, 0x3e, 0x99, 0xc5, 0x25, 0x42, 0x2a, 0x5c, 0x1c, 0x45, 0x10, 0x6e, 0xb1, 0x04, 0xa3,
	0x02, 0xb3, 0x06, 0xb7, 0x11, 0x87, 0x1e, 0x54, 0x3e, 0x08, 0x2e, 0x63, 0x64, 0xc9, 0xc5, 0xe7,
	0x02, 0x31, 0x2e, 0x18, 0xe8, 0x8e, 0xcc, 0xe4, 0x54, 0x21, 0x75, 0x2e, 0x5e, 0xf7, 0xd4, 0x12,
	0xe7, 0xfc, 0x9c, 0x9c, 0xd4, 0xe4, 0x92, 0xcc, 0xfc, 0x3c, 0x21, 0x36, 0x3d, 0xb0, 0xf9, 0x52,
	0x3c, 0x7a, 0x48, 0xc6, 0x2b, 0x31, 0x24, 0xb1, 0x81, 0xed, 0x37, 0x06, 0x04, 0x00, 0x00, 0xff,
	0xff, 0xb3, 0xff, 0xa9, 0x91, 0xc2, 0x00, 0x00, 0x00,
}
