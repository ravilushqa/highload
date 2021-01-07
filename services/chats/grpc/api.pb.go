// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.13.0
// source: services/chats/grpc/api.proto

package grpc

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type GetUserChatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetUserChatsRequest) Reset() {
	*x = GetUserChatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserChatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserChatsRequest) ProtoMessage() {}

func (x *GetUserChatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserChatsRequest.ProtoReflect.Descriptor instead.
func (*GetUserChatsRequest) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserChatsRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetUserChatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatIds []int64 `protobuf:"varint,1,rep,packed,name=chat_ids,json=chatIds,proto3" json:"chat_ids,omitempty"`
}

func (x *GetUserChatsResponse) Reset() {
	*x = GetUserChatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserChatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserChatsResponse) ProtoMessage() {}

func (x *GetUserChatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserChatsResponse.ProtoReflect.Descriptor instead.
func (*GetUserChatsResponse) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserChatsResponse) GetChatIds() []int64 {
	if x != nil {
		return x.ChatIds
	}
	return nil
}

type GetChatMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatId int64 `protobuf:"varint,1,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
}

func (x *GetChatMessagesRequest) Reset() {
	*x = GetChatMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChatMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChatMessagesRequest) ProtoMessage() {}

func (x *GetChatMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChatMessagesRequest.ProtoReflect.Descriptor instead.
func (*GetChatMessagesRequest) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{2}
}

func (x *GetChatMessagesRequest) GetChatId() int64 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

type GetChatMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*Message `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *GetChatMessagesResponse) Reset() {
	*x = GetChatMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChatMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChatMessagesResponse) ProtoMessage() {}

func (x *GetChatMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChatMessagesResponse.ProtoReflect.Descriptor instead.
func (*GetChatMessagesResponse) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{3}
}

func (x *GetChatMessagesResponse) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

type StoreMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ChatId int64  `protobuf:"varint,2,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
	Text   string `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *StoreMessageRequest) Reset() {
	*x = StoreMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreMessageRequest) ProtoMessage() {}

func (x *StoreMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreMessageRequest.ProtoReflect.Descriptor instead.
func (*StoreMessageRequest) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{4}
}

func (x *StoreMessageRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *StoreMessageRequest) GetChatId() int64 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

func (x *StoreMessageRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type FindOrCreateChatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId_1 int64 `protobuf:"varint,1,opt,name=user_id_1,json=userId1,proto3" json:"user_id_1,omitempty"`
	UserId_2 int64 `protobuf:"varint,2,opt,name=user_id_2,json=userId2,proto3" json:"user_id_2,omitempty"`
}

func (x *FindOrCreateChatRequest) Reset() {
	*x = FindOrCreateChatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindOrCreateChatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindOrCreateChatRequest) ProtoMessage() {}

func (x *FindOrCreateChatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindOrCreateChatRequest.ProtoReflect.Descriptor instead.
func (*FindOrCreateChatRequest) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{5}
}

func (x *FindOrCreateChatRequest) GetUserId_1() int64 {
	if x != nil {
		return x.UserId_1
	}
	return 0
}

func (x *FindOrCreateChatRequest) GetUserId_2() int64 {
	if x != nil {
		return x.UserId_2
	}
	return 0
}

type FindOrCreateChatResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatId int64 `protobuf:"varint,1,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
}

func (x *FindOrCreateChatResponse) Reset() {
	*x = FindOrCreateChatResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindOrCreateChatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindOrCreateChatResponse) ProtoMessage() {}

func (x *FindOrCreateChatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindOrCreateChatResponse.ProtoReflect.Descriptor instead.
func (*FindOrCreateChatResponse) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{6}
}

func (x *FindOrCreateChatResponse) GetChatId() int64 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid      string               `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	UserId    int64                `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ChatId    int64                `protobuf:"varint,3,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
	Text      string               `protobuf:"bytes,4,opt,name=text,proto3" json:"text,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt *timestamp.Timestamp `protobuf:"bytes,7,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_chats_grpc_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_services_chats_grpc_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_services_chats_grpc_api_proto_rawDescGZIP(), []int{7}
}

func (x *Message) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Message) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Message) GetChatId() int64 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

func (x *Message) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Message) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Message) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Message) GetDeletedAt() *timestamp.Timestamp {
	if x != nil {
		return x.DeletedAt
	}
	return nil
}

var File_services_chats_grpc_api_proto protoreflect.FileDescriptor

