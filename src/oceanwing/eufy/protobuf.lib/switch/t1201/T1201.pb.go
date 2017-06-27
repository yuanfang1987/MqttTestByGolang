// Code generated by protoc-gen-go.
// source: Switch/T1201/T1201.proto
// DO NOT EDIT!

/*
Package T1201_message is a generated protocol buffer package.

It is generated from these files:
	Switch/T1201/T1201.proto

It has these top-level messages:
	DevStateMessage
	SyncAlarmRecordMessage
	ServerMessage
	HeartBeatMessage
	APPDataMessage
	DeviceMessage
*/
package t1201

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import alarm_message "oceanwing/eufy/protobuf.lib/common/server/alarm"
import away_mode_message "oceanwing/eufy/protobuf.lib/common/server/awaymode"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CmdType int32

const (
	// 设备向服务器请求时间及闹钟信息
	CmdType_DEV_REQUSET_TIME_Alarm      CmdType = 0
	CmdType_SERVER_RESPONSE_TIME_Alarm  CmdType = 1
	CmdType_DEV_RESPONSE_TIME_Alarm_ACK CmdType = 2
	// 设备向服务器请求时间信息
	CmdType_DEV_REQUEST_TIME      CmdType = 3
	CmdType_SERVER_RESPONSE_TIME  CmdType = 4
	CmdType_DEV_RESPONSE_TIME_ACK CmdType = 5
	// 设备主动上报设备状态，空闲时每20S主动上报一次，有状态变化立即上报
	CmdType_DEV_REPORT_STATUS CmdType = 6
	// 远程设置设备的开关，支持服务器及APP两种控制
	CmdType_REMOTE_SET_PLUG_STATE CmdType = 7
	// 设置设备在掉电后重新上电时的状态，仅支持服务器更改
	CmdType_DEV_REQUEST_POWERUP_PLUG_STATUS CmdType = 8
	CmdType_REMOTE_SET_POWERUP_PLUG_STATUS  CmdType = 9
	CmdType_DEV_RESPONE_POWERUP_PLUG_ACK    CmdType = 10
	//  上报电量 20min 上报一次
	CmdType_DEV_REPORT_ELECTRIC CmdType = 11
	// 设置离家模式状态
	CmdType_DEV_REQUEST_AWAYMODE_STATUS CmdType = 12
	CmdType_REMOTE_SET_AWAYMODE_STATUS  CmdType = 13
	CmdType_DEV_RESPONE_AWAYMODE_ACK    CmdType = 14
)

var CmdType_name = map[int32]string{
	0:  "DEV_REQUSET_TIME_Alarm",
	1:  "SERVER_RESPONSE_TIME_Alarm",
	2:  "DEV_RESPONSE_TIME_Alarm_ACK",
	3:  "DEV_REQUEST_TIME",
	4:  "SERVER_RESPONSE_TIME",
	5:  "DEV_RESPONSE_TIME_ACK",
	6:  "DEV_REPORT_STATUS",
	7:  "REMOTE_SET_PLUG_STATE",
	8:  "DEV_REQUEST_POWERUP_PLUG_STATUS",
	9:  "REMOTE_SET_POWERUP_PLUG_STATUS",
	10: "DEV_RESPONE_POWERUP_PLUG_ACK",
	11: "DEV_REPORT_ELECTRIC",
	12: "DEV_REQUEST_AWAYMODE_STATUS",
	13: "REMOTE_SET_AWAYMODE_STATUS",
	14: "DEV_RESPONE_AWAYMODE_ACK",
}
var CmdType_value = map[string]int32{
	"DEV_REQUSET_TIME_Alarm":          0,
	"SERVER_RESPONSE_TIME_Alarm":      1,
	"DEV_RESPONSE_TIME_Alarm_ACK":     2,
	"DEV_REQUEST_TIME":                3,
	"SERVER_RESPONSE_TIME":            4,
	"DEV_RESPONSE_TIME_ACK":           5,
	"DEV_REPORT_STATUS":               6,
	"REMOTE_SET_PLUG_STATE":           7,
	"DEV_REQUEST_POWERUP_PLUG_STATUS": 8,
	"REMOTE_SET_POWERUP_PLUG_STATUS":  9,
	"DEV_RESPONE_POWERUP_PLUG_ACK":    10,
	"DEV_REPORT_ELECTRIC":             11,
	"DEV_REQUEST_AWAYMODE_STATUS":     12,
	"REMOTE_SET_AWAYMODE_STATUS":      13,
	"DEV_RESPONE_AWAYMODE_ACK":        14,
}

