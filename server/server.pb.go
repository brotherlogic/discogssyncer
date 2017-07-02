// Code generated by protoc-gen-go.
// source: server.proto
// DO NOT EDIT!

/*
Package discogsserver is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
	Token
	RecordCollection
	CollectionFolder
	ReleaseMetadata
	Empty
	FolderList
	ReleaseList
	ReleaseMove
	MetadataUpdate
	Want
	Wantlist
	SpendRequest
	SpendResponse
	SearchRequest
*/
package discogsserver

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

type Token struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type RecordCollection struct {
	Folders  []*CollectionFolder `protobuf:"bytes,1,rep,name=folders" json:"folders,omitempty"`
	Metadata []*ReleaseMetadata  `protobuf:"bytes,2,rep,name=metadata" json:"metadata,omitempty"`
	Wantlist *Wantlist           `protobuf:"bytes,3,opt,name=wantlist" json:"wantlist,omitempty"`
}

func (m *RecordCollection) Reset()                    { *m = RecordCollection{} }
func (m *RecordCollection) String() string            { return proto.CompactTextString(m) }
func (*RecordCollection) ProtoMessage()               {}
func (*RecordCollection) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RecordCollection) GetFolders() []*CollectionFolder {
	if m != nil {
		return m.Folders
	}
	return nil
}

func (m *RecordCollection) GetMetadata() []*ReleaseMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *RecordCollection) GetWantlist() *Wantlist {
	if m != nil {
		return m.Wantlist
	}
	return nil
}

type CollectionFolder struct {
	Folder   *godiscogs.Folder `protobuf:"bytes,1,opt,name=folder" json:"folder,omitempty"`
	Releases *ReleaseList      `protobuf:"bytes,2,opt,name=releases" json:"releases,omitempty"`
}

func (m *CollectionFolder) Reset()                    { *m = CollectionFolder{} }
func (m *CollectionFolder) String() string            { return proto.CompactTextString(m) }
func (*CollectionFolder) ProtoMessage()               {}
func (*CollectionFolder) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CollectionFolder) GetFolder() *godiscogs.Folder {
	if m != nil {
		return m.Folder
	}
	return nil
}

func (m *CollectionFolder) GetReleases() *ReleaseList {
	if m != nil {
		return m.Releases
	}
	return nil
}

type ReleaseMetadata struct {
	// The date the release was added
	DateAdded int64 `protobuf:"varint,1,opt,name=date_added,json=dateAdded" json:"date_added,omitempty"`
	// The date the release was last refreshed
	DateRefreshed int64 `protobuf:"varint,2,opt,name=date_refreshed,json=dateRefreshed" json:"date_refreshed,omitempty"`
	// The path to the file on iTunes if available
	FilePath string `protobuf:"bytes,3,opt,name=file_path,json=filePath" json:"file_path,omitempty"`
	// The cost of the record in pence
	Cost int32 `protobuf:"varint,4,opt,name=cost" json:"cost,omitempty"`
	// If we have other copies of this
	Others bool `protobuf:"varint,5,opt,name=others" json:"others,omitempty"`
	// The id of the release this relates to
	Id int32 `protobuf:"varint,6,opt,name=id" json:"id,omitempty"`
}