var file_services_chats_grpc_api_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x73,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2e, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43,
	0x68, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x22, 0x31, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43,
	0x68, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08,
	0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x07,
	0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x73, 0x22, 0x31, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x43, 0x68,
	0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x22, 0x45, 0x0a, 0x17, 0x47, 0x65,
	0x74, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x22, 0x5b, 0x0a, 0x13, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x51,
	0x0a, 0x17, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68,
	0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x09, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x5f, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x31, 0x12, 0x1a, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x5f, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x32, 0x22, 0x33, 0x0a, 0x18, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x22, 0x94, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x39, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x32, 0xbb, 0x02,
	0x0a, 0x05, 0x43, 0x68, 0x61, 0x74, 0x73, 0x12, 0x47, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x43, 0x68, 0x61, 0x74, 0x73, 0x12, 0x1a, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x43, 0x68, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x50, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x43,
	0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68,
	0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x42, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x1a, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x53, 0x0a, 0x10, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x72,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x12, 0x1e, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x73, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x73, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43,
	0x68, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x15, 0x5a, 0x13, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_chats_grpc_api_proto_rawDescOnce sync.Once
	file_services_chats_grpc_api_proto_rawDescData = file_services_chats_grpc_api_proto_rawDesc
)

func file_services_chats_grpc_api_proto_rawDescGZIP() []byte {
	file_services_chats_grpc_api_proto_rawDescOnce.Do(func() {
		file_services_chats_grpc_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_chats_grpc_api_proto_rawDescData)
	})
	return file_services_chats_grpc_api_proto_rawDescData
}