func (x CmdType) Enum() *CmdType {
	p := new(CmdType)
	*p = x
	return p
}
func (x CmdType) String() string {
	return proto.EnumName(CmdType_name, int32(x))
}
func (x *CmdType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(CmdType_value, data, "CmdType")
	if err != nil {
		return err
	}
	*x = CmdType(value)
	return nil
}
func (CmdType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type DevStateMessage struct {
	Type             *CmdType `protobuf:"varint,1,req,name=type,enum=T1201.message.CmdType" json:"type,omitempty"`
	RelayState       *uint32  `protobuf:"varint,2,req,name=RelayState" json:"RelayState,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *DevStateMessage) Reset()                    { *m = DevStateMessage{} }
func (m *DevStateMessage) String() string            { return proto.CompactTextString(m) }
func (*DevStateMessage) ProtoMessage()               {}
func (*DevStateMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DevStateMessage) GetType() CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return CmdType_DEV_REQUSET_TIME_Alarm
}

func (m *DevStateMessage) GetRelayState() uint32 {
	if m != nil && m.RelayState != nil {
		return *m.RelayState
	}
	return 0
}

type SyncAlarmRecordMessage struct {
	AlarmRecordData  []*SyncAlarmRecordMessage_AlarmRecord `protobuf:"bytes,1,rep,name=AlarmRecordData" json:"AlarmRecordData,omitempty"`
	XXX_unrecognized []byte                                `json:"-"`
}

func (m *SyncAlarmRecordMessage) Reset()                    { *m = SyncAlarmRecordMessage{} }
func (m *SyncAlarmRecordMessage) String() string            { return proto.CompactTextString(m) }
func (*SyncAlarmRecordMessage) ProtoMessage()               {}
func (*SyncAlarmRecordMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SyncAlarmRecordMessage) GetAlarmRecordData() []*SyncAlarmRecordMessage_AlarmRecord {
	if m != nil {
		return m.AlarmRecordData
	}
	return nil
}

type SyncAlarmRecordMessage_AlarmRecord struct {
	AlarmMesage      *alarm_message.Alarm `protobuf:"bytes,1,req,name=AlarmMesage" json:"AlarmMesage,omitempty"`
	RelayState       *uint32              `protobuf:"varint,2,req,name=RelayState" json:"RelayState,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *SyncAlarmRecordMessage_AlarmRecord) Reset()         { *m = SyncAlarmRecordMessage_AlarmRecord{} }
func (m *SyncAlarmRecordMessage_AlarmRecord) String() string { return proto.CompactTextString(m) }
func (*SyncAlarmRecordMessage_AlarmRecord) ProtoMessage()    {}
func (*SyncAlarmRecordMessage_AlarmRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1, 0}
}

func (m *SyncAlarmRecordMessage_AlarmRecord) GetAlarmMesage() *alarm_message.Alarm {
	if m != nil {
		return m.AlarmMesage
	}
	return nil
}

func (m *SyncAlarmRecordMessage_AlarmRecord) GetRelayState() uint32 {
	if m != nil && m.RelayState != nil {
		return *m.RelayState
	}
	return 0
}

type ServerMessage struct {
	SessionId *int32 `protobuf:"varint,1,req,name=sessionId" json:"sessionId,omitempty"`
	// Types that are valid to be assigned to RemoteMessage:
	//	*ServerMessage_Sync_Time_Alarm_
	//	*ServerMessage_DevState
	RemoteMessage    isServerMessage_RemoteMessage `protobuf_oneof:"RemoteMessage"`
	XXX_unrecognized []byte                        `json:"-"`
}

func (m *ServerMessage) Reset()                    { *m = ServerMessage{} }
func (m *ServerMessage) String() string            { return proto.CompactTextString(m) }
func (*ServerMessage) ProtoMessage()               {}
func (*ServerMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isServerMessage_RemoteMessage interface {
	isServerMessage_RemoteMessage()
}

type ServerMessage_Sync_Time_Alarm_ struct {
	Sync_Time_Alarm *ServerMessage_Sync_Time_Alarm `protobuf:"bytes,2,opt,name=sync_Time_Alarm,json=syncTimeAlarm,oneof"`
}
type ServerMessage_DevState struct {
	DevState *DevStateMessage `protobuf:"bytes,3,opt,name=devState,oneof"`
}

func (*ServerMessage_Sync_Time_Alarm_) isServerMessage_RemoteMessage() {}
func (*ServerMessage_DevState) isServerMessage_RemoteMessage()         {}

func (m *ServerMessage) GetRemoteMessage() isServerMessage_RemoteMessage {
	if m != nil {
		return m.RemoteMessage
	}
	return nil
}

func (m *ServerMessage) GetSessionId() int32 {
	if m != nil && m.SessionId != nil {
		return *m.SessionId
	}
	return 0
}

func (m *ServerMessage) GetSync_Time_Alarm() *ServerMessage_Sync_Time_Alarm {
	if x, ok := m.GetRemoteMessage().(*ServerMessage_Sync_Time_Alarm_); ok {
		return x.Sync_Time_Alarm
	}
	return nil
}

func (m *ServerMessage) GetDevState() *DevStateMessage {
	if x, ok := m.GetRemoteMessage().(*ServerMessage_DevState); ok {
		return x.DevState
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ServerMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ServerMessage_OneofMarshaler, _ServerMessage_OneofUnmarshaler, _ServerMessage_OneofSizer, []interface{}{
		(*ServerMessage_Sync_Time_Alarm_)(nil),
		(*ServerMessage_DevState)(nil),
	}
}

func _ServerMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ServerMessage)
	// RemoteMessage
	switch x := m.RemoteMessage.(type) {
	case *ServerMessage_Sync_Time_Alarm_:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Sync_Time_Alarm); err != nil {
			return err
		}
	case *ServerMessage_DevState:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DevState); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ServerMessage.RemoteMessage has unexpected type %T", x)
	}
	return nil
}

