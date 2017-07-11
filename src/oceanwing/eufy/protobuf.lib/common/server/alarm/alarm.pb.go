// Code generated by protoc-gen-go.
// source: common/server/alarm/alarm.proto
// DO NOT EDIT!

/*
Package alarm_message is a generated protocol buffer package.

It is generated from these files:
	common/server/alarm/alarm.proto

It has these top-level messages:
	Alarm
	SyncTime
*/
package alarm

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

type Alarm struct {
	Hours            *uint32 `protobuf:"varint,1,req,name=hours" json:"hours,omitempty"`
	Minutes          *uint32 `protobuf:"varint,2,req,name=minutes" json:"minutes,omitempty"`
	Repetiton        *bool   `protobuf:"varint,3,req,name=repetiton" json:"repetiton,omitempty"`
	WeekInfo         *uint32 `protobuf:"varint,4,req,name=week_info,json=weekInfo" json:"week_info,omitempty"`
	Seconds          *uint32 `protobuf:"varint,5,opt,name=seconds" json:"seconds,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Alarm) Reset()                    { *m = Alarm{} }
func (m *Alarm) String() string            { return proto.CompactTextString(m) }
func (*Alarm) ProtoMessage()               {}
func (*Alarm) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Alarm) GetHours() uint32 {
	if m != nil && m.Hours != nil {
		return *m.Hours
	}
	return 0
}

func (m *Alarm) GetMinutes() uint32 {
	if m != nil && m.Minutes != nil {
		return *m.Minutes
	}
	return 0
}

func (m *Alarm) GetRepetiton() bool {
	if m != nil && m.Repetiton != nil {
		return *m.Repetiton
	}
	return false
}

func (m *Alarm) GetWeekInfo() uint32 {
	if m != nil && m.WeekInfo != nil {
		return *m.WeekInfo
	}
	return 0
}

func (m *Alarm) GetSeconds() uint32 {
	if m != nil && m.Seconds != nil {
		return *m.Seconds
	}
	return 0
}

type SyncTime struct {
	Year             *uint32 `protobuf:"varint,1,req,name=year" json:"year,omitempty"`
	Month            *uint32 `protobuf:"varint,2,req,name=month" json:"month,omitempty"`
	Day              *uint32 `protobuf:"varint,3,req,name=day" json:"day,omitempty"`
	Weekday          *uint32 `protobuf:"varint,4,req,name=weekday" json:"weekday,omitempty"`
	Hours            *uint32 `protobuf:"varint,5,req,name=hours" json:"hours,omitempty"`
	Minutes          *uint32 `protobuf:"varint,6,req,name=minutes" json:"minutes,omitempty"`
	Seconds          *uint32 `protobuf:"varint,7,req,name=seconds" json:"seconds,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SyncTime) Reset()                    { *m = SyncTime{} }
func (m *SyncTime) String() string            { return proto.CompactTextString(m) }
func (*SyncTime) ProtoMessage()               {}
func (*SyncTime) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SyncTime) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *SyncTime) GetMonth() uint32 {
	if m != nil && m.Month != nil {
		return *m.Month
	}
	return 0
}

func (m *SyncTime) GetDay() uint32 {
	if m != nil && m.Day != nil {
		return *m.Day
	}
	return 0
}

func (m *SyncTime) GetWeekday() uint32 {
	if m != nil && m.Weekday != nil {
		return *m.Weekday
	}
	return 0
}

func (m *SyncTime) GetHours() uint32 {
	if m != nil && m.Hours != nil {
		return *m.Hours
	}
	return 0
}

func (m *SyncTime) GetMinutes() uint32 {
	if m != nil && m.Minutes != nil {
		return *m.Minutes
	}
	return 0
}

func (m *SyncTime) GetSeconds() uint32 {
	if m != nil && m.Seconds != nil {
		return *m.Seconds
	}
	return 0
}

func init() {
	proto.RegisterType((*Alarm)(nil), "alarm.message.Alarm")
	proto.RegisterType((*SyncTime)(nil), "alarm.message.SyncTime")
}

func init() { proto.RegisterFile("common/server/alarm/alarm.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0xd9, 0xb4, 0x6b, 0xb7, 0x03, 0x0b, 0x12, 0x3c, 0x04, 0x14, 0x5c, 0x7a, 0xea, 0xc9,
	0x3e, 0x83, 0x9e, 0xf4, 0xba, 0x7a, 0x97, 0xb0, 0x9d, 0xda, 0x45, 0x33, 0x53, 0x92, 0x54, 0xd9,
	0x77, 0xf0, 0x49, 0x7c, 0x4a, 0x99, 0x4d, 0xca, 0xea, 0xa1, 0x97, 0x30, 0xdf, 0x1f, 0x66, 0xfe,
	0xff, 0x87, 0xdb, 0x8e, 0x9d, 0x63, 0xda, 0x04, 0xf4, 0x9f, 0xe8, 0x37, 0xf6, 0xc3, 0x7a, 0x97,
	0xde, 0xbb, 0x83, 0xe7, 0xc8, 0xba, 0x4e, 0xe0, 0x30, 0x04, 0xfb, 0x86, 0xab, 0xef, 0x02, 0xca,
	0x7b, 0x51, 0xf4, 0x15, 0x94, 0x7b, 0x3e, 0xfa, 0x60, 0x8a, 0x46, 0xad, 0xeb, 0x36, 0x81, 0x36,
	0xb0, 0x70, 0x3d, 0x1d, 0x23, 0x06, 0xa3, 0x46, 0xfd, 0x84, 0xfa, 0x06, 0x96, 0x1e, 0x0f, 0x18,
	0xfb, 0xc8, 0x64, 0x66, 0x8d, 0x5a, 0x57, 0xed, 0x24, 0xe8, 0x6b, 0x58, 0x7e, 0x21, 0xbe, 0xbf,
	0xf6, 0xb4, 0x63, 0x33, 0x1f, 0x37, 0x2b, 0x11, 0x9e, 0x68, 0xc7, 0x72, 0x34, 0x60, 0xc7, 0xb4,
	0x0d, 0xa6, 0x6c, 0x0a, 0x39, 0x9a, 0x71, 0xf5, 0x53, 0x40, 0xf5, 0x3c, 0x50, 0xf7, 0xd2, 0x3b,
	0xd4, 0x1a, 0xe6, 0x03, 0x5a, 0x9f, 0x03, 0x8d, 0xb3, 0xa4, 0x74, 0x4c, 0x71, 0x9f, 0xd3, 0x24,
	0xd0, 0x97, 0x30, 0xdb, 0xda, 0x61, 0x4c, 0x51, 0xb7, 0x32, 0x8a, 0x85, 0xd8, 0x89, 0x9a, 0xdc,
	0x4f, 0x38, 0xf5, 0x2c, 0xcf, 0xf4, 0xbc, 0xf8, 0xdf, 0xf3, 0x4f, 0xd8, 0x45, 0xfa, 0xc9, 0xf8,
	0xa0, 0x1e, 0xd5, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x84, 0x30, 0x63, 0x87, 0x70, 0x01, 0x00,
	0x00,
}
