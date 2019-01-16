// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package rpctransport

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type EmptyRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyRequest) Reset()         { *m = EmptyRequest{} }
func (m *EmptyRequest) String() string { return proto.CompactTextString(m) }
func (*EmptyRequest) ProtoMessage()    {}
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *EmptyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyRequest.Unmarshal(m, b)
}
func (m *EmptyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyRequest.Marshal(b, m, deterministic)
}
func (m *EmptyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyRequest.Merge(m, src)
}
func (m *EmptyRequest) XXX_Size() int {
	return xxx_messageInfo_EmptyRequest.Size(m)
}
func (m *EmptyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyRequest proto.InternalMessageInfo

type StatusReply struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusReply) Reset()         { *m = StatusReply{} }
func (m *StatusReply) String() string { return proto.CompactTextString(m) }
func (*StatusReply) ProtoMessage()    {}
func (*StatusReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *StatusReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusReply.Unmarshal(m, b)
}
func (m *StatusReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusReply.Marshal(b, m, deterministic)
}
func (m *StatusReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusReply.Merge(m, src)
}
func (m *StatusReply) XXX_Size() int {
	return xxx_messageInfo_StatusReply.Size(m)
}
func (m *StatusReply) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusReply.DiscardUnknown(m)
}

var xxx_messageInfo_StatusReply proto.InternalMessageInfo

func (m *StatusReply) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type GetCategoryRequest struct {
	CategoryID           string   `protobuf:"bytes,1,opt,name=categoryID,proto3" json:"categoryID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCategoryRequest) Reset()         { *m = GetCategoryRequest{} }
func (m *GetCategoryRequest) String() string { return proto.CompactTextString(m) }
func (*GetCategoryRequest) ProtoMessage()    {}
func (*GetCategoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{2}
}

func (m *GetCategoryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCategoryRequest.Unmarshal(m, b)
}
func (m *GetCategoryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCategoryRequest.Marshal(b, m, deterministic)
}
func (m *GetCategoryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCategoryRequest.Merge(m, src)
}
func (m *GetCategoryRequest) XXX_Size() int {
	return xxx_messageInfo_GetCategoryRequest.Size(m)
}
func (m *GetCategoryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCategoryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetCategoryRequest proto.InternalMessageInfo

func (m *GetCategoryRequest) GetCategoryID() string {
	if m != nil {
		return m.CategoryID
	}
	return ""
}

type GetCategoryReply struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ParentID             string   `protobuf:"bytes,3,opt,name=parentID,proto3" json:"parentID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCategoryReply) Reset()         { *m = GetCategoryReply{} }
func (m *GetCategoryReply) String() string { return proto.CompactTextString(m) }
func (*GetCategoryReply) ProtoMessage()    {}
func (*GetCategoryReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{3}
}

func (m *GetCategoryReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCategoryReply.Unmarshal(m, b)
}
func (m *GetCategoryReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCategoryReply.Marshal(b, m, deterministic)
}
func (m *GetCategoryReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCategoryReply.Merge(m, src)
}
func (m *GetCategoryReply) XXX_Size() int {
	return xxx_messageInfo_GetCategoryReply.Size(m)
}
func (m *GetCategoryReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCategoryReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetCategoryReply proto.InternalMessageInfo

func (m *GetCategoryReply) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *GetCategoryReply) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetCategoryReply) GetParentID() string {
	if m != nil {
		return m.ParentID
	}
	return ""
}

