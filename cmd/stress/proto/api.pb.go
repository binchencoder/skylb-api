// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cmd/stress/proto/api.proto

package proto

import (
	context "context"
	fmt "fmt"
	_ "github.com/binchencoder/ease-gateway/httpoptions"
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

type Gender int32

const (
	Gender_MALE   Gender = 0
	Gender_FEMALE Gender = 1
)

var Gender_name = map[int32]string{
	0: "MALE",
	1: "FEMALE",
}

var Gender_value = map[string]int32{
	"MALE":   0,
	"FEMALE": 1,
}

func (x Gender) String() string {
	return proto.EnumName(Gender_name, int32(x))
}

func (Gender) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_17fe16740fdb1c07, []int{0}
}

type SayHelloReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Gender               Gender   `protobuf:"varint,2,opt,name=gender,proto3,enum=proto.Gender" json:"gender,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SayHelloReq) Reset()         { *m = SayHelloReq{} }
func (m *SayHelloReq) String() string { return proto.CompactTextString(m) }
func (*SayHelloReq) ProtoMessage()    {}
func (*SayHelloReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_17fe16740fdb1c07, []int{0}
}

func (m *SayHelloReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayHelloReq.Unmarshal(m, b)
}
func (m *SayHelloReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayHelloReq.Marshal(b, m, deterministic)
}
func (m *SayHelloReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayHelloReq.Merge(m, src)
}
func (m *SayHelloReq) XXX_Size() int {
	return xxx_messageInfo_SayHelloReq.Size(m)
}
func (m *SayHelloReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SayHelloReq.DiscardUnknown(m)
}

var xxx_messageInfo_SayHelloReq proto.InternalMessageInfo

func (m *SayHelloReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SayHelloReq) GetGender() Gender {
	if m != nil {
		return m.Gender
	}
	return Gender_MALE
}

type SayHelloResp struct {
	Greeting             string   `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"`
	Peer                 string   `protobuf:"bytes,2,opt,name=peer,proto3" json:"peer,omitempty"`
	ServiceId            int32    `protobuf:"varint,3,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	Tid                  string   `protobuf:"bytes,4,opt,name=tid,proto3" json:"tid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SayHelloResp) Reset()         { *m = SayHelloResp{} }
func (m *SayHelloResp) String() string { return proto.CompactTextString(m) }
func (*SayHelloResp) ProtoMessage()    {}
func (*SayHelloResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_17fe16740fdb1c07, []int{1}
}

func (m *SayHelloResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayHelloResp.Unmarshal(m, b)
}
func (m *SayHelloResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayHelloResp.Marshal(b, m, deterministic)
}
func (m *SayHelloResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayHelloResp.Merge(m, src)
}
func (m *SayHelloResp) XXX_Size() int {
	return xxx_messageInfo_SayHelloResp.Size(m)
}
func (m *SayHelloResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SayHelloResp.DiscardUnknown(m)
}

var xxx_messageInfo_SayHelloResp proto.InternalMessageInfo

func (m *SayHelloResp) GetGreeting() string {
	if m != nil {
		return m.Greeting
	}
	return ""
}

func (m *SayHelloResp) GetPeer() string {
	if m != nil {
		return m.Peer
	}
	return ""
}

func (m *SayHelloResp) GetServiceId() int32 {
	if m != nil {
		return m.ServiceId
	}
	return 0
}

func (m *SayHelloResp) GetTid() string {
	if m != nil {
		return m.Tid
	}
	return ""
}

func init() {
	proto.RegisterEnum("proto.Gender", Gender_name, Gender_value)
	proto.RegisterType((*SayHelloReq)(nil), "proto.SayHelloReq")
	proto.RegisterType((*SayHelloResp)(nil), "proto.SayHelloResp")
}

func init() { proto.RegisterFile("cmd/stress/proto/api.proto", fileDescriptor_17fe16740fdb1c07) }

