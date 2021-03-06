// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cmd/skytest/proto/api.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GreetingRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Age                  int32    `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GreetingRequest) Reset()         { *m = GreetingRequest{} }
func (m *GreetingRequest) String() string { return proto.CompactTextString(m) }
func (*GreetingRequest) ProtoMessage()    {}
func (*GreetingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6c2e35c6833a0d4, []int{0}
}

func (m *GreetingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GreetingRequest.Unmarshal(m, b)
}
func (m *GreetingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GreetingRequest.Marshal(b, m, deterministic)
}
func (m *GreetingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GreetingRequest.Merge(m, src)
}
func (m *GreetingRequest) XXX_Size() int {
	return xxx_messageInfo_GreetingRequest.Size(m)
}
func (m *GreetingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GreetingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GreetingRequest proto.InternalMessageInfo

func (m *GreetingRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GreetingRequest) GetAge() int32 {
	if m != nil {
		return m.Age
	}
	return 0
}

type GreetingResponse struct {
	Greeting             string   `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GreetingResponse) Reset()         { *m = GreetingResponse{} }
func (m *GreetingResponse) String() string { return proto.CompactTextString(m) }
func (*GreetingResponse) ProtoMessage()    {}
func (*GreetingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6c2e35c6833a0d4, []int{1}
}

func (m *GreetingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GreetingResponse.Unmarshal(m, b)
}
func (m *GreetingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GreetingResponse.Marshal(b, m, deterministic)
}
func (m *GreetingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GreetingResponse.Merge(m, src)
}
func (m *GreetingResponse) XXX_Size() int {
	return xxx_messageInfo_GreetingResponse.Size(m)
}
func (m *GreetingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GreetingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GreetingResponse proto.InternalMessageInfo

func (m *GreetingResponse) GetGreeting() string {
	if m != nil {
		return m.Greeting
	}
	return ""
}

func init() {
	proto.RegisterType((*GreetingRequest)(nil), "proto.GreetingRequest")
	proto.RegisterType((*GreetingResponse)(nil), "proto.GreetingResponse")
}

func init() { proto.RegisterFile("cmd/skytest/proto/api.proto", fileDescriptor_a6c2e35c6833a0d4) }

var fileDescriptor_a6c2e35c6833a0d4 = []byte{
	// 201 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4e, 0xce, 0x4d, 0xd1,
	0x2f, 0xce, 0xae, 0x2c, 0x49, 0x2d, 0x2e, 0xd1, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0x2c,
	0xc8, 0xd4, 0x03, 0xb3, 0x84, 0x58, 0xc1, 0x94, 0x92, 0x39, 0x17, 0xbf, 0x7b, 0x51, 0x6a, 0x6a,
	0x49, 0x66, 0x5e, 0x7a, 0x50, 0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89, 0x90, 0x10, 0x17, 0x4b, 0x5e,
	0x62, 0x6e, 0xaa, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x98, 0x2d, 0x24, 0xc0, 0xc5, 0x9c,
	0x98, 0x9e, 0x2a, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x1a, 0x04, 0x62, 0x2a, 0xe9, 0x71, 0x09, 0x20,
	0x34, 0x16, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x49, 0x71, 0x71, 0xa4, 0x43, 0xc5, 0xa0, 0xba,
	0xe1, 0x7c, 0x23, 0x37, 0x2e, 0xf6, 0x60, 0x88, 0x53, 0x84, 0xac, 0xb9, 0x38, 0x60, 0x5a, 0x85,
	0xc4, 0x20, 0xce, 0xd1, 0x43, 0x73, 0x84, 0x94, 0x38, 0x86, 0x38, 0xc4, 0x0e, 0x27, 0x73, 0x2e,
	0x95, 0xe4, 0xfc, 0x5c, 0xbd, 0xa4, 0xcc, 0xbc, 0xe4, 0x8c, 0xd4, 0xbc, 0xe4, 0xfc, 0x94, 0xd4,
	0x22, 0xbd, 0xe2, 0xec, 0xca, 0x9c, 0x24, 0x3d, 0xa8, 0x4f, 0x21, 0x1a, 0x9d, 0x78, 0xa1, 0xb6,
	0x05, 0x80, 0x78, 0xc5, 0x01, 0x0c, 0x49, 0x6c, 0x60, 0x71, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x62, 0x5b, 0xa3, 0x94, 0x16, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SkytestClient is the client API for Skytest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SkytestClient interface {
	Greeting(ctx context.Context, in *GreetingRequest, opts ...grpc.CallOption) (*GreetingResponse, error)
}

type skytestClient struct {
	cc grpc.ClientConnInterface
}

func NewSkytestClient(cc grpc.ClientConnInterface) SkytestClient {
	return &skytestClient{cc}
}

func (c *skytestClient) Greeting(ctx context.Context, in *GreetingRequest, opts ...grpc.CallOption) (*GreetingResponse, error) {
	out := new(GreetingResponse)
	err := c.cc.Invoke(ctx, "/proto.Skytest/Greeting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SkytestServer is the server API for Skytest service.
type SkytestServer interface {
	Greeting(context.Context, *GreetingRequest) (*GreetingResponse, error)
}

// UnimplementedSkytestServer can be embedded to have forward compatible implementations.
type UnimplementedSkytestServer struct {
}

func (*UnimplementedSkytestServer) Greeting(ctx context.Context, req *GreetingRequest) (*GreetingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Greeting not implemented")
}

func RegisterSkytestServer(s *grpc.Server, srv SkytestServer) {
	s.RegisterService(&_Skytest_serviceDesc, srv)
}

func _Skytest_Greeting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GreetingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SkytestServer).Greeting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Skytest/Greeting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SkytestServer).Greeting(ctx, req.(*GreetingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Skytest_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Skytest",
	HandlerType: (*SkytestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Greeting",
			Handler:    _Skytest_Greeting_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cmd/skytest/proto/api.proto",
}