func _ServerMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ServerMessage)
	switch tag {
	case 2: // RemoteMessage.sync_Time_Alarm
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ServerMessage_Sync_Time_Alarm)
		err := b.DecodeMessage(msg)
		m.RemoteMessage = &ServerMessage_Sync_Time_Alarm_{msg}
		return true, err
	case 3: // RemoteMessage.devState
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DevStateMessage)
		err := b.DecodeMessage(msg)
		m.RemoteMessage = &ServerMessage_DevState{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ServerMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ServerMessage)
	// RemoteMessage
	switch x := m.RemoteMessage.(type) {
	case *ServerMessage_Sync_Time_Alarm_:
		s := proto.Size(x.Sync_Time_Alarm)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ServerMessage_DevState:
		s := proto.Size(x.DevState)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ServerMessage_Sync_Time_Alarm struct {
	Type             *CmdType                            `protobuf:"varint,1,req,name=type,enum=T1201.message.CmdType" json:"type,omitempty"`
	Time             *alarm_message.SyncTime             `protobuf:"bytes,2,opt,name=time" json:"time,omitempty"`
	Alarm            *SyncAlarmRecordMessage             `protobuf:"bytes,3,opt,name=Alarm" json:"Alarm,omitempty"`
	SyncLeaveModeMsg *away_mode_message.LeaveHomeMessage `protobuf:"bytes,4,opt,name=syncLeaveModeMsg" json:"syncLeaveModeMsg,omitempty"`
	XXX_unrecognized []byte                              `json:"-"`
}

func (m *ServerMessage_Sync_Time_Alarm) Reset()         { *m = ServerMessage_Sync_Time_Alarm{} }
func (m *ServerMessage_Sync_Time_Alarm) String() string { return proto.CompactTextString(m) }
func (*ServerMessage_Sync_Time_Alarm) ProtoMessage()    {}
func (*ServerMessage_Sync_Time_Alarm) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2, 0}
}

func (m *ServerMessage_Sync_Time_Alarm) GetType() CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return CmdType_DEV_REQUSET_TIME_Alarm
}

func (m *ServerMessage_Sync_Time_Alarm) GetTime() *alarm_message.SyncTime {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *ServerMessage_Sync_Time_Alarm) GetAlarm() *SyncAlarmRecordMessage {
	if m != nil {
		return m.Alarm
	}
	return nil
}

func (m *ServerMessage_Sync_Time_Alarm) GetSyncLeaveModeMsg() *away_mode_message.LeaveHomeMessage {
	if m != nil {
		return m.SyncLeaveModeMsg
	}
	return nil
}