var fileDescriptor_17fe16740fdb1c07 = []byte{
	// 345 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xd1, 0x4a, 0x32, 0x41,
	0x14, 0xc7, 0xbf, 0xfd, 0xd4, 0xfd, 0xd6, 0xf3, 0x69, 0xc8, 0xe9, 0x66, 0x59, 0x30, 0x16, 0x21,
	0x90, 0x42, 0x07, 0x6d, 0xe9, 0xa2, 0xbb, 0x20, 0xcb, 0xa0, 0x6e, 0xd6, 0x07, 0x88, 0xd5, 0x99,
	0x74, 0x40, 0x77, 0xa6, 0x9d, 0xc1, 0x10, 0xf1, 0x31, 0x7a, 0x90, 0x5e, 0xc3, 0x47, 0xe9, 0xd6,
	0xfb, 0x8c, 0x9d, 0x5d, 0x2b, 0xa8, 0xa0, 0xab, 0xf9, 0xef, 0x39, 0x7b, 0x7e, 0xff, 0x73, 0xfe,
	0xe0, 0x8d, 0x66, 0x94, 0x28, 0x9d, 0x30, 0xa5, 0x88, 0x4c, 0x84, 0x16, 0x24, 0x92, 0xbc, 0x6d,
	0x14, 0x96, 0xcc, 0xe3, 0xd5, 0x27, 0x5a, 0x4b, 0x21, 0x35, 0x17, 0xb1, 0x22, 0x51, 0x1c, 0x0b,
	0x1d, 0x19, 0x9d, 0xfd, 0xd5, 0xe8, 0xc3, 0xff, 0x41, 0xb4, 0xe8, 0xb3, 0xe9, 0x54, 0x84, 0xec,
	0x01, 0x11, 0x8a, 0x71, 0x34, 0x63, 0xae, 0xe5, 0x5b, 0xcd, 0x72, 0x68, 0x34, 0x1e, 0x82, 0x3d,
	0x66, 0x31, 0x65, 0x89, 0xfb, 0xd7, 0xb7, 0x9a, 0x7b, 0xdd, 0x6a, 0x36, 0xda, 0xbe, 0x32, 0xc5,
	0x30, 0x6f, 0x36, 0x04, 0x54, 0x3e, 0x48, 0x4a, 0xa2, 0x07, 0xce, 0x38, 0x61, 0x4c, 0xf3, 0x78,
	0x9c, 0xe3, 0xde, 0xbf, 0x53, 0x1b, 0xc9, 0x72, 0x60, 0x39, 0x34, 0x1a, 0xeb, 0x00, 0x8a, 0x25,
	0x73, 0x3e, 0x62, 0x77, 0x9c, 0xba, 0x05, 0xdf, 0x6a, 0x96, 0xc2, 0x72, 0x5e, 0xb9, 0xa6, 0x58,
	0x83, 0x82, 0xe6, 0xd4, 0x2d, 0x9a, 0x89, 0x54, 0x1e, 0x1d, 0x80, 0x9d, 0xad, 0x80, 0x0e, 0x14,
	0x6f, 0xcf, 0x6f, 0x7a, 0xb5, 0x3f, 0x08, 0x60, 0x5f, 0xf6, 0x8c, 0xb6, 0xba, 0xaf, 0x16, 0x54,
	0x07, 0x26, 0x9b, 0x41, 0x46, 0xc1, 0x7b, 0x70, 0x76, 0x2b, 0x22, 0xe6, 0x57, 0x7c, 0xba, 0xde,
	0xdb, 0xff, 0x52, 0x53, 0xb2, 0xd1, 0x59, 0x6f, 0x82, 0x16, 0x1e, 0xef, 0x62, 0x9e, 0x77, 0x88,
	0x8a, 0x16, 0xad, 0x49, 0xda, 0x27, 0xcb, 0x34, 0xa0, 0x15, 0xc9, 0x12, 0x20, 0xcb, 0xec, 0x5d,
	0xe1, 0x23, 0xd4, 0x76, 0x88, 0x0b, 0xae, 0xa2, 0xe1, 0x94, 0xd1, 0xdf, 0xfb, 0x9d, 0xad, 0x37,
	0xc1, 0x29, 0x06, 0xdf, 0xf9, 0xb5, 0x68, 0x4e, 0xfb, 0xc1, 0xd8, 0xab, 0xbc, 0x6c, 0x02, 0xc7,
	0x79, 0x7e, 0xda, 0x6e, 0xff, 0xf9, 0xd6, 0xd0, 0x36, 0xf4, 0x93, 0xb7, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x29, 0x2f, 0x88, 0x4e, 0x26, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StressServiceClient is the client API for StressService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StressServiceClient interface {
	SayHello(ctx context.Context, in *SayHelloReq, opts ...grpc.CallOption) (*SayHelloResp, error)
	SayHelloDisabled(ctx context.Context, in *SayHelloReq, opts ...grpc.CallOption) (*SayHelloResp, error)
}

type stressServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStressServiceClient(cc grpc.ClientConnInterface) StressServiceClient {
	return &stressServiceClient{cc}
}

func (c *stressServiceClient) SayHello(ctx context.Context, in *SayHelloReq, opts ...grpc.CallOption) (*SayHelloResp, error) {
	out := new(SayHelloResp)
	err := c.cc.Invoke(ctx, "/proto.StressService/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stressServiceClient) SayHelloDisabled(ctx context.Context, in *SayHelloReq, opts ...grpc.CallOption) (*SayHelloResp, error) {
	out := new(SayHelloResp)
	err := c.cc.Invoke(ctx, "/proto.StressService/SayHelloDisabled", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StressServiceServer is the server API for StressService service.
type StressServiceServer interface {
	SayHello(context.Context, *SayHelloReq) (*SayHelloResp, error)
	SayHelloDisabled(context.Context, *SayHelloReq) (*SayHelloResp, error)
}

// UnimplementedStressServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStressServiceServer struct {
}

func (*UnimplementedStressServiceServer) SayHello(ctx context.Context, req *SayHelloReq) (*SayHelloResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (*UnimplementedStressServiceServer) SayHelloDisabled(ctx context.Context, req *SayHelloReq) (*SayHelloResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHelloDisabled not implemented")
}

func RegisterStressServiceServer(s *grpc.Server, srv StressServiceServer) {
	s.RegisterService(&_StressService_serviceDesc, srv)
}

func _StressService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayHelloReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StressServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.StressService/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StressServiceServer).SayHello(ctx, req.(*SayHelloReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _StressService_SayHelloDisabled_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayHelloReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StressServiceServer).SayHelloDisabled(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.StressService/SayHelloDisabled",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StressServiceServer).SayHelloDisabled(ctx, req.(*SayHelloReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _StressService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.StressService",
	HandlerType: (*StressServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _StressService_SayHello_Handler,
		},
		{
			MethodName: "SayHelloDisabled",
			Handler:    _StressService_SayHelloDisabled_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cmd/stress/proto/api.proto",
}