func init() {
	proto.RegisterType((*EmptyRequest)(nil), "rpctransport.EmptyRequest")
	proto.RegisterType((*StatusReply)(nil), "rpctransport.StatusReply")
	proto.RegisterType((*GetCategoryRequest)(nil), "rpctransport.GetCategoryRequest")
	proto.RegisterType((*GetCategoryReply)(nil), "rpctransport.GetCategoryReply")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x4d, 0x94, 0x62, 0xa7, 0xb5, 0xca, 0x1c, 0x24, 0xe6, 0x50, 0xca, 0x82, 0xe0, 0x29,
	0x07, 0xf5, 0x0d, 0x8c, 0x84, 0xdc, 0x34, 0xfa, 0x02, 0xeb, 0x76, 0x28, 0x01, 0x9b, 0x5d, 0x67,
	0xa7, 0x42, 0x5f, 0xc9, 0xa7, 0x14, 0xd7, 0x18, 0x36, 0x08, 0xbd, 0xed, 0xff, 0xef, 0xff, 0xef,
	0xce, 0xc7, 0xc0, 0x99, 0x27, 0xfe, 0x6c, 0x0d, 0x15, 0x8e, 0xad, 0x58, 0x9c, 0xb3, 0x33, 0xc2,
	0xba, 0xf3, 0xce, 0xb2, 0xa8, 0x05, 0xcc, 0x1f, 0xb7, 0x4e, 0xf6, 0x0d, 0x7d, 0xec, 0xc8, 0x8b,
	0xba, 0x86, 0xd9, 0x8b, 0x68, 0xd9, 0xf9, 0x86, 0xdc, 0xfb, 0x1e, 0x2f, 0x61, 0xe2, 0x83, 0xcc,
	0x92, 0x55, 0x72, 0x33, 0x6d, 0x7a, 0xa5, 0xee, 0x01, 0x2b, 0x92, 0x07, 0x2d, 0xb4, 0xb1, 0xfc,
	0x57, 0xc6, 0x25, 0x80, 0xe9, 0xad, 0xba, 0xec, 0x1b, 0x91, 0xa3, 0x1a, 0xb8, 0x18, 0xb5, 0x7e,
	0x7e, 0x58, 0x40, 0x3a, 0x64, 0xd3, 0xba, 0x44, 0x84, 0x93, 0x4e, 0x6f, 0x29, 0x4b, 0x83, 0x13,
	0xce, 0x98, 0xc3, 0xa9, 0xd3, 0x4c, 0x9d, 0xd4, 0x65, 0x76, 0x1c, 0xfc, 0x41, 0xdf, 0x7e, 0x25,
	0x70, 0xfe, 0xc4, 0xda, 0x48, 0x6b, 0xda, 0x6e, 0x53, 0xd9, 0xd7, 0xf5, 0x1a, 0x4b, 0x98, 0x56,
	0x24, 0xbf, 0x1c, 0x98, 0x17, 0x31, 0x70, 0x11, 0xd3, 0xe6, 0x57, 0xe3, 0xbb, 0x88, 0x5c, 0x1d,
	0xe1, 0x33, 0xcc, 0xa2, 0x69, 0x71, 0x35, 0xce, 0xfe, 0xc7, 0xcf, 0x97, 0x07, 0x12, 0xe1, 0xc9,
	0xb7, 0x49, 0x58, 0xc1, 0xdd, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf2, 0x42, 0x0e, 0x3f, 0x93,
	0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PracticingGoTddClient is the client API for PracticingGoTdd service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PracticingGoTddClient interface {
	GetStatus(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*StatusReply, error)
	GetCategory(ctx context.Context, in *GetCategoryRequest, opts ...grpc.CallOption) (*GetCategoryReply, error)
}

type practicingGoTddClient struct {
	cc *grpc.ClientConn
}

func NewPracticingGoTddClient(cc *grpc.ClientConn) PracticingGoTddClient {
	return &practicingGoTddClient{cc}
}

func (c *practicingGoTddClient) GetStatus(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := c.cc.Invoke(ctx, "/rpctransport.PracticingGoTdd/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *practicingGoTddClient) GetCategory(ctx context.Context, in *GetCategoryRequest, opts ...grpc.CallOption) (*GetCategoryReply, error) {
	out := new(GetCategoryReply)
	err := c.cc.Invoke(ctx, "/rpctransport.PracticingGoTdd/GetCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PracticingGoTddServer is the server API for PracticingGoTdd service.
type PracticingGoTddServer interface {
	GetStatus(context.Context, *EmptyRequest) (*StatusReply, error)
	GetCategory(context.Context, *GetCategoryRequest) (*GetCategoryReply, error)
}

func RegisterPracticingGoTddServer(s *grpc.Server, srv PracticingGoTddServer) {
	s.RegisterService(&_PracticingGoTdd_serviceDesc, srv)
}

func _PracticingGoTdd_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PracticingGoTddServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpctransport.PracticingGoTdd/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PracticingGoTddServer).GetStatus(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PracticingGoTdd_GetCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PracticingGoTddServer).GetCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpctransport.PracticingGoTdd/GetCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PracticingGoTddServer).GetCategory(ctx, req.(*GetCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PracticingGoTdd_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpctransport.PracticingGoTdd",
	HandlerType: (*PracticingGoTddServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _PracticingGoTdd_GetStatus_Handler,
		},
		{
			MethodName: "GetCategory",
			Handler:    _PracticingGoTdd_GetCategory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
