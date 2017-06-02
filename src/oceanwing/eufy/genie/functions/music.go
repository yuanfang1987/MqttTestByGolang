package functions

import (
	"oceanwing/eufy/genie/results"

	log "github.com/cihub/seelog"
)

const categoryMusic = "Music"

// GetPlayerStatus hh. 4.1
func (b *BaseEufyGenie) GetPlayerStatus(key, expValue string) {
	b.sendGet("/httpapi.asp?command=getPlayerStatus")
	log.Info("execute get music play status.")
	myJSON := b.convertJSON()
	str, err := myJSON.Get(key).String()
	if err != nil {
		return
	}
	results.WriteToResultFile(categoryMusic, "verify key "+key+" is "+expValue, passOrFail(expValue == str))
	log.Infof("verify key: %s, expected value: %s, actual value: %s, test case passed or not? ---> %t",
		key, expValue, str, expValue == str)
}

// PlayMusic ask the device to play music  4.2
func (b *BaseEufyGenie) PlayMusic() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:play")
	log.Info("execute play music.")
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "play music", re)
	log.Infof("play music OK ---> %s", re)
}

// PlayPrev hh.  4.3
func (b *BaseEufyGenie) PlayPrev() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:prev")
	log.Info("execute play previous song.")
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "play previous song", re)
	log.Infof("play previous song OK? ---> %s", re)
}

// PlayNext hh. 4.4
func (b *BaseEufyGenie) PlayNext() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:next")
	log.Info("execute play next song.")
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "play next song", re)
	log.Infof("play next song OK? ---> %s", re)
}

// FastMoveForwardOrBack hh.  4.5
func (b *BaseEufyGenie) FastMoveForwardOrBack(position string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:seek:" + position)
	log.Infof("execute fast move forward or back with position: %s", position)
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "fast forward or back with position: "+position, re)
	log.Infof("fast move forward or back OK? ---> %s", re)
}

// StopMusic hh. 4.6
func (b *BaseEufyGenie) StopMusic() {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:stop")
	log.Info("execute stop play music")
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "stop play music", re)
	log.Infof("stop play music OK? ---> %s", re)
}

// SetVolume hh. 4.7
func (b *BaseEufyGenie) SetVolume(value string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:vol:" + value)
	log.Infof("execute set music volume with value: %s", value)
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "set volume to "+value, re)
	log.Infof("set vol OK ? ---> %s", re)
}

// SetMute hh, 1: mute, 0: unmute 4.8
func (b *BaseEufyGenie) SetMute(value string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:mute:" + value)
	log.Infof("execute set music mute with value: %s", value)
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "set mute to "+value, re)
	log.Infof("set mute the volume OK? ---> %s", re)
}

// SetPlayMode hh, 4.9  0 列表循环; 1 单曲循环; 2 随机播放; -1 列表出现
func (b *BaseEufyGenie) SetPlayMode(mod string) {
	b.sendGet("/httpapi.asp?command=setPlayerCmd:loopmode:" + mod)
	log.Infof("execute set play mode with value: %s", mod)
	okStr := b.getStringResult()
	re := passOrFail(okStr == "OK")
	results.WriteToResultFile(categoryMusic, "set play mode to "+mod, re)
	log.Infof("set play mode OK? ---> %s", re)
}

// GetPlayerStatusValue hh.
func (b *BaseEufyGenie) GetPlayerStatusValue(key string) string {
	b.sendGet("/httpapi.asp?command=getPlayerStatus")
	log.Info("execute get music play status.")
	myJSON := b.convertJSON()
	iuri, err := myJSON.Get(key).String()
	if err == nil {
		return hexToString(iuri)
	}
	return ""
}
