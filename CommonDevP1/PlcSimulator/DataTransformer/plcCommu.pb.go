// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.15.8
// source: plcCommu.proto

package DataTransformer

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StRealDt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RealForce float32 `protobuf:"fixed32,1,opt,name=RealForce,proto3" json:"RealForce,omitempty"`
}

func (x *StRealDt) Reset() {
	*x = StRealDt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plcCommu_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StRealDt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StRealDt) ProtoMessage() {}

func (x *StRealDt) ProtoReflect() protoreflect.Message {
	mi := &file_plcCommu_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StRealDt.ProtoReflect.Descriptor instead.
func (*StRealDt) Descriptor() ([]byte, []int) {
	return file_plcCommu_proto_rawDescGZIP(), []int{0}
}

func (x *StRealDt) GetRealForce() float32 {
	if x != nil {
		return x.RealForce
	}
	return 0
}

type StCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CmdCode uint32  `protobuf:"varint,1,opt,name=CmdCode,proto3" json:"CmdCode,omitempty"`
	Param1  int32   `protobuf:"varint,2,opt,name=Param1,proto3" json:"Param1,omitempty"`
	Param2  int32   `protobuf:"varint,3,opt,name=Param2,proto3" json:"Param2,omitempty"`
	Param3  int32   `protobuf:"varint,4,opt,name=Param3,proto3" json:"Param3,omitempty"`
	Param4  int32   `protobuf:"varint,5,opt,name=Param4,proto3" json:"Param4,omitempty"`
	Param5  float32 `protobuf:"fixed32,6,opt,name=Param5,proto3" json:"Param5,omitempty"`
	Param6  float32 `protobuf:"fixed32,7,opt,name=Param6,proto3" json:"Param6,omitempty"`
	Param7  float32 `protobuf:"fixed32,8,opt,name=Param7,proto3" json:"Param7,omitempty"`
	Param8  float32 `protobuf:"fixed32,9,opt,name=Param8,proto3" json:"Param8,omitempty"`
}

func (x *StCmd) Reset() {
	*x = StCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plcCommu_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StCmd) ProtoMessage() {}

func (x *StCmd) ProtoReflect() protoreflect.Message {
	mi := &file_plcCommu_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StCmd.ProtoReflect.Descriptor instead.
func (*StCmd) Descriptor() ([]byte, []int) {
	return file_plcCommu_proto_rawDescGZIP(), []int{1}
}

func (x *StCmd) GetCmdCode() uint32 {
	if x != nil {
		return x.CmdCode
	}
	return 0
}

func (x *StCmd) GetParam1() int32 {
	if x != nil {
		return x.Param1
	}
	return 0
}

func (x *StCmd) GetParam2() int32 {
	if x != nil {
		return x.Param2
	}
	return 0
}

func (x *StCmd) GetParam3() int32 {
	if x != nil {
		return x.Param3
	}
	return 0
}

func (x *StCmd) GetParam4() int32 {
	if x != nil {
		return x.Param4
	}
	return 0
}

func (x *StCmd) GetParam5() float32 {
	if x != nil {
		return x.Param5
	}
	return 0
}

func (x *StCmd) GetParam6() float32 {
	if x != nil {
		return x.Param6
	}
	return 0
}

func (x *StCmd) GetParam7() float32 {
	if x != nil {
		return x.Param7
	}
	return 0
}

func (x *StCmd) GetParam8() float32 {
	if x != nil {
		return x.Param8
	}
	return 0
}

type MsgCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CmdCnt uint32   `protobuf:"varint,1,opt,name=CmdCnt,proto3" json:"CmdCnt,omitempty"`
	Cmds   []*StCmd `protobuf:"bytes,2,rep,name=Cmds,proto3" json:"Cmds,omitempty"`
	CmdIdx uint32   `protobuf:"varint,3,opt,name=CmdIdx,proto3" json:"CmdIdx,omitempty"`
}

func (x *MsgCmd) Reset() {
	*x = MsgCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plcCommu_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgCmd) ProtoMessage() {}

func (x *MsgCmd) ProtoReflect() protoreflect.Message {
	mi := &file_plcCommu_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgCmd.ProtoReflect.Descriptor instead.
func (*MsgCmd) Descriptor() ([]byte, []int) {
	return file_plcCommu_proto_rawDescGZIP(), []int{2}
}

func (x *MsgCmd) GetCmdCnt() uint32 {
	if x != nil {
		return x.CmdCnt
	}
	return 0
}

func (x *MsgCmd) GetCmds() []*StCmd {
	if x != nil {
		return x.Cmds
	}
	return nil
}

func (x *MsgCmd) GetCmdIdx() uint32 {
	if x != nil {
		return x.CmdIdx
	}
	return 0
}

type MsgRealDt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CmdIdx  uint32      `protobuf:"varint,1,opt,name=CmdIdx,proto3" json:"CmdIdx,omitempty"`
	RealDts []*StRealDt `protobuf:"bytes,2,rep,name=RealDts,proto3" json:"RealDts,omitempty"`
}

func (x *MsgRealDt) Reset() {
	*x = MsgRealDt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plcCommu_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgRealDt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgRealDt) ProtoMessage() {}