var file_services_chats_grpc_api_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_services_chats_grpc_api_proto_goTypes = []interface{}{
	(*GetUserChatsRequest)(nil),      // 0: chats.GetUserChatsRequest
	(*GetUserChatsResponse)(nil),     // 1: chats.GetUserChatsResponse
	(*GetChatMessagesRequest)(nil),   // 2: chats.GetChatMessagesRequest
	(*GetChatMessagesResponse)(nil),  // 3: chats.GetChatMessagesResponse
	(*StoreMessageRequest)(nil),      // 4: chats.StoreMessageRequest
	(*FindOrCreateChatRequest)(nil),  // 5: chats.FindOrCreateChatRequest
	(*FindOrCreateChatResponse)(nil), // 6: chats.FindOrCreateChatResponse
	(*Message)(nil),                  // 7: chats.Message
	(*timestamp.Timestamp)(nil),      // 8: google.protobuf.Timestamp
	(*empty.Empty)(nil),              // 9: google.protobuf.Empty
}
var file_services_chats_grpc_api_proto_depIdxs = []int32{
	7, // 0: chats.GetChatMessagesResponse.messages:type_name -> chats.Message
	8, // 1: chats.Message.created_at:type_name -> google.protobuf.Timestamp
	8, // 2: chats.Message.updated_at:type_name -> google.protobuf.Timestamp
	8, // 3: chats.Message.deleted_at:type_name -> google.protobuf.Timestamp
	0, // 4: chats.Chats.GetUserChats:input_type -> chats.GetUserChatsRequest
	2, // 5: chats.Chats.GetChatMessages:input_type -> chats.GetChatMessagesRequest
	4, // 6: chats.Chats.StoreMessage:input_type -> chats.StoreMessageRequest
	5, // 7: chats.Chats.FindOrCreateChat:input_type -> chats.FindOrCreateChatRequest
	1, // 8: chats.Chats.GetUserChats:output_type -> chats.GetUserChatsResponse
	3, // 9: chats.Chats.GetChatMessages:output_type -> chats.GetChatMessagesResponse
	9, // 10: chats.Chats.StoreMessage:output_type -> google.protobuf.Empty
	6, // 11: chats.Chats.FindOrCreateChat:output_type -> chats.FindOrCreateChatResponse
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_services_chats_grpc_api_proto_init() }
func file_services_chats_grpc_api_proto_init() {
	if File_services_chats_grpc_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_chats_grpc_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserChatsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_chats_grpc_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserChatsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_chats_grpc_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChatMessagesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_chats_grpc_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChatMessagesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_chats_grpc_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreMessageRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_chats_grpc_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindOrCreateChatRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_chats_grpc_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindOrCreateChatResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_services_chats_grpc_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_chats_grpc_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_chats_grpc_api_proto_goTypes,
		DependencyIndexes: file_services_chats_grpc_api_proto_depIdxs,
		MessageInfos:      file_services_chats_grpc_api_proto_msgTypes,
	}.Build()
	File_services_chats_grpc_api_proto = out.File
	file_services_chats_grpc_api_proto_rawDesc = nil
	file_services_chats_grpc_api_proto_goTypes = nil
	file_services_chats_grpc_api_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ChatsClient is the client API for Chats service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatsClient interface {
	GetUserChats(ctx context.Context, in *GetUserChatsRequest, opts ...grpc.CallOption) (*GetUserChatsResponse, error)
	GetChatMessages(ctx context.Context, in *GetChatMessagesRequest, opts ...grpc.CallOption) (*GetChatMessagesResponse, error)
	StoreMessage(ctx context.Context, in *StoreMessageRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	FindOrCreateChat(ctx context.Context, in *FindOrCreateChatRequest, opts ...grpc.CallOption) (*FindOrCreateChatResponse, error)
}

type chatsClient struct {
	cc grpc.ClientConnInterface
}

func NewChatsClient(cc grpc.ClientConnInterface) ChatsClient {
	return &chatsClient{cc}
}

func (c *chatsClient) GetUserChats(ctx context.Context, in *GetUserChatsRequest, opts ...grpc.CallOption) (*GetUserChatsResponse, error) {
	out := new(GetUserChatsResponse)
	err := c.cc.Invoke(ctx, "/chats.Chats/GetUserChats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatsClient) GetChatMessages(ctx context.Context, in *GetChatMessagesRequest, opts ...grpc.CallOption) (*GetChatMessagesResponse, error) {
	out := new(GetChatMessagesResponse)
	err := c.cc.Invoke(ctx, "/chats.Chats/GetChatMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatsClient) StoreMessage(ctx context.Context, in *StoreMessageRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/chats.Chats/StoreMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatsClient) FindOrCreateChat(ctx context.Context, in *FindOrCreateChatRequest, opts ...grpc.CallOption) (*FindOrCreateChatResponse, error) {
	out := new(FindOrCreateChatResponse)
	err := c.cc.Invoke(ctx, "/chats.Chats/FindOrCreateChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatsServer is the server API for Chats service.
type ChatsServer interface {
	GetUserChats(context.Context, *GetUserChatsRequest) (*GetUserChatsResponse, error)
	GetChatMessages(context.Context, *GetChatMessagesRequest) (*GetChatMessagesResponse, error)
	StoreMessage(context.Context, *StoreMessageRequest) (*empty.Empty, error)
	FindOrCreateChat(context.Context, *FindOrCreateChatRequest) (*FindOrCreateChatResponse, error)
}

// UnimplementedChatsServer can be embedded to have forward compatible implementations.
type UnimplementedChatsServer struct {
}

func (*UnimplementedChatsServer) GetUserChats(context.Context, *GetUserChatsRequest) (*GetUserChatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserChats not implemented")
}
func (*UnimplementedChatsServer) GetChatMessages(context.Context, *GetChatMessagesRequest) (*GetChatMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatMessages not implemented")
}
func (*UnimplementedChatsServer) StoreMessage(context.Context, *StoreMessageRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreMessage not implemented")
}
func (*UnimplementedChatsServer) FindOrCreateChat(context.Context, *FindOrCreateChatRequest) (*FindOrCreateChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindOrCreateChat not implemented")
}

func RegisterChatsServer(s *grpc.Server, srv ChatsServer) {
	s.RegisterService(&_Chats_serviceDesc, srv)
}

func _Chats_GetUserChats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserChatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatsServer).GetUserChats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chats.Chats/GetUserChats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatsServer).GetUserChats(ctx, req.(*GetUserChatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chats_GetChatMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatsServer).GetChatMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chats.Chats/GetChatMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatsServer).GetChatMessages(ctx, req.(*GetChatMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chats_StoreMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatsServer).StoreMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chats.Chats/StoreMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatsServer).StoreMessage(ctx, req.(*StoreMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chats_FindOrCreateChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindOrCreateChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatsServer).FindOrCreateChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chats.Chats/FindOrCreateChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatsServer).FindOrCreateChat(ctx, req.(*FindOrCreateChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Chats_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chats.Chats",
	HandlerType: (*ChatsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserChats",
			Handler:    _Chats_GetUserChats_Handler,
		},
		{
			MethodName: "GetChatMessages",
			Handler:    _Chats_GetChatMessages_Handler,
		},
		{
			MethodName: "StoreMessage",
			Handler:    _Chats_StoreMessage_Handler,
		},
		{
			MethodName: "FindOrCreateChat",
			Handler:    _Chats_FindOrCreateChat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/chats/grpc/api.proto",
}