type HeartBeatMessage struct {
	Type             *CmdType `protobuf:"varint,1,req,name=type,enum=T1201.message.CmdType" json:"type,omitempty"`
	RelayState       *uint32  `protobuf:"varint,2,req,name=RelayState" json:"RelayState,omitempty"`
	Power            *uint32  `protobuf:"varint,3,req,name=power" json:"power,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *HeartBeatMessage) Reset()                    { *m = HeartBeatMessage{} }
func (m *HeartBeatMessage) String() string            { return proto.CompactTextString(m) }
func (*HeartBeatMessage) ProtoMessage()               {}
func (*HeartBeatMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *HeartBeatMessage) GetType() CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return CmdType_DEV_REQUSET_TIME_Alarm
}

func (m *HeartBeatMessage) GetRelayState() uint32 {
	if m != nil && m.RelayState != nil {
		return *m.RelayState
	}
	return 0
}

func (m *HeartBeatMessage) GetPower() uint32 {
	if m != nil && m.Power != nil {
		return *m.Power
	}
	return 0
}

type APPDataMessage struct {
	Type             *CmdType `protobuf:"varint,1,req,name=type,enum=T1201.message.CmdType" json:"type,omitempty"`
	RelayState       *uint32  `protobuf:"varint,2,req,name=RelayState" json:"RelayState,omitempty"`
	Power            *uint32  `protobuf:"varint,3,req,name=power" json:"power,omitempty"`
	Electric         *uint32  `protobuf:"varint,4,req,name=electric" json:"electric,omitempty"`
	WorkingTime      *uint32  `protobuf:"varint,5,req,name=workingTime" json:"workingTime,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *APPDataMessage) Reset()                    { *m = APPDataMessage{} }
func (m *APPDataMessage) String() string            { return proto.CompactTextString(m) }
func (*APPDataMessage) ProtoMessage()               {}
func (*APPDataMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *APPDataMessage) GetType() CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return CmdType_DEV_REQUSET_TIME_Alarm
}

func (m *APPDataMessage) GetRelayState() uint32 {
	if m != nil && m.RelayState != nil {
		return *m.RelayState
	}
	return 0
}

func (m *APPDataMessage) GetPower() uint32 {
	if m != nil && m.Power != nil {
		return *m.Power
	}
	return 0
}

func (m *APPDataMessage) GetElectric() uint32 {
	if m != nil && m.Electric != nil {
		return *m.Electric
	}
	return 0
}

func (m *APPDataMessage) GetWorkingTime() uint32 {
	if m != nil && m.WorkingTime != nil {
		return *m.WorkingTime
	}
	return 0
}

type DeviceMessage struct {
	SessionId *int32 `protobuf:"varint,1,req,name=sessionId" json:"sessionId,omitempty"`
	// Types that are valid to be assigned to DevMessage:
	//	*DeviceMessage_NonParaMsg
	//	*DeviceMessage_HeartBeat
	//	*DeviceMessage_ElectricMessage_
	//	*DeviceMessage_AppData
	DevMessage       isDeviceMessage_DevMessage `protobuf_oneof:"devMessage"`
	XXX_unrecognized []byte                     `json:"-"`
}

