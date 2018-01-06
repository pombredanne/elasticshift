// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/vcs.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type VCSAccount struct {
	ID        string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Type      string `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
	OwnerType string `protobuf:"bytes,4,opt,name=owner_type,json=ownerType" json:"owner_type,omitempty"`
	AvatarUrl string `protobuf:"bytes,5,opt,name=avatar_url,json=avatarUrl" json:"avatar_url,omitempty"`
}

func (m *VCSAccount) Reset()                    { *m = VCSAccount{} }
func (m *VCSAccount) String() string            { return proto.CompactTextString(m) }
func (*VCSAccount) ProtoMessage()               {}
func (*VCSAccount) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *VCSAccount) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *VCSAccount) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *VCSAccount) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *VCSAccount) GetOwnerType() string {
	if m != nil {
		return m.OwnerType
	}
	return ""
}

func (m *VCSAccount) GetAvatarUrl() string {
	if m != nil {
		return m.AvatarUrl
	}
	return ""
}

type GetVCSReq struct {
}

func (m *GetVCSReq) Reset()                    { *m = GetVCSReq{} }
func (m *GetVCSReq) String() string            { return proto.CompactTextString(m) }
func (*GetVCSReq) ProtoMessage()               {}
func (*GetVCSReq) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{1} }

type GetVCSRes struct {
	Accounts []*VCSAccount `protobuf:"bytes,1,rep,name=accounts" json:"accounts,omitempty"`
}

func (m *GetVCSRes) Reset()                    { *m = GetVCSRes{} }
func (m *GetVCSRes) String() string            { return proto.CompactTextString(m) }
func (*GetVCSRes) ProtoMessage()               {}
func (*GetVCSRes) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{2} }

func (m *GetVCSRes) GetAccounts() []*VCSAccount {
	if m != nil {
		return m.Accounts
	}
	return nil
}

type SyncVCSReq struct {
	VcsId string `protobuf:"bytes,1,opt,name=vcs_id,json=vcsId" json:"vcs_id,omitempty"`
}

func (m *SyncVCSReq) Reset()                    { *m = SyncVCSReq{} }
func (m *SyncVCSReq) String() string            { return proto.CompactTextString(m) }
func (*SyncVCSReq) ProtoMessage()               {}
func (*SyncVCSReq) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{3} }

func (m *SyncVCSReq) GetVcsId() string {
	if m != nil {
		return m.VcsId
	}
	return ""
}

type SyncVCSRes struct {
	SyncCompleted string `protobuf:"bytes,1,opt,name=sync_completed,json=syncCompleted" json:"sync_completed,omitempty"`
}

func (m *SyncVCSRes) Reset()                    { *m = SyncVCSRes{} }
func (m *SyncVCSRes) String() string            { return proto.CompactTextString(m) }
func (*SyncVCSRes) ProtoMessage()               {}
func (*SyncVCSRes) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{4} }

func (m *SyncVCSRes) GetSyncCompleted() string {
	if m != nil {
		return m.SyncCompleted
	}
	return ""
}

func init() {
	proto.RegisterType((*VCSAccount)(nil), "api.VCSAccount")
	proto.RegisterType((*GetVCSReq)(nil), "api.GetVCSReq")
	proto.RegisterType((*GetVCSRes)(nil), "api.GetVCSRes")
	proto.RegisterType((*SyncVCSReq)(nil), "api.SyncVCSReq")
	proto.RegisterType((*SyncVCSRes)(nil), "api.SyncVCSRes")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for VCS service

type VCSClient interface {
	GetVCS(ctx context.Context, in *GetVCSReq, opts ...grpc.CallOption) (*GetVCSRes, error)
	Sync(ctx context.Context, in *SyncVCSReq, opts ...grpc.CallOption) (*SyncVCSRes, error)
}

type vCSClient struct {
	cc *grpc.ClientConn
}

func NewVCSClient(cc *grpc.ClientConn) VCSClient {
	return &vCSClient{cc}
}

func (c *vCSClient) GetVCS(ctx context.Context, in *GetVCSReq, opts ...grpc.CallOption) (*GetVCSRes, error) {
	out := new(GetVCSRes)
	err := grpc.Invoke(ctx, "/api.VCS/GetVCS", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vCSClient) Sync(ctx context.Context, in *SyncVCSReq, opts ...grpc.CallOption) (*SyncVCSRes, error) {
	out := new(SyncVCSRes)
	err := grpc.Invoke(ctx, "/api.VCS/Sync", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for VCS service

type VCSServer interface {
	GetVCS(context.Context, *GetVCSReq) (*GetVCSRes, error)
	Sync(context.Context, *SyncVCSReq) (*SyncVCSRes, error)
}

func RegisterVCSServer(s *grpc.Server, srv VCSServer) {
	s.RegisterService(&_VCS_serviceDesc, srv)
}

func _VCS_GetVCS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVCSReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VCSServer).GetVCS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VCS/GetVCS",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VCSServer).GetVCS(ctx, req.(*GetVCSReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _VCS_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncVCSReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VCSServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VCS/Sync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VCSServer).Sync(ctx, req.(*SyncVCSReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _VCS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.VCS",
	HandlerType: (*VCSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVCS",
			Handler:    _VCS_GetVCS_Handler,
		},
		{
			MethodName: "Sync",
			Handler:    _VCS_Sync_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/vcs.proto",
}

func init() { proto.RegisterFile("api/vcs.proto", fileDescriptor4) }

var fileDescriptor4 = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xd1, 0x4a, 0xeb, 0x40,
	0x10, 0x86, 0x49, 0xd2, 0xf6, 0x9c, 0x4c, 0x69, 0x7b, 0x18, 0x38, 0x10, 0x8a, 0x42, 0x59, 0x11,
	0x0a, 0x42, 0x83, 0xed, 0x8d, 0x5e, 0x6a, 0x0a, 0xd2, 0xdb, 0x46, 0x7b, 0x1b, 0xd6, 0xed, 0x52,
	0x02, 0xe9, 0x6e, 0xcc, 0x6e, 0x23, 0x45, 0xbc, 0x29, 0xbe, 0x81, 0x8f, 0xe6, 0x2b, 0xf8, 0x20,
	0x92, 0xdd, 0xb4, 0x55, 0xef, 0x66, 0xbe, 0xff, 0x9f, 0xdd, 0x7f, 0x18, 0xe8, 0xd0, 0x3c, 0x0d,
	0x4b, 0xa6, 0x46, 0x79, 0x21, 0xb5, 0x44, 0x8f, 0xe6, 0x69, 0xff, 0x64, 0x25, 0xe5, 0x2a, 0xe3,
	0x61, 0x25, 0x51, 0x21, 0xa4, 0xa6, 0x3a, 0x95, 0xa2, 0xb6, 0x90, 0x9d, 0x03, 0xb0, 0x88, 0xe2,
	0x1b, 0xc6, 0xe4, 0x46, 0x68, 0xec, 0x82, 0x3b, 0x9b, 0x06, 0xce, 0xc0, 0x19, 0xfa, 0x73, 0x77,
	0x36, 0x45, 0x84, 0x86, 0xa0, 0x6b, 0x1e, 0xb8, 0x86, 0x98, 0xba, 0x62, 0x7a, 0x9b, 0xf3, 0xc0,
	0xb3, 0xac, 0xaa, 0xf1, 0x14, 0x40, 0x3e, 0x0b, 0x5e, 0x24, 0x46, 0x69, 0x18, 0xc5, 0x37, 0xe4,
	0xbe, 0x96, 0x69, 0x49, 0x35, 0x2d, 0x92, 0x4d, 0x91, 0x05, 0x4d, 0x2b, 0x5b, 0xf2, 0x50, 0x64,
	0xa4, 0x0d, 0xfe, 0x1d, 0xd7, 0x8b, 0x28, 0x9e, 0xf3, 0x27, 0x72, 0x75, 0x6c, 0x14, 0x5e, 0xc0,
	0x5f, 0x6a, 0xa3, 0xa9, 0xc0, 0x19, 0x78, 0xc3, 0xf6, 0xb8, 0x37, 0xa2, 0x79, 0x3a, 0x3a, 0x46,
	0x9e, 0x1f, 0x0c, 0xe4, 0x0c, 0x20, 0xde, 0x0a, 0x66, 0xdf, 0xc1, 0xff, 0xd0, 0x2a, 0x99, 0x4a,
	0xd2, 0x65, 0xbd, 0x4e, 0xb3, 0x64, 0x6a, 0xb6, 0x24, 0x93, 0x6f, 0x26, 0x85, 0xe7, 0xd0, 0x55,
	0x5b, 0xc1, 0x12, 0x26, 0xd7, 0x79, 0xc6, 0x35, 0xdf, 0x9b, 0x3b, 0x15, 0x8d, 0xf6, 0x70, 0xfc,
	0xe6, 0x80, 0xb7, 0x88, 0x62, 0xbc, 0x86, 0x96, 0xcd, 0x86, 0x5d, 0x13, 0xe3, 0x90, 0xba, 0xff,
	0xb3, 0x57, 0xa4, 0xb7, 0xfb, 0xf8, 0x7c, 0x77, 0x7d, 0xfc, 0x13, 0x96, 0x97, 0xd5, 0x45, 0xf0,
	0x16, 0x1a, 0xd5, 0xbf, 0x68, 0xf3, 0x1f, 0x73, 0xf6, 0x7f, 0x01, 0x45, 0x02, 0x33, 0x8a, 0xf8,
	0xaf, 0x1e, 0x0d, 0x5f, 0xec, 0x1e, 0xaf, 0x8f, 0x2d, 0x73, 0xb3, 0xc9, 0x57, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x86, 0x1a, 0xc6, 0x03, 0xe7, 0x01, 0x00, 0x00,
}
