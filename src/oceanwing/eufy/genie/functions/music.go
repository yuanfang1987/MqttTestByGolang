package functions

import (
	log "github.com/cihub/seelog"
)

// GetPlayerStatus hh. 4.1
func (b *BaseEufyGenie) GetPlayerStatus(key, expValue string) {
	b.sendGet("/httpapi.asp?command=getPlayerStatus")
	log.Info("execute get music play status.")
	myJSON := b.convertJSON()
	str, err := myJSON.Get(key).String()
	if err != nil {
		return
	}
	log.Infof("verify key: %s, expected value: %s, actual value: %s, test case passed or not? ---> %t",
		key, expValue, str, expValue == str)
}

// PlayMusic ask the device to play music  4.2
func (b *BaseEufyGenie) PlayMusic() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:play")
	log.Info("execute play music.")
	okStr := b.getStringResult()
	log.Infof("play music OK ---> %t", okStr == "OK")
}

// PlayPrev hh.  4.3
func (b *BaseEufyGenie) PlayPrev() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:prev")
	log.Info("execute play previous song.")
	okStr := b.getStringResult()
	log.Infof("play previous song OK? ---> %t", okStr == "OK")
}

// PlayNext hh. 4.4
func (b *BaseEufyGenie) PlayNext() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:next")
	log.Info("execute play next song.")
	okStr := b.getStringResult()
	log.Infof("play next song OK? ---> %t", okStr == "OK")
}

// FastMoveForwardOrBack hh.  4.5
func (b *BaseEufyGenie) FastMoveForwardOrBack(position string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:seek:" + position)
	log.Infof("execute fast move forward or back with position: %s", position)
	okStr := b.getStringResult()
	log.Infof("fast move forward or back OK? ---> %t", okStr == "OK")
}

// StopMusic hh. 4.6
func (b *BaseEufyGenie) StopMusic() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:stop")
	log.Info("execute stop play music")
	okStr := b.getStringResult()
	log.Infof("stop play music OK? ---> %t", okStr == "OK")
}

// SetVolume hh. 4.7
func (b *BaseEufyGenie) SetVolume(value string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:vol:" + value)
	log.Infof("execute set music volume with value: %s", value)
	okStr := b.getStringResult()
	log.Infof("set vol OK ? ---> %t", okStr == "OK")
}

// SetMute hh, 1: mute, 0: unmute 4.8
func (b *BaseEufyGenie) SetMute(value string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:mute:" + value)
	log.Infof("execute set music mute with value: %s", value)
	okStr := b.getStringResult()
	log.Infof("set mute the volume OK? ---> %t", okStr == "OK")
}

// SetPlayMode hh, 4.9  0 列表循环; 1 单曲循环; 2 随机播放; -1 列表出现
func (b *BaseEufyGenie) SetPlayMode(mod string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:loopmode:" + mod)
	log.Infof("execute set play mode with value: %s", mod)
	okStr := b.getStringResult()
	log.Infof("set play mode OK? ---> %t", okStr == "OK")
}
