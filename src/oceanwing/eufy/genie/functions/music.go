package functions

import (
	log "github.com/cihub/seelog"
)

// GetPlayerStatus hh. 4.1
func (b *BaseEufyGenie) GetPlayerStatus(key, expValue string) {
	b.sendGet("/httpapi.asp?command=getPlayerStatus")
	myJSON := b.convertJSON()
	str, err := myJSON.Get(key).String()
	if err != nil {
		return
	}
	if str != expValue {
		log.Infof("verify that play status [%s] should be %s, but actual value is: %s", key,
			expValue, str)
	} else {
		log.Infof("Test passed, play status [%s] is %s", key, str)
	}
}

// PlayMusic ask the device to play music  4.2
func (b *BaseEufyGenie) PlayMusic() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:play")
	okStr := b.getStringResult()
	log.Infof("play music OK? ---> %s", okStr)
}

// PlayPrev hh.  4.3
func (b *BaseEufyGenie) PlayPrev() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:prev")
	okStr := b.getStringResult()
	log.Infof("play previous song OK? ---> %s", okStr)
}

// PlayNext hh. 4.4
func (b *BaseEufyGenie) PlayNext() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:next")
	okStr := b.getStringResult()
	log.Infof("play next song OK? ---> %s", okStr)
}

// FastMoveForwardOrBack hh.  4.5
func (b *BaseEufyGenie) FastMoveForwardOrBack(position string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:seek:" + position)
	okStr := b.getStringResult()
	log.Infof("fast move forward or back OK? ---> %s", okStr)
}

// StopMusic hh. 4.6
func (b *BaseEufyGenie) StopMusic() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:stop")
	okStr := b.getStringResult()
	log.Infof("stop play music OK? ---> %s", okStr)
}

// SetVolume hh. 4.7
func (b *BaseEufyGenie) SetVolume(value string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:vol:" + value)
	okStr := b.getStringResult()
	log.Infof("set vol OK ? ---> %s", okStr)
}

// SetMute hh, 1: mute, 0: unmute 4.8
func (b *BaseEufyGenie) SetMute(value string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:mute:" + value)
	okStr := b.getStringResult()
	log.Infof("set mute the volume OK? ---> %s", okStr)
}

// SetPlayMode hh, 4.9  0 列表循环; 1 单曲循环; 2 随机播放; -1 列表出现
func (b *BaseEufyGenie) SetPlayMode(mod string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:loopmode:" + mod)
	okStr := b.getStringResult()
	log.Infof("set play mode OK? ---> %s", okStr)
}