func (m *ReleaseMetadata) Reset()                    { *m = ReleaseMetadata{} }
func (m *ReleaseMetadata) String() string            { return proto.CompactTextString(m) }
func (*ReleaseMetadata) ProtoMessage()               {}
func (*ReleaseMetadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type FolderList struct {
	Folders []*godiscogs.Folder `protobuf:"bytes,1,rep,name=folders" json:"folders,omitempty"`
}

func (m *FolderList) Reset()                    { *m = FolderList{} }
func (m *FolderList) String() string            { return proto.CompactTextString(m) }
func (*FolderList) ProtoMessage()               {}
func (*FolderList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *FolderList) GetFolders() []*godiscogs.Folder {
	if m != nil {
		return m.Folders
	}
	return nil
}

type ReleaseList struct {
	Releases []*godiscogs.Release `protobuf:"bytes,1,rep,name=releases" json:"releases,omitempty"`
}

func (m *ReleaseList) Reset()                    { *m = ReleaseList{} }
func (m *ReleaseList) String() string            { return proto.CompactTextString(m) }
func (*ReleaseList) ProtoMessage()               {}
func (*ReleaseList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ReleaseList) GetReleases() []*godiscogs.Release {
	if m != nil {
		return m.Releases
	}
	return nil
}

type ReleaseMove struct {
	Release     *godiscogs.Release `protobuf:"bytes,1,opt,name=release" json:"release,omitempty"`
	NewFolderId int32              `protobuf:"varint,2,opt,name=new_folder_id,json=newFolderId" json:"new_folder_id,omitempty"`
}

func (m *ReleaseMove) Reset()                    { *m = ReleaseMove{} }
func (m *ReleaseMove) String() string            { return proto.CompactTextString(m) }
func (*ReleaseMove) ProtoMessage()               {}
func (*ReleaseMove) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ReleaseMove) GetRelease() *godiscogs.Release {
	if m != nil {
		return m.Release
	}
	return nil
}

type MetadataUpdate struct {
	Release *godiscogs.Release `protobuf:"bytes,1,opt,name=release" json:"release,omitempty"`
	Update  *ReleaseMetadata   `protobuf:"bytes,2,opt,name=update" json:"update,omitempty"`
}

func (m *MetadataUpdate) Reset()                    { *m = MetadataUpdate{} }
func (m *MetadataUpdate) String() string            { return proto.CompactTextString(m) }
func (*MetadataUpdate) ProtoMessage()               {}
func (*MetadataUpdate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *MetadataUpdate) GetRelease() *godiscogs.Release {
	if m != nil {
		return m.Release
	}
	return nil
}

func (m *MetadataUpdate) GetUpdate() *ReleaseMetadata {
	if m != nil {
		return m.Update
	}
	return nil
}

type Want struct {
	ReleaseId int32 `protobuf:"varint,1,opt,name=release_id,json=releaseId" json:"release_id,omitempty"`
	Valued    bool  `protobuf:"varint,2,opt,name=valued" json:"valued,omitempty"`
	Wanted    bool  `protobuf:"varint,3,opt,name=wanted" json:"wanted,omitempty"`
}

func (m *Want) Reset()                    { *m = Want{} }
func (m *Want) String() string            { return proto.CompactTextString(m) }
func (*Want) ProtoMessage()               {}
func (*Want) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type Wantlist struct {
	Want []*Want `protobuf:"bytes,1,rep,name=want" json:"want,omitempty"`
}

func (m *Wantlist) Reset()                    { *m = Wantlist{} }
func (m *Wantlist) String() string            { return proto.CompactTextString(m) }
func (*Wantlist) ProtoMessage()               {}
func (*Wantlist) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Wantlist) GetWant() []*Want {
	if m != nil {
		return m.Want
	}
	return nil
}

type SpendRequest struct {
	Month int32 `protobuf:"varint,1,opt,name=month" json:"month,omitempty"`
	Year  int32 `protobuf:"varint,2,opt,name=year" json:"year,omitempty"`
	Lower int64 `protobuf:"varint,3,opt,name=lower" json:"lower,omitempty"`
	Upper int64 `protobuf:"varint,4,opt,name=upper" json:"upper,omitempty"`
}

