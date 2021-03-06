// Code generated by protoc-gen-go.
// source: common.proto
// DO NOT EDIT!

/*
Package eufy is a generated protocol buffer package.

It is generated from these files:
	common.proto

It has these top-level messages:
	Server2DevMessage
	DevInfo
	Dev2ServerMessage
*/
package commonpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Server2Dev_CmdType int32

const (
	Server2Dev_CmdType_RESTORE_FACTORY_SETTING Server2Dev_CmdType = 1
	Server2Dev_CmdType_RESPONSE_LOCAL_CODE     Server2Dev_CmdType = 2
)

var Server2Dev_CmdType_name = map[int32]string{
	1: "RESTORE_FACTORY_SETTING",
	2: "RESPONSE_LOCAL_CODE",
}
var Server2Dev_CmdType_value = map[string]int32{
	"RESTORE_FACTORY_SETTING": 1,
	"RESPONSE_LOCAL_CODE":     2,
}

func (x Server2Dev_CmdType) Enum() *Server2Dev_CmdType {
	p := new(Server2Dev_CmdType)
	*p = x
	return p
}
func (x Server2Dev_CmdType) String() string {
	return proto.EnumName(Server2Dev_CmdType_name, int32(x))
}
func (x *Server2Dev_CmdType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Server2Dev_CmdType_value, data, "Server2Dev_CmdType")
	if err != nil {
		return err
	}
	*x = Server2Dev_CmdType(value)
	return nil
}
func (Server2Dev_CmdType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Dev2Server_CmdType int32

const (
	Dev2Server_CmdType_CLEAR_ALL_CONNECT  Dev2Server_CmdType = 1
	Dev2Server_CmdType_REQUEST_LOCAL_CODE Dev2Server_CmdType = 2
	Dev2Server_CmdType_REPORT_DEV_INFO    Dev2Server_CmdType = 3
)

var Dev2Server_CmdType_name = map[int32]string{
	1: "CLEAR_ALL_CONNECT",
	2: "REQUEST_LOCAL_CODE",
	3: "REPORT_DEV_INFO",
}
var Dev2Server_CmdType_value = map[string]int32{
	"CLEAR_ALL_CONNECT":  1,
	"REQUEST_LOCAL_CODE": 2,
	"REPORT_DEV_INFO":    3,
}

func (x Dev2Server_CmdType) Enum() *Dev2Server_CmdType {
	p := new(Dev2Server_CmdType)
	*p = x
	return p
}
func (x Dev2Server_CmdType) String() string {
	return proto.EnumName(Dev2Server_CmdType_name, int32(x))
}
func (x *Dev2Server_CmdType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Dev2Server_CmdType_value, data, "Dev2Server_CmdType")
	if err != nil {
		return err
	}
	*x = Dev2Server_CmdType(value)
	return nil
}
func (Dev2Server_CmdType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Server2DevMessage struct {
	Type             *Server2Dev_CmdType `protobuf:"varint,1,req,name=Type,enum=eufy.Server2Dev_CmdType" json:"Type,omitempty"`
	LocalCode        *string             `protobuf:"bytes,2,opt,name=local_code,json=localCode" json:"local_code,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *Server2DevMessage) Reset()                    { *m = Server2DevMessage{} }
func (m *Server2DevMessage) String() string            { return proto.CompactTextString(m) }
func (*Server2DevMessage) ProtoMessage()               {}
func (*Server2DevMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Server2DevMessage) GetType() Server2Dev_CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Server2Dev_CmdType_RESTORE_FACTORY_SETTING
}

func (m *Server2DevMessage) GetLocalCode() string {
	if m != nil && m.LocalCode != nil {
		return *m.LocalCode
	}
	return ""
}

type DevInfo struct {
	MacAddr          *string `protobuf:"bytes,1,req,name=mac_addr,json=macAddr" json:"mac_addr,omitempty"`
	FwVersion        *uint32 `protobuf:"varint,2,req,name=fw_version,json=fwVersion" json:"fw_version,omitempty"`
	LanIpaddr        *uint32 `protobuf:"varint,3,req,name=lan_ipaddr,json=lanIpaddr" json:"lan_ipaddr,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DevInfo) Reset()                    { *m = DevInfo{} }
func (m *DevInfo) String() string            { return proto.CompactTextString(m) }
func (*DevInfo) ProtoMessage()               {}
func (*DevInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DevInfo) GetMacAddr() string {
	if m != nil && m.MacAddr != nil {
		return *m.MacAddr
	}
	return ""
}

func (m *DevInfo) GetFwVersion() uint32 {
	if m != nil && m.FwVersion != nil {
		return *m.FwVersion
	}
	return 0
}

func (m *DevInfo) GetLanIpaddr() uint32 {
	if m != nil && m.LanIpaddr != nil {
		return *m.LanIpaddr
	}
	return 0
}

type Dev2ServerMessage struct {
	Type             *Dev2Server_CmdType `protobuf:"varint,1,req,name=Type,enum=eufy.Dev2Server_CmdType" json:"Type,omitempty"`
	Info             *DevInfo            `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *Dev2ServerMessage) Reset()                    { *m = Dev2ServerMessage{} }
func (m *Dev2ServerMessage) String() string            { return proto.CompactTextString(m) }
func (*Dev2ServerMessage) ProtoMessage()               {}
func (*Dev2ServerMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Dev2ServerMessage) GetType() Dev2Server_CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Dev2Server_CmdType_CLEAR_ALL_CONNECT
}

func (m *Dev2ServerMessage) GetInfo() *DevInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func init() {
	proto.RegisterType((*Server2DevMessage)(nil), "eufy.Server2DevMessage")
	proto.RegisterType((*DevInfo)(nil), "eufy.DevInfo")
	proto.RegisterType((*Dev2ServerMessage)(nil), "eufy.Dev2ServerMessage")
	proto.RegisterEnum("eufy.Server2Dev_CmdType", Server2Dev_CmdType_name, Server2Dev_CmdType_value)
	proto.RegisterEnum("eufy.Dev2Server_CmdType", Dev2Server_CmdType_name, Dev2Server_CmdType_value)
}

func init() { proto.RegisterFile("common.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 362 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x51, 0xab, 0xd3, 0x30,
	0x1c, 0xc5, 0x69, 0xef, 0xe0, 0xda, 0xbf, 0xf7, 0x6a, 0x97, 0x8b, 0xae, 0x22, 0xc2, 0xec, 0xd3,
	0x18, 0xd2, 0x87, 0x7e, 0x83, 0xae, 0xcd, 0xa4, 0x52, 0xdb, 0x99, 0xc6, 0xa1, 0x4f, 0x31, 0x6b,
	0xd2, 0x51, 0x58, 0x93, 0xd1, 0xb9, 0x8e, 0x7d, 0x7b, 0x69, 0x06, 0x0e, 0x1d, 0xdc, 0xd7, 0x73,
	0x7e, 0x39, 0x39, 0x27, 0x81, 0x87, 0x4a, 0xb7, 0xad, 0x56, 0xc1, 0xbe, 0xd3, 0xbf, 0x35, 0x1a,
	0xc9, 0x63, 0x7d, 0xf6, 0x7f, 0xc1, 0xb8, 0x94, 0x5d, 0x2f, 0xbb, 0x30, 0x91, 0xfd, 0x57, 0x79,
	0x38, 0xf0, 0xad, 0x44, 0x9f, 0x60, 0x44, 0xcf, 0x7b, 0xe9, 0x59, 0x53, 0x7b, 0xf6, 0x2a, 0xf4,
	0x82, 0x81, 0x0c, 0xae, 0x18, 0x8b, 0x5b, 0x31, 0xf8, 0xc4, 0x50, 0xe8, 0x03, 0xc0, 0x4e, 0x57,
	0x7c, 0xc7, 0x2a, 0x2d, 0xa4, 0x67, 0x4f, 0xad, 0x99, 0x43, 0x1c, 0xa3, 0xc4, 0x5a, 0x48, 0x7f,
	0x03, 0xf7, 0x89, 0xec, 0x53, 0x55, 0x6b, 0xf4, 0x0e, 0x5e, 0xb4, 0xbc, 0x62, 0x5c, 0x88, 0xce,
	0x64, 0x3b, 0xe4, 0xbe, 0xe5, 0x55, 0x24, 0x44, 0x37, 0x84, 0xd4, 0x27, 0xd6, 0xcb, 0xee, 0xd0,
	0x68, 0xe5, 0xd9, 0x53, 0x7b, 0xf6, 0x48, 0x9c, 0xfa, 0xb4, 0xbe, 0x08, 0xe6, 0x0e, 0xae, 0x58,
	0xb3, 0x37, 0x67, 0xef, 0x2e, 0xf6, 0x8e, 0xab, 0xd4, 0x08, 0xbe, 0x80, 0x71, 0x22, 0xfb, 0xf0,
	0x52, 0xf1, 0xd9, 0x15, 0x57, 0xec, 0xbf, 0x15, 0x1f, 0x61, 0xd4, 0xa8, 0x5a, 0x9b, 0xfe, 0x2f,
	0xc3, 0xc7, 0xbf, 0xf4, 0x50, 0x9c, 0x18, 0x6b, 0xfe, 0x05, 0xd0, 0xed, 0x23, 0xa0, 0xf7, 0x30,
	0x21, 0xb8, 0xa4, 0x05, 0xc1, 0x6c, 0x19, 0xc5, 0xb4, 0x20, 0x3f, 0x59, 0x89, 0x29, 0x4d, 0xf3,
	0xcf, 0xae, 0x85, 0x26, 0xf0, 0x44, 0x70, 0xb9, 0x2a, 0xf2, 0x12, 0xb3, 0xac, 0x88, 0xa3, 0x8c,
	0xc5, 0x45, 0x82, 0x5d, 0x7b, 0xfe, 0x03, 0xd0, 0x6d, 0x15, 0xf4, 0x06, 0xc6, 0x71, 0x86, 0x23,
	0xc2, 0xa2, 0x6c, 0x20, 0xf3, 0x1c, 0xc7, 0xd4, 0xb5, 0xd0, 0x5b, 0x40, 0x04, 0x7f, 0xfb, 0x8e,
	0x4b, 0xfa, 0x4f, 0x08, 0x7a, 0x82, 0xd7, 0x04, 0xaf, 0x0a, 0x42, 0x59, 0x82, 0xd7, 0x2c, 0xcd,
	0x97, 0x85, 0x7b, 0xb7, 0xf0, 0x01, 0x55, 0xba, 0x0d, 0x74, 0x25, 0xb9, 0x3a, 0x35, 0x6a, 0x6b,
	0x96, 0x2c, 0x1e, 0xf0, 0x71, 0x79, 0x5e, 0x0d, 0x1f, 0xbf, 0x39, 0xd6, 0x7f, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x6f, 0x62, 0xab, 0x7c, 0x0a, 0x02, 0x00, 0x00,
}