func (m *DeviceMessage) Reset()                    { *m = DeviceMessage{} }
func (m *DeviceMessage) String() string            { return proto.CompactTextString(m) }
func (*DeviceMessage) ProtoMessage()               {}
func (*DeviceMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type isDeviceMessage_DevMessage interface {
	isDeviceMessage_DevMessage()
}

type DeviceMessage_NonParaMsg struct {
	NonParaMsg *DeviceMessage_Non_ParamMsg `protobuf:"bytes,2,opt,name=NonParaMsg,oneof"`
}
type DeviceMessage_HeartBeat struct {
	HeartBeat *HeartBeatMessage `protobuf:"bytes,3,opt,name=heartBeat,oneof"`
}
type DeviceMessage_ElectricMessage_ struct {
	ElectricMessage *DeviceMessage_ElectricMessage `protobuf:"bytes,4,opt,name=electricMessage,oneof"`
}
type DeviceMessage_AppData struct {
	AppData *APPDataMessage `protobuf:"bytes,5,opt,name=appData,oneof"`
}

func (*DeviceMessage_NonParaMsg) isDeviceMessage_DevMessage()       {}
func (*DeviceMessage_HeartBeat) isDeviceMessage_DevMessage()        {}
func (*DeviceMessage_ElectricMessage_) isDeviceMessage_DevMessage() {}
func (*DeviceMessage_AppData) isDeviceMessage_DevMessage()          {}

func (m *DeviceMessage) GetDevMessage() isDeviceMessage_DevMessage {
	if m != nil {
		return m.DevMessage
	}
	return nil
}

func (m *DeviceMessage) GetSessionId() int32 {
	if m != nil && m.SessionId != nil {
		return *m.SessionId
	}
	return 0
}

func (m *DeviceMessage) GetNonParaMsg() *DeviceMessage_Non_ParamMsg {
	if x, ok := m.GetDevMessage().(*DeviceMessage_NonParaMsg); ok {
		return x.NonParaMsg
	}
	return nil
}

func (m *DeviceMessage) GetHeartBeat() *HeartBeatMessage {
	if x, ok := m.GetDevMessage().(*DeviceMessage_HeartBeat); ok {
		return x.HeartBeat
	}
	return nil
}

func (m *DeviceMessage) GetElectricMessage() *DeviceMessage_ElectricMessage {
	if x, ok := m.GetDevMessage().(*DeviceMessage_ElectricMessage_); ok {
		return x.ElectricMessage
	}
	return nil
}

func (m *DeviceMessage) GetAppData() *APPDataMessage {
	if x, ok := m.GetDevMessage().(*DeviceMessage_AppData); ok {
		return x.AppData
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*DeviceMessage) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _DeviceMessage_OneofMarshaler, _DeviceMessage_OneofUnmarshaler, _DeviceMessage_OneofSizer, []interface{}{
		(*DeviceMessage_NonParaMsg)(nil),
		(*DeviceMessage_HeartBeat)(nil),
		(*DeviceMessage_ElectricMessage_)(nil),
		(*DeviceMessage_AppData)(nil),
	}
}

func _DeviceMessage_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*DeviceMessage)
	// devMessage
	switch x := m.DevMessage.(type) {
	case *DeviceMessage_NonParaMsg:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.NonParaMsg); err != nil {
			return err
		}
	case *DeviceMessage_HeartBeat:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.HeartBeat); err != nil {
			return err
		}
	case *DeviceMessage_ElectricMessage_:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ElectricMessage); err != nil {
			return err
		}
	case *DeviceMessage_AppData:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AppData); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("DeviceMessage.DevMessage has unexpected type %T", x)
	}
	return nil
}

func _DeviceMessage_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*DeviceMessage)
	switch tag {
	case 2: // devMessage.NonParaMsg
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DeviceMessage_Non_ParamMsg)
		err := b.DecodeMessage(msg)
		m.DevMessage = &DeviceMessage_NonParaMsg{msg}
		return true, err
	case 3: // devMessage.heartBeat
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(HeartBeatMessage)
		err := b.DecodeMessage(msg)
		m.DevMessage = &DeviceMessage_HeartBeat{msg}
		return true, err
	case 4: // devMessage.electricMessage
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DeviceMessage_ElectricMessage)
		err := b.DecodeMessage(msg)
		m.DevMessage = &DeviceMessage_ElectricMessage_{msg}
		return true, err
	case 5: // devMessage.appData
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(APPDataMessage)
		err := b.DecodeMessage(msg)
		m.DevMessage = &DeviceMessage_AppData{msg}
		return true, err
	default:
		return false, nil
	}
}

func _DeviceMessage_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*DeviceMessage)
	// devMessage
	switch x := m.DevMessage.(type) {
	case *DeviceMessage_NonParaMsg:
		s := proto.Size(x.NonParaMsg)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *DeviceMessage_HeartBeat:
		s := proto.Size(x.HeartBeat)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *DeviceMessage_ElectricMessage_:
		s := proto.Size(x.ElectricMessage)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *DeviceMessage_AppData:
		s := proto.Size(x.AppData)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type DeviceMessage_ElectricMessage struct {
	Type             *CmdType `protobuf:"varint,1,req,name=type,enum=T1201.message.CmdType" json:"type,omitempty"`
	Electric         *uint32  `protobuf:"varint,2,req,name=electric" json:"electric,omitempty"`
	WorkingTime      *uint32  `protobuf:"varint,3,req,name=workingTime" json:"workingTime,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *DeviceMessage_ElectricMessage) Reset()         { *m = DeviceMessage_ElectricMessage{} }
func (m *DeviceMessage_ElectricMessage) String() string { return proto.CompactTextString(m) }
func (*DeviceMessage_ElectricMessage) ProtoMessage()    {}
func (*DeviceMessage_ElectricMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{5, 0}
}

func (m *DeviceMessage_ElectricMessage) GetType() CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return CmdType_DEV_REQUSET_TIME_Alarm
}

func (m *DeviceMessage_ElectricMessage) GetElectric() uint32 {
	if m != nil && m.Electric != nil {
		return *m.Electric
	}
	return 0
}

func (m *DeviceMessage_ElectricMessage) GetWorkingTime() uint32 {
	if m != nil && m.WorkingTime != nil {
		return *m.WorkingTime
	}
	return 0
}

type DeviceMessage_Non_ParamMsg struct {
	Type             *CmdType `protobuf:"varint,1,req,name=type,enum=T1201.message.CmdType" json:"type,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *DeviceMessage_Non_ParamMsg) Reset()                    { *m = DeviceMessage_Non_ParamMsg{} }
func (m *DeviceMessage_Non_ParamMsg) String() string            { return proto.CompactTextString(m) }
func (*DeviceMessage_Non_ParamMsg) ProtoMessage()               {}
func (*DeviceMessage_Non_ParamMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 1} }

func (m *DeviceMessage_Non_ParamMsg) GetType() CmdType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return CmdType_DEV_REQUSET_TIME_Alarm
}

func init() {
	proto.RegisterType((*DevStateMessage)(nil), "T1201.message.DevStateMessage")
	proto.RegisterType((*SyncAlarmRecordMessage)(nil), "T1201.message.SyncAlarmRecordMessage")
	proto.RegisterType((*SyncAlarmRecordMessage_AlarmRecord)(nil), "T1201.message.SyncAlarmRecordMessage.AlarmRecord")
	proto.RegisterType((*ServerMessage)(nil), "T1201.message.ServerMessage")
	proto.RegisterType((*ServerMessage_Sync_Time_Alarm)(nil), "T1201.message.ServerMessage.Sync_Time_Alarm")
	proto.RegisterType((*HeartBeatMessage)(nil), "T1201.message.HeartBeatMessage")
	proto.RegisterType((*APPDataMessage)(nil), "T1201.message.APPDataMessage")
	proto.RegisterType((*DeviceMessage)(nil), "T1201.message.DeviceMessage")
	proto.RegisterType((*DeviceMessage_ElectricMessage)(nil), "T1201.message.DeviceMessage.ElectricMessage")
	proto.RegisterType((*DeviceMessage_Non_ParamMsg)(nil), "T1201.message.DeviceMessage.Non_ParamMsg")
	proto.RegisterEnum("T1201.message.CmdType", CmdType_name, CmdType_value)
}

