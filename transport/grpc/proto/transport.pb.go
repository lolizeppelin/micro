// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transport/grpc/proto/transport.proto

package go_micro_transport_grpc

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

type Message struct {
	Header               map[string]string `protobuf:"bytes,1,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 []byte            `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_651718cd7c7ae974, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetHeader() map[string]string {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Message) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "go.micro.transport.grpc.Message")
	proto.RegisterMapType((map[string]string)(nil), "go.micro.transport.grpc.Message.HeaderEntry")
}

func init() {
	proto.RegisterFile("transport/grpc/proto/transport.proto", fileDescriptor_651718cd7c7ae974)
}

var fileDescriptor_651718cd7c7ae974 = []byte{
	// 209 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x29, 0x29, 0x4a, 0xcc,
	0x2b, 0x2e, 0xc8, 0x2f, 0x2a, 0xd1, 0x4f, 0x2f, 0x2a, 0x48, 0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9,
	0xd7, 0x87, 0x0b, 0xea, 0x81, 0xf9, 0x42, 0xe2, 0xe9, 0xf9, 0x7a, 0xb9, 0x99, 0xc9, 0x45, 0xf9,
	0x7a, 0x08, 0x19, 0x90, 0x72, 0xa5, 0x79, 0x8c, 0x5c, 0xec, 0xbe, 0xa9, 0xc5, 0xc5, 0x89, 0xe9,
	0xa9, 0x42, 0x2e, 0x5c, 0x6c, 0x19, 0xa9, 0x89, 0x29, 0xa9, 0x45, 0x12, 0x8c, 0x0a, 0xcc, 0x1a,
	0xdc, 0x46, 0x3a, 0x7a, 0x38, 0x74, 0xe9, 0x41, 0x75, 0xe8, 0x79, 0x80, 0x95, 0xbb, 0xe6, 0x95,
	0x14, 0x55, 0x06, 0x41, 0xf5, 0x0a, 0x09, 0x71, 0xb1, 0x24, 0xe5, 0xa7, 0x54, 0x4a, 0x30, 0x29,
	0x30, 0x6a, 0xf0, 0x04, 0x81, 0xd9, 0x52, 0x96, 0x5c, 0xdc, 0x48, 0x4a, 0x85, 0x04, 0xb8, 0x98,
	0xb3, 0x53, 0x2b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x40, 0x4c, 0x21, 0x11, 0x2e, 0xd6,
	0xb2, 0xc4, 0x9c, 0xd2, 0x54, 0xb0, 0x2e, 0xce, 0x20, 0x08, 0xc7, 0x8a, 0xc9, 0x82, 0xd1, 0x28,
	0x9e, 0x8b, 0x33, 0x04, 0x66, 0xb9, 0x50, 0x10, 0x17, 0x5b, 0x70, 0x49, 0x51, 0x6a, 0x62, 0xae,
	0x90, 0x02, 0x21, 0xb7, 0x49, 0x11, 0x54, 0xa1, 0xc4, 0xa0, 0xc1, 0x68, 0xc0, 0x98, 0xc4, 0x06,
	0x0e, 0x21, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd4, 0xd0, 0x4b, 0x4b, 0x49, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TransportClient is the client API for Transport service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TransportClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (Transport_StreamClient, error)
}

type transportClient struct {
	cc *grpc.ClientConn
}

func NewTransportClient(cc *grpc.ClientConn) TransportClient {
	return &transportClient{cc}
}

func (c *transportClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Transport_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Transport_serviceDesc.Streams[0], "/go.micro.transport.grpc.Transport/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &transportStreamClient{stream}
	return x, nil
}

type Transport_StreamClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type transportStreamClient struct {
	grpc.ClientStream
}

func (x *transportStreamClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *transportStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TransportServer is the server API for Transport service.
type TransportServer interface {
	Stream(Transport_StreamServer) error
}

// UnimplementedTransportServer can be embedded to have forward compatible implementations.
type UnimplementedTransportServer struct {
}

func (*UnimplementedTransportServer) Stream(srv Transport_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}

func RegisterTransportServer(s *grpc.Server, srv TransportServer) {
	s.RegisterService(&_Transport_serviceDesc, srv)
}

func _Transport_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TransportServer).Stream(&transportStreamServer{stream})
}

type Transport_StreamServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type transportStreamServer struct {
	grpc.ServerStream
}

func (x *transportStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *transportStreamServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Transport_serviceDesc = grpc.ServiceDesc{
	ServiceName: "go.micro.transport.grpc.Transport",
	HandlerType: (*TransportServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Transport_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "transport/grpc/proto/transport.proto",
}