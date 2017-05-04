package deviceactionapi

import (
	log "github.com/cihub/seelog"
)

// ActionRobotCleanerReq 发送MQTT命令的BODY
type ActionRobotCleanerReq struct {
	CleanSpeed    CleanSpeedOption    `json:"clean_speed,omitempty"`
	FindMe        FindMeOption        `json:"find_me,omitempty"`
	StopClean     StopCleanOption     `json:"stop_clean,omitempty"`
	TimerOption   TimerOption         `json:"timer_option,omitempty"`
	TurnDirection TurnDirectionOption `json:"turn_direction,omitempty"`
	WorkMode      WorkModeOption      `json:"work_mode,omitempty"`
}

// CleanSpeedOption hh
type CleanSpeedOption struct {
	Speed int `json:"speed"` // 速度选项：0:日常 1:强力 2:地毯 3:静音
}

// FindMeOption hh.
type FindMeOption struct {
	OnOff int `json:"on_off"` // 0: 关闭findme功能 扫地机停止发声; 1:开启findme功能,扫地机持续鸣叫
}

// StopCleanOption hh.
type StopCleanOption struct {
	Stop int `json:"stop "` // 1:停止; 0:对应工作模式
}

// TimerOption hh.
type TimerOption struct {
	ScheduleType string            `json:"schedule_type "` //Timer调度类型, 目前可选值为 weekly, 后续版本会增加daily等
	WeeklyOption WeeklyTimerOption `json:"weekly_option"`
}

// WeeklyTimerOption hh.
type WeeklyTimerOption struct {
	StartHour   int `json:"start_hour"`   //开始执行的小时, 24小时制, 可选值为 0~23
	StartMinute int `json:"start_minute"` //开始执行的分钟, 60分钟制, 可选值为 0~59 ,
	Weekday     int `json:"weekday"`      //分别对应(Sunday=0, Monday=1, ..., Saturday=6)
}

// TurnDirectionOption hh.
type TurnDirectionOption struct {
	Direction int `json:"direction"` //0:Forward, 1:Backward, 2:Left, 3:Right
}

// WorkModeOption hh.
type WorkModeOption struct {
	Mode int `json:"mode"` //0:暂停, 1:定点, 2:自动, 3:返回充电, 4:沿边, 5:精扫
}

// ActionRobotCleanerResp 解析返回值
type ActionRobotCleanerResp struct {
	Message string `json:"message"`
	ResCode int    `json:"res_code"` //1成功 0失败, 2001, Device Action Result is pending, need check status later
}

//---------------------------------------------------------------------------------------------------------------------

// NewActionRobotCleanerReq hh.
func NewActionRobotCleanerReq(option string, d ...int) interface{} {
	if option == "findMe" {
		log.Debugf("do something for %s", option)
	}
	if option == "turnDirection" {
		log.Debugf("do something for %s", option)
	}
	req := &ActionRobotCleanerReq{
		CleanSpeed: CleanSpeedOption{
			Speed: d[0],
		},
		FindMe: FindMeOption{
			OnOff: d[1],
		},
		StopClean: StopCleanOption{
			Stop: d[2],
		},
		WorkMode: WorkModeOption{
			Mode: d[3],
		},
	}
	return req
}