func init() { proto.RegisterFile("Switch/T1201/T1201.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 824 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x54, 0xdb, 0x6e, 0xe3, 0x44,
	0x18, 0x4e, 0x9c, 0x84, 0xb6, 0x7f, 0x9a, 0xc6, 0x0c, 0xdd, 0xae, 0x31, 0xa5, 0x89, 0xb2, 0x42,
	0x84, 0x65, 0x95, 0xd2, 0x5c, 0x20, 0x71, 0x90, 0x50, 0x0e, 0x23, 0x5c, 0xb5, 0x6e, 0xcc, 0xd8,
	0xe9, 0x82, 0x10, 0xb2, 0x2c, 0x67, 0x94, 0x8d, 0xa8, 0xe3, 0xc8, 0xb6, 0x12, 0x45, 0x3c, 0x08,
	0xaf, 0xc1, 0x53, 0x70, 0xcb, 0x0b, 0xf0, 0x20, 0x5c, 0x22, 0xcf, 0xd8, 0xf1, 0xc4, 0x89, 0x96,
	0xf6, 0x82, 0xbd, 0x89, 0x3a, 0xff, 0xff, 0x1d, 0xfe, 0x53, 0x0d, 0x8a, 0xb9, 0x9a, 0x45, 0xee,
	0x9b, 0x4b, 0xeb, 0xaa, 0xfb, 0xc5, 0x15, 0xff, 0xed, 0x2c, 0x02, 0x3f, 0xf2, 0x51, 0x8d, 0x3f,
	0x3c, 0x1a, 0x86, 0xce, 0x94, 0xaa, 0x0d, 0xd7, 0xf7, 0x3c, 0x7f, 0x7e, 0x19, 0xd2, 0x60, 0x49,
	0x83, 0x4b, 0xe7, 0xc1, 0x09, 0x3c, 0xfe, 0xcb, 0xf1, 0xea, 0xa7, 0x39, 0xc0, 0xca, 0x59, 0xdb,
	0x9e, 0x3f, 0xa1, 0xd9, 0x5f, 0x1c, 0xd8, 0xfa, 0x05, 0xea, 0x43, 0xba, 0x34, 0x23, 0x27, 0xa2,
	0x3a, 0x17, 0x47, 0x2f, 0xa1, 0x1c, 0xad, 0x17, 0x54, 0x29, 0x36, 0xa5, 0xf6, 0x49, 0xf7, 0xac,
	0xb3, 0x65, 0xdd, 0x19, 0x78, 0x13, 0x6b, 0xbd, 0xa0, 0x84, 0x61, 0xd0, 0x05, 0x00, 0xa1, 0x0f,
	0xce, 0x9a, 0x09, 0x28, 0x52, 0x53, 0x6a, 0xd7, 0x88, 0x10, 0x69, 0xfd, 0x5d, 0x84, 0x33, 0x73,
	0x3d, 0x77, 0x7b, 0x71, 0x6d, 0x84, 0xba, 0x7e, 0x30, 0x49, 0x6d, 0x7e, 0x86, 0xba, 0x10, 0x1d,
	0x3a, 0x91, 0xa3, 0x14, 0x9b, 0xa5, 0x76, 0xb5, 0x7b, 0x95, 0x73, 0xdc, 0xcf, 0xef, 0x08, 0x21,
	0x92, 0x57, 0x52, 0x29, 0x54, 0x85, 0x10, 0xfa, 0x32, 0x79, 0xea, 0x34, 0xa6, 0xb2, 0xce, 0xaa,
	0xdd, 0xd3, 0x0e, 0x9f, 0x98, 0xb7, 0x25, 0x28, 0x02, 0xff, 0xb3, 0xbd, 0xbf, 0x4a, 0x50, 0x33,
	0xd9, 0x88, 0xd3, 0xae, 0xce, 0xe1, 0x28, 0xa4, 0x61, 0x38, 0xf3, 0xe7, 0xd7, 0x13, 0xe6, 0x53,
	0x21, 0x59, 0x00, 0xdd, 0x43, 0x3d, 0x5c, 0xcf, 0x5d, 0xdb, 0x9a, 0x79, 0xd4, 0x66, 0x46, 0x8a,
	0xd4, 0x2c, 0xb6, 0xab, 0xdd, 0x57, 0xf9, 0x9e, 0x45, 0x51, 0x36, 0x01, 0x81, 0xa3, 0x15, 0x48,
	0x2d, 0x96, 0x89, 0x23, 0x2c, 0x80, 0xbe, 0x85, 0xc3, 0x49, 0xb2, 0x45, 0xa5, 0xc4, 0x04, 0x2f,
	0x72, 0x82, 0xb9, 0x25, 0x6b, 0x05, 0xb2, 0x61, 0xa8, 0xff, 0x14, 0xa1, 0x9e, 0xb3, 0x78, 0xd2,
	0x11, 0x7c, 0x0e, 0xe5, 0x68, 0xe6, 0xd1, 0xa4, 0x95, 0xe7, 0xb9, 0xb1, 0x9a, 0x49, 0xa5, 0x84,
	0x81, 0xd0, 0x37, 0x50, 0xe1, 0x8d, 0xf3, 0x3a, 0x3f, 0x79, 0xd4, 0xb2, 0x09, 0xe7, 0xa0, 0x11,
	0xc8, 0x71, 0xe3, 0xb7, 0xd4, 0x59, 0x52, 0xdd, 0x9f, 0x50, 0x3d, 0x9c, 0x2a, 0x65, 0xa6, 0xf3,
	0xa2, 0x93, 0x5d, 0x76, 0xaa, 0xc5, 0x60, 0x9a, 0xef, 0xa5, 0x4d, 0x93, 0x1d, 0x72, 0xbf, 0x0e,
	0x35, 0x42, 0x3d, 0x7f, 0x33, 0x97, 0x56, 0x04, 0xb2, 0x46, 0x9d, 0x20, 0xea, 0x53, 0x27, 0xfa,
	0x1f, 0xfe, 0x21, 0xd0, 0x29, 0x54, 0x16, 0xfe, 0x8a, 0x06, 0x4a, 0x89, 0xa5, 0xf8, 0xa3, 0xf5,
	0x47, 0x11, 0x4e, 0x7a, 0x86, 0x11, 0x9f, 0xee, 0x3b, 0x33, 0x45, 0x2a, 0x1c, 0xd2, 0x07, 0xea,
	0x46, 0xc1, 0xcc, 0x55, 0xca, 0x2c, 0xb1, 0x79, 0xa3, 0x26, 0x54, 0x57, 0x7e, 0xf0, 0xeb, 0x6c,
	0x3e, 0x8d, 0x57, 0xa7, 0x54, 0x58, 0x5a, 0x0c, 0xb5, 0x7e, 0x2f, 0x43, 0x6d, 0x48, 0x97, 0x33,
	0x97, 0x3e, 0xee, 0xf4, 0x6f, 0x00, 0xee, 0xfc, 0xb9, 0xe1, 0x04, 0x4e, 0xbc, 0x34, 0x7e, 0x2a,
	0x9f, 0xed, 0x1e, 0x69, 0xa6, 0xd7, 0xb9, 0xf3, 0xe7, 0x76, 0x8c, 0xf7, 0xf4, 0x70, 0xaa, 0x15,
	0x88, 0x40, 0x47, 0xdf, 0xc1, 0xd1, 0x9b, 0x74, 0x4b, 0xc9, 0x21, 0x35, 0x72, 0x5a, 0xf9, 0x2d,
	0x6a, 0x05, 0x92, 0x71, 0xd0, 0x8f, 0x50, 0x4f, 0x7b, 0x4d, 0xf2, 0xc9, 0x1d, 0xbd, 0x7a, 0x6b,
	0x49, 0x78, 0x9b, 0xa3, 0x15, 0x48, 0x5e, 0x06, 0x7d, 0x05, 0x07, 0xce, 0x62, 0xc1, 0x3e, 0x67,
	0x15, 0xa6, 0xf8, 0x71, 0x4e, 0x71, 0x7b, 0xcf, 0x5a, 0x81, 0xa4, 0x78, 0xf5, 0x37, 0xa8, 0xe7,
	0x0c, 0x9e, 0x74, 0x05, 0xe2, 0x3e, 0xa5, 0xb7, 0xef, 0xb3, 0xb4, 0xb3, 0x4f, 0xf5, 0x6b, 0x38,
	0x16, 0x07, 0xfe, 0x14, 0xe7, 0xfe, 0x31, 0xc0, 0x84, 0x2e, 0x93, 0x9a, 0x5f, 0xfe, 0x59, 0x82,
	0x83, 0x24, 0x8f, 0x54, 0x38, 0x1b, 0xe2, 0x7b, 0x9b, 0xe0, 0x1f, 0xc6, 0x26, 0xb6, 0x6c, 0xeb,
	0x5a, 0xc7, 0xfc, 0x03, 0x23, 0x17, 0xd0, 0x05, 0xa8, 0x26, 0x26, 0xf7, 0x98, 0xd8, 0x04, 0x9b,
	0xc6, 0xe8, 0xce, 0xc4, 0x62, 0xbe, 0x88, 0x1a, 0xf0, 0x11, 0xe7, 0xee, 0x24, 0xed, 0xde, 0xe0,
	0x46, 0x96, 0xd0, 0x29, 0xc8, 0xa9, 0x38, 0x36, 0xb9, 0xb8, 0x5c, 0x42, 0x0a, 0x9c, 0xee, 0x93,
	0x95, 0xcb, 0xe8, 0x43, 0x78, 0xb6, 0x47, 0x70, 0x70, 0x23, 0x57, 0xd0, 0x33, 0x78, 0x9f, 0xa7,
	0x8c, 0x11, 0xb1, 0x6c, 0xd3, 0xea, 0x59, 0x63, 0x53, 0x7e, 0x2f, 0x66, 0x10, 0xac, 0x8f, 0x2c,
	0x6c, 0xc7, 0xd5, 0x1b, 0xb7, 0xe3, 0xef, 0x59, 0x0e, 0xcb, 0x07, 0xe8, 0x05, 0x34, 0x44, 0x73,
	0x63, 0xf4, 0x1a, 0x93, 0xb1, 0x91, 0x61, 0xc6, 0xa6, 0x7c, 0x88, 0x5a, 0x70, 0x21, 0xf2, 0xf7,
	0x60, 0x8e, 0x50, 0x13, 0xce, 0xb3, 0xaa, 0xf0, 0x36, 0x28, 0x2e, 0x0e, 0xd0, 0x73, 0xf8, 0x40,
	0x28, 0x0e, 0xdf, 0xe2, 0x81, 0x45, 0xae, 0x07, 0x72, 0x35, 0x9b, 0x10, 0xaf, 0xa1, 0xf7, 0xba,
	0xf7, 0x93, 0x3e, 0x1a, 0xe2, 0x54, 0xfb, 0x38, 0x1e, 0xb1, 0xe0, 0x9f, 0xcf, 0xd7, 0xd0, 0x39,
	0x28, 0xa2, 0xf7, 0x06, 0x10, 0xfb, 0x9e, 0xf4, 0x25, 0x4d, 0xfa, 0x37, 0x00, 0x00, 0xff, 0xff,
	0xec, 0x8c, 0x57, 0x39, 0x93, 0x08, 0x00, 0x00,
}