func (m *SpendRequest) Reset()                    { *m = SpendRequest{} }
func (m *SpendRequest) String() string            { return proto.CompactTextString(m) }
func (*SpendRequest) ProtoMessage()               {}
func (*SpendRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type SpendResponse struct {
	TotalSpend int32             `protobuf:"varint,1,opt,name=total_spend,json=totalSpend" json:"total_spend,omitempty"`
	Spends     []*MetadataUpdate `protobuf:"bytes,2,rep,name=spends" json:"spends,omitempty"`
}

func (m *SpendResponse) Reset()                    { *m = SpendResponse{} }
func (m *SpendResponse) String() string            { return proto.CompactTextString(m) }
func (*SpendResponse) ProtoMessage()               {}
func (*SpendResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *SpendResponse) GetSpends() []*MetadataUpdate {
	if m != nil {
		return m.Spends
	}
	return nil
}

type SearchRequest struct {
	Query string `protobuf:"bytes,1,opt,name=query" json:"query,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func init() {
	proto.RegisterType((*Token)(nil), "discogsserver.Token")
	proto.RegisterType((*RecordCollection)(nil), "discogsserver.RecordCollection")
	proto.RegisterType((*CollectionFolder)(nil), "discogsserver.CollectionFolder")
	proto.RegisterType((*ReleaseMetadata)(nil), "discogsserver.ReleaseMetadata")
	proto.RegisterType((*Empty)(nil), "discogsserver.Empty")
	proto.RegisterType((*FolderList)(nil), "discogsserver.FolderList")
	proto.RegisterType((*ReleaseList)(nil), "discogsserver.ReleaseList")
	proto.RegisterType((*ReleaseMove)(nil), "discogsserver.ReleaseMove")
	proto.RegisterType((*MetadataUpdate)(nil), "discogsserver.MetadataUpdate")
	proto.RegisterType((*Want)(nil), "discogsserver.Want")
	proto.RegisterType((*Wantlist)(nil), "discogsserver.Wantlist")
	proto.RegisterType((*SpendRequest)(nil), "discogsserver.SpendRequest")
	proto.RegisterType((*SpendResponse)(nil), "discogsserver.SpendResponse")
	proto.RegisterType((*SearchRequest)(nil), "discogsserver.SearchRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DiscogsService service

type DiscogsServiceClient interface {
	GetCollection(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReleaseList, error)
	GetReleasesInFolder(ctx context.Context, in *FolderList, opts ...grpc.CallOption) (*ReleaseList, error)
	MoveToFolder(ctx context.Context, in *ReleaseMove, opts ...grpc.CallOption) (*Empty, error)
	AddToFolder(ctx context.Context, in *ReleaseMove, opts ...grpc.CallOption) (*Empty, error)
	UpdateMetadata(ctx context.Context, in *MetadataUpdate, opts ...grpc.CallOption) (*ReleaseMetadata, error)
	GetMetadata(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*ReleaseMetadata, error)
	UpdateRating(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*Empty, error)
	GetSingleRelease(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*godiscogs.Release, error)
	GetWantlist(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Wantlist, error)
	CollapseWantlist(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Wantlist, error)
	RebuildWantlist(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Wantlist, error)
	GetSpend(ctx context.Context, in *SpendRequest, opts ...grpc.CallOption) (*SpendResponse, error)
	EditWant(ctx context.Context, in *Want, opts ...grpc.CallOption) (*Want, error)
	DeleteWant(ctx context.Context, in *Want, opts ...grpc.CallOption) (*Wantlist, error)
	AddWant(ctx context.Context, in *Want, opts ...grpc.CallOption) (*Empty, error)
	SyncWithDiscogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type discogsServiceClient struct {
	cc *grpc.ClientConn
}

func NewDiscogsServiceClient(cc *grpc.ClientConn) DiscogsServiceClient {
	return &discogsServiceClient{cc}
}

func (c *discogsServiceClient) GetCollection(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReleaseList, error) {
	out := new(ReleaseList)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/GetCollection", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) GetReleasesInFolder(ctx context.Context, in *FolderList, opts ...grpc.CallOption) (*ReleaseList, error) {
	out := new(ReleaseList)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/GetReleasesInFolder", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) MoveToFolder(ctx context.Context, in *ReleaseMove, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/MoveToFolder", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) AddToFolder(ctx context.Context, in *ReleaseMove, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/AddToFolder", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) UpdateMetadata(ctx context.Context, in *MetadataUpdate, opts ...grpc.CallOption) (*ReleaseMetadata, error) {
	out := new(ReleaseMetadata)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/UpdateMetadata", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) GetMetadata(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*ReleaseMetadata, error) {
	out := new(ReleaseMetadata)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/GetMetadata", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) UpdateRating(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/UpdateRating", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) GetSingleRelease(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*godiscogs.Release, error) {
	out := new(godiscogs.Release)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/GetSingleRelease", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) GetWantlist(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Wantlist, error) {
	out := new(Wantlist)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/GetWantlist", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) CollapseWantlist(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Wantlist, error) {
	out := new(Wantlist)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/CollapseWantlist", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) RebuildWantlist(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Wantlist, error) {
	out := new(Wantlist)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/RebuildWantlist", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) GetSpend(ctx context.Context, in *SpendRequest, opts ...grpc.CallOption) (*SpendResponse, error) {
	out := new(SpendResponse)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/GetSpend", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) EditWant(ctx context.Context, in *Want, opts ...grpc.CallOption) (*Want, error) {
	out := new(Want)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/EditWant", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) DeleteWant(ctx context.Context, in *Want, opts ...grpc.CallOption) (*Wantlist, error) {
	out := new(Wantlist)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/DeleteWant", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) AddWant(ctx context.Context, in *Want, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/AddWant", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discogsServiceClient) SyncWithDiscogs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/discogsserver.DiscogsService/SyncWithDiscogs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DiscogsService service

type DiscogsServiceServer interface {
	GetCollection(context.Context, *Empty) (*ReleaseList, error)
	GetReleasesInFolder(context.Context, *FolderList) (*ReleaseList, error)
	MoveToFolder(context.Context, *ReleaseMove) (*Empty, error)
	AddToFolder(context.Context, *ReleaseMove) (*Empty, error)
	UpdateMetadata(context.Context, *MetadataUpdate) (*ReleaseMetadata, error)
	GetMetadata(context.Context, *godiscogs.Release) (*ReleaseMetadata, error)
	UpdateRating(context.Context, *godiscogs.Release) (*Empty, error)
	GetSingleRelease(context.Context, *godiscogs.Release) (*godiscogs.Release, error)
	GetWantlist(context.Context, *Empty) (*Wantlist, error)
	CollapseWantlist(context.Context, *Empty) (*Wantlist, error)
	RebuildWantlist(context.Context, *Empty) (*Wantlist, error)
	GetSpend(context.Context, *SpendRequest) (*SpendResponse, error)
	EditWant(context.Context, *Want) (*Want, error)
	DeleteWant(context.Context, *Want) (*Wantlist, error)
	AddWant(context.Context, *Want) (*Empty, error)
	SyncWithDiscogs(context.Context, *Empty) (*Empty, error)
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
		FullMethod: "/discogsserver.DiscogsService/GetCollection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).GetCollection(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_GetReleasesInFolder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FolderList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).GetReleasesInFolder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/GetReleasesInFolder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).GetReleasesInFolder(ctx, req.(*FolderList))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_MoveToFolder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseMove)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).MoveToFolder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/MoveToFolder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).MoveToFolder(ctx, req.(*ReleaseMove))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_AddToFolder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseMove)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).AddToFolder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/AddToFolder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).AddToFolder(ctx, req.(*ReleaseMove))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_UpdateMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetadataUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).UpdateMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/UpdateMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).UpdateMetadata(ctx, req.(*MetadataUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_GetMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(godiscogs.Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).GetMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/GetMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).GetMetadata(ctx, req.(*godiscogs.Release))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_UpdateRating_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(godiscogs.Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).UpdateRating(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/UpdateRating",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).UpdateRating(ctx, req.(*godiscogs.Release))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_GetSingleRelease_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(godiscogs.Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).GetSingleRelease(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/GetSingleRelease",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).GetSingleRelease(ctx, req.(*godiscogs.Release))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_GetWantlist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).GetWantlist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/GetWantlist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).GetWantlist(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_CollapseWantlist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).CollapseWantlist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/CollapseWantlist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).CollapseWantlist(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_RebuildWantlist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).RebuildWantlist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/RebuildWantlist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).RebuildWantlist(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_GetSpend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).GetSpend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/GetSpend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).GetSpend(ctx, req.(*SpendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_EditWant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Want)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).EditWant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/EditWant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).EditWant(ctx, req.(*Want))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_DeleteWant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Want)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).DeleteWant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/DeleteWant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).DeleteWant(ctx, req.(*Want))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_AddWant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Want)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).AddWant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/AddWant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).AddWant(ctx, req.(*Want))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiscogsService_SyncWithDiscogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscogsServiceServer).SyncWithDiscogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discogsserver.DiscogsService/SyncWithDiscogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscogsServiceServer).SyncWithDiscogs(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _DiscogsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discogsserver.DiscogsService",
	HandlerType: (*DiscogsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCollection",
			Handler:    _DiscogsService_GetCollection_Handler,
		},
		{
			MethodName: "GetReleasesInFolder",
			Handler:    _DiscogsService_GetReleasesInFolder_Handler,
		},
		{
			MethodName: "MoveToFolder",
			Handler:    _DiscogsService_MoveToFolder_Handler,
		},
		{
			MethodName: "AddToFolder",
			Handler:    _DiscogsService_AddToFolder_Handler,
		},
		{
			MethodName: "UpdateMetadata",
			Handler:    _DiscogsService_UpdateMetadata_Handler,
		},
		{
			MethodName: "GetMetadata",
			Handler:    _DiscogsService_GetMetadata_Handler,
		},
		{
			MethodName: "UpdateRating",
			Handler:    _DiscogsService_UpdateRating_Handler,
		},
		{
			MethodName: "GetSingleRelease",
			Handler:    _DiscogsService_GetSingleRelease_Handler,
		},
		{
			MethodName: "GetWantlist",
			Handler:    _DiscogsService_GetWantlist_Handler,
		},
		{
			MethodName: "CollapseWantlist",
			Handler:    _DiscogsService_CollapseWantlist_Handler,
		},
		{
			MethodName: "RebuildWantlist",
			Handler:    _DiscogsService_RebuildWantlist_Handler,
		},
		{
			MethodName: "GetSpend",
			Handler:    _DiscogsService_GetSpend_Handler,
		},
		{
			MethodName: "EditWant",
			Handler:    _DiscogsService_EditWant_Handler,
		},
		{
			MethodName: "DeleteWant",
			Handler:    _DiscogsService_DeleteWant_Handler,
		},
		{
			MethodName: "AddWant",
			Handler:    _DiscogsService_AddWant_Handler,
		},
		{
			MethodName: "SyncWithDiscogs",
			Handler:    _DiscogsService_SyncWithDiscogs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}

func init() { proto.RegisterFile("server.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 882 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x56, 0x6d, 0x6f, 0xdb, 0x36,
	0x10, 0xb6, 0xe2, 0xd8, 0x91, 0xcf, 0xb1, 0x9b, 0xb1, 0xc5, 0xe6, 0xb9, 0xcb, 0x6a, 0x10, 0x28,
	0xe6, 0x61, 0x83, 0x8b, 0x25, 0x58, 0x80, 0x16, 0x7b, 0x41, 0x9a, 0x76, 0x46, 0x80, 0x75, 0xd8,
	0xe8, 0x16, 0xfd, 0x68, 0x28, 0xe2, 0xc5, 0x16, 0xa6, 0x88, 0x2a, 0x49, 0x3b, 0xc8, 0xbf, 0xda,
	0x2f, 0xd8, 0xfe, 0xda, 0xc0, 0x17, 0xd9, 0xae, 0x20, 0x27, 0x5b, 0xfb, 0x4d, 0x3c, 0x3e, 0x77,
	0x7c, 0xf8, 0xdc, 0x0b, 0x05, 0xfb, 0x0a, 0xe5, 0x12, 0xe5, 0x28, 0x97, 0x42, 0x0b, 0xd2, 0xe1,
	0x89, 0x8a, 0xc5, 0x4c, 0x39, 0x63, 0xff, 0xbb, 0x59, 0xa2, 0xe7, 0x8b, 0x8b, 0x51, 0x2c, 0xae,
	0x9e, 0x5c, 0x48, 0xa1, 0xe7, 0x28, 0x53, 0x31, 0x4b, 0xe2, 0x27, 0x33, 0xe1, 0x81, 0xeb, 0x2f,
	0x17, 0x81, 0x1e, 0x42, 0xe3, 0xb5, 0xf8, 0x13, 0x33, 0xf2, 0x00, 0x1a, 0xda, 0x7c, 0xf4, 0x82,
	0x41, 0x30, 0x6c, 0x31, 0xb7, 0xa0, 0x7f, 0x07, 0x70, 0xc0, 0x30, 0x16, 0x92, 0x9f, 0x89, 0x34,
	0xc5, 0x58, 0x27, 0x22, 0x23, 0x4f, 0x61, 0xef, 0x52, 0xa4, 0x1c, 0xa5, 0xea, 0x05, 0x83, 0xfa,
	0xb0, 0x7d, 0xf4, 0x68, 0xf4, 0x1e, 0x8f, 0xd1, 0x1a, 0xfb, 0x8b, 0xc5, 0xb1, 0x02, 0x4f, 0x9e,
	0x41, 0x78, 0x85, 0x3a, 0xe2, 0x91, 0x8e, 0x7a, 0x3b, 0xd6, 0xf7, 0xcb, 0x92, 0x2f, 0xc3, 0x14,
	0x23, 0x85, 0xaf, 0x3c, 0x8a, 0xad, 0xf0, 0xe4, 0x18, 0xc2, 0xeb, 0x28, 0xd3, 0x69, 0xa2, 0x74,
	0xaf, 0x3e, 0x08, 0x86, 0xed, 0xa3, 0xcf, 0x4a, 0xbe, 0x6f, 0xfd, 0x36, 0x5b, 0x01, 0xe9, 0x02,
	0x0e, 0xca, 0x6c, 0xc8, 0xd7, 0xd0, 0x74, 0x7c, 0xec, 0x5d, 0xdb, 0x47, 0x9f, 0x8c, 0xd6, 0xaa,
	0x78, 0xc2, 0x1e, 0x40, 0x4e, 0x20, 0x94, 0x8e, 0x90, 0xea, 0xed, 0x58, 0x70, 0xbf, 0x9a, 0xef,
	0xaf, 0xf6, 0xd8, 0x02, 0x4b, 0xff, 0x0a, 0xe0, 0x5e, 0xe9, 0x26, 0xe4, 0x10, 0x80, 0x47, 0x1a,
	0xa7, 0x11, 0xe7, 0xc8, 0xed, 0xd1, 0x75, 0xd6, 0x32, 0x96, 0x53, 0x63, 0x20, 0x8f, 0xa1, 0x6b,
	0xb7, 0x25, 0x5e, 0x4a, 0x54, 0x73, 0xe4, 0xf6, 0xc0, 0x3a, 0xeb, 0x18, 0x2b, 0x2b, 0x8c, 0xe4,
	0x21, 0xb4, 0x2e, 0x93, 0x14, 0xa7, 0x79, 0xa4, 0xe7, 0x56, 0x86, 0x16, 0x0b, 0x8d, 0xe1, 0xf7,
	0x48, 0xcf, 0x09, 0x81, 0xdd, 0x58, 0x28, 0xdd, 0xdb, 0x1d, 0x04, 0xc3, 0x06, 0xb3, 0xdf, 0xe4,
	0x53, 0x68, 0xda, 0x4a, 0x50, 0xbd, 0xc6, 0x20, 0x18, 0x86, 0xcc, 0xaf, 0x48, 0x17, 0x76, 0x12,
	0xde, 0x6b, 0x5a, 0xe4, 0x4e, 0xc2, 0xe9, 0x1e, 0x34, 0x5e, 0x5e, 0xe5, 0xfa, 0x86, 0x3e, 0x05,
	0x70, 0x2a, 0x98, 0x3b, 0x91, 0x6f, 0xca, 0xc9, 0xae, 0x50, 0xab, 0x40, 0xd0, 0x1f, 0xa1, 0xbd,
	0xa1, 0x07, 0x19, 0x6d, 0xa8, 0xe7, 0x9c, 0xc9, 0x86, 0xb3, 0x47, 0x6e, 0xa8, 0x36, 0x5d, 0xb9,
	0xbf, 0x12, 0x4b, 0x24, 0xdf, 0xc2, 0x9e, 0xdf, 0xf2, 0x89, 0xaa, 0xf2, 0x2e, 0x20, 0x84, 0x42,
	0x27, 0xc3, 0xeb, 0xa9, 0xa3, 0x32, 0x4d, 0x9c, 0x7c, 0x0d, 0xd6, 0xce, 0xf0, 0xda, 0xd1, 0x3c,
	0xe7, 0x74, 0x09, 0xdd, 0x22, 0x1d, 0x6f, 0x72, 0xa3, 0xeb, 0xff, 0x3c, 0xe3, 0x04, 0x9a, 0x0b,
	0xeb, 0xe7, 0x8b, 0xe1, 0xae, 0xe2, 0xf5, 0x68, 0xfa, 0x06, 0x76, 0x4d, 0x6d, 0x9a, 0x12, 0xf0,
	0xa1, 0x0c, 0xc1, 0xc0, 0x12, 0x6c, 0x79, 0xcb, 0x39, 0x37, 0xa9, 0x5a, 0x46, 0xe9, 0xc2, 0xa7,
	0x3e, 0x64, 0x7e, 0x65, 0xec, 0xa6, 0xa0, 0x91, 0xdb, 0x84, 0x87, 0xcc, 0xaf, 0xe8, 0x31, 0x84,
	0x45, 0xc9, 0x93, 0xaf, 0x60, 0xd7, 0x58, 0xbd, 0xce, 0xf7, 0x2b, 0x3a, 0x83, 0x59, 0x00, 0xe5,
	0xb0, 0x3f, 0xc9, 0x31, 0xe3, 0x0c, 0xdf, 0x2d, 0x50, 0x69, 0xd3, 0xf8, 0x57, 0x22, 0xd3, 0x73,
	0x4f, 0xc7, 0x2d, 0x4c, 0x25, 0xdd, 0x60, 0x24, 0xbd, 0x88, 0xf6, 0xdb, 0x20, 0x53, 0x71, 0x8d,
	0xd2, 0xb2, 0xa8, 0x33, 0xb7, 0x30, 0xd6, 0x45, 0x9e, 0xa3, 0xb4, 0x45, 0x57, 0x67, 0x6e, 0x41,
	0x67, 0xd0, 0xf1, 0xa7, 0xa8, 0x5c, 0x64, 0x0a, 0xc9, 0x23, 0x68, 0x6b, 0xa1, 0xa3, 0x74, 0xaa,
	0x8c, 0xd9, 0x1f, 0x06, 0xd6, 0x64, 0x81, 0xe4, 0x7b, 0x68, 0xda, 0x2d, 0xe5, 0x07, 0xc3, 0x61,
	0xe9, 0x0a, 0xef, 0x27, 0x8e, 0x79, 0x30, 0x7d, 0x0c, 0x9d, 0x09, 0x46, 0x32, 0x9e, 0x6f, 0xdc,
	0xe7, 0xdd, 0x02, 0xe5, 0x4d, 0x31, 0xc8, 0xec, 0xe2, 0xe8, 0x9f, 0x10, 0xba, 0x2f, 0x5c, 0xbc,
	0x09, 0xca, 0x65, 0x12, 0x23, 0x39, 0x83, 0xce, 0x18, 0xf5, 0xc6, 0x5c, 0x7b, 0x50, 0x3a, 0xd1,
	0xb6, 0x43, 0xff, 0x96, 0x86, 0xa7, 0x35, 0xf2, 0x1b, 0xdc, 0x1f, 0xa3, 0xf6, 0x36, 0x75, 0x5e,
	0x8c, 0x98, 0xcf, 0x4b, 0x4e, 0xeb, 0x86, 0xba, 0x23, 0xde, 0x73, 0xd8, 0x37, 0xb5, 0xff, 0x5a,
	0xf8, 0x40, 0x5b, 0xd0, 0x06, 0xd3, 0xaf, 0xe4, 0x4b, 0x6b, 0xe4, 0x14, 0xda, 0xa7, 0x9c, 0x7f,
	0x54, 0x88, 0x3f, 0xa0, 0xeb, 0x74, 0x5e, 0x4f, 0xaf, 0x5b, 0xd3, 0xd1, 0xbf, 0xa3, 0x13, 0x68,
	0x8d, 0x9c, 0x41, 0x7b, 0x8c, 0x7a, 0x15, 0xaf, 0xa2, 0xcf, 0xfe, 0x43, 0x90, 0x67, 0xb0, 0xef,
	0xf3, 0x1f, 0xe9, 0x24, 0x9b, 0x55, 0x46, 0xd9, 0x76, 0xa7, 0x1f, 0xe0, 0x60, 0x8c, 0x7a, 0x92,
	0x64, 0xb3, 0x14, 0x3d, 0xb6, 0xd2, 0xbf, 0xc2, 0x46, 0x6b, 0xe4, 0x27, 0x4b, 0x7f, 0xd5, 0x6e,
	0xd5, 0xb5, 0xb2, 0xed, 0x41, 0xb2, 0xd7, 0xb7, 0x0f, 0x51, 0x94, 0x2b, 0xfc, 0xf0, 0x20, 0xcf,
	0xcd, 0xab, 0x72, 0xb1, 0x48, 0x52, 0xfe, 0xe1, 0x31, 0xc6, 0x10, 0x1a, 0x19, 0x6c, 0xcf, 0x3d,
	0x2c, 0xc1, 0x36, 0x07, 0x43, 0xff, 0x8b, 0xea, 0x4d, 0xd7, 0xcf, 0xb4, 0x66, 0xde, 0xc6, 0x97,
	0x3c, 0xb1, 0x92, 0x90, 0xaa, 0x79, 0xd3, 0xaf, 0x32, 0xda, 0x3c, 0xc0, 0x0b, 0x4c, 0x51, 0xe3,
	0x76, 0xcf, 0x5b, 0xe8, 0x9f, 0xc0, 0xde, 0x29, 0xe7, 0xdb, 0x5d, 0xb7, 0x65, 0xff, 0x67, 0xb8,
	0x37, 0xb9, 0xc9, 0xe2, 0xb7, 0x89, 0x9e, 0xfb, 0x39, 0xb0, 0x45, 0xba, 0x2d, 0x01, 0x2e, 0x9a,
	0xf6, 0x87, 0xe9, 0xf8, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x1f, 0x81, 0x67, 0x82, 0x09,
	0x00, 0x00,
}