func (x *MsgRealDt) ProtoReflect() protoreflect.Message {
	mi := &file_plcCommu_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgRealDt.ProtoReflect.Descriptor instead.
func (*MsgRealDt) Descriptor() ([]byte, []int) {
	return file_plcCommu_proto_rawDescGZIP(), []int{3}
}

func (x *MsgRealDt) GetCmdIdx() uint32 {
	if x != nil {
		return x.CmdIdx
	}
	return 0
}

func (x *MsgRealDt) GetRealDts() []*StRealDt {
	if x != nil {
		return x.RealDts
	}
	return nil
}

type MsgRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reserve uint32 `protobuf:"varint,1,opt,name=Reserve,proto3" json:"Reserve,omitempty"`
}

func (x *MsgRequest) Reset() {
	*x = MsgRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plcCommu_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgRequest) ProtoMessage() {}

func (x *MsgRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plcCommu_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgRequest.ProtoReflect.Descriptor instead.
func (*MsgRequest) Descriptor() ([]byte, []int) {
	return file_plcCommu_proto_rawDescGZIP(), []int{4}
}

func (x *MsgRequest) GetReserve() uint32 {
	if x != nil {
		return x.Reserve
	}
	return 0
}

var File_plcCommu_proto protoreflect.FileDescriptor

var file_plcCommu_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x6c, 0x63, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x28, 0x0a, 0x08, 0x53, 0x74, 0x52, 0x65, 0x61, 0x6c, 0x44, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x52, 0x65, 0x61, 0x6c, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x09, 0x52, 0x65, 0x61, 0x6c, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x22, 0xe1, 0x01, 0x0a, 0x05, 0x53,
	0x74, 0x43, 0x6d, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6d, 0x64, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x43, 0x6d, 0x64, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x32,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x32, 0x12, 0x16,
	0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x33, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x33, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x34,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x34, 0x12, 0x16,
	0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x35, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x35, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x36,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x36, 0x12, 0x16,
	0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x37, 0x18, 0x08, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x37, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x38,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x38, 0x22, 0x54,
	0x0a, 0x06, 0x4d, 0x73, 0x67, 0x43, 0x6d, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x6d, 0x64, 0x43,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x43, 0x6d, 0x64, 0x43, 0x6e, 0x74,
	0x12, 0x1a, 0x0a, 0x04, 0x43, 0x6d, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06,
	0x2e, 0x53, 0x74, 0x43, 0x6d, 0x64, 0x52, 0x04, 0x43, 0x6d, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x43, 0x6d, 0x64, 0x49, 0x64, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x43, 0x6d,
	0x64, 0x49, 0x64, 0x78, 0x22, 0x48, 0x0a, 0x09, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x61, 0x6c, 0x44,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x6d, 0x64, 0x49, 0x64, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x43, 0x6d, 0x64, 0x49, 0x64, 0x78, 0x12, 0x23, 0x0a, 0x07, 0x52, 0x65, 0x61,
	0x6c, 0x44, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x53, 0x74, 0x52,
	0x65, 0x61, 0x6c, 0x44, 0x74, 0x52, 0x07, 0x52, 0x65, 0x61, 0x6c, 0x44, 0x74, 0x73, 0x22, 0x26,
	0x0a, 0x0a, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x52,
	0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x42, 0x20, 0x5a, 0x1e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x44, 0x65, 0x76, 0x50, 0x31, 0x2f, 0x50, 0x6c, 0x63, 0x53, 0x69, 0x6d, 0x75, 0x6c, 0x61, 0x74,
	0x6f, 0x72, 0x2f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plcCommu_proto_rawDescOnce sync.Once
	file_plcCommu_proto_rawDescData = file_plcCommu_proto_rawDesc
)

func file_plcCommu_proto_rawDescGZIP() []byte {
	file_plcCommu_proto_rawDescOnce.Do(func() {
		file_plcCommu_proto_rawDescData = protoimpl.X.CompressGZIP(file_plcCommu_proto_rawDescData)
	})
	return file_plcCommu_proto_rawDescData
}

var file_plcCommu_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_plcCommu_proto_goTypes = []interface{}{
	(*StRealDt)(nil),   // 0: StRealDt
	(*StCmd)(nil),      // 1: StCmd
	(*MsgCmd)(nil),     // 2: MsgCmd
	(*MsgRealDt)(nil),  // 3: MsgRealDt
	(*MsgRequest)(nil), // 4: MsgRequest
}
var file_plcCommu_proto_depIdxs = []int32{
	1, // 0: MsgCmd.Cmds:type_name -> StCmd
	0, // 1: MsgRealDt.RealDts:type_name -> StRealDt
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_plcCommu_proto_init() }
func file_plcCommu_proto_init() {
	if File_plcCommu_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plcCommu_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StRealDt); i {
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
		file_plcCommu_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StCmd); i {
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
		file_plcCommu_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgCmd); i {
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
		file_plcCommu_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgRealDt); i {
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
		file_plcCommu_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgRequest); i {
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
			RawDescriptor: file_plcCommu_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_plcCommu_proto_goTypes,
		DependencyIndexes: file_plcCommu_proto_depIdxs,
		MessageInfos:      file_plcCommu_proto_msgTypes,
	}.Build()
	File_plcCommu_proto = out.File
	file_plcCommu_proto_rawDesc = nil
	file_plcCommu_proto_goTypes = nil
	file_plcCommu_proto_depIdxs = nil
}
