package cases

import (
	log "github.com/cihub/seelog"
)

// RunMusicCases hh.
func RunMusicCases(url string) {
	log.Info("========= Running music test cases ==========")
	// new instance
	newTestInstance(url)
	log.Info("create a new eufygenie instance.")

	// play
	log.Info("start play music")
	Instance.PlayMusic()
	Instance.GetPlayerStatus("status", "play")

	// play prev
	log.Info("start play previous song")
	Instance.PlayPrev()
	Instance.GetPlayerStatus("status", "play")

	// play next
	log.Info("start play next song")
	Instance.PlayNext()
	Instance.GetPlayerStatus("status", "play")

	// fast forward 5 sec
	log.Info("start fast move forward to 5 second.")
	Instance.FastMoveForwardOrBack("5")
	Instance.GetPlayerStatus("status", "play")

	// fast back  3 sec
	log.Info("start fast move back to 3 second")
	Instance.FastMoveForwardOrBack("3")
	Instance.GetPlayerStatus("status", "play")

	// set volume 80 percent?
	log.Info("start set volume to 80")
	Instance.SetVolume("80")
	Instance.GetPlayerStatus("vol", "80")

	// set mute
	log.Info("start set mute to 1")
	Instance.SetMute("1")
	Instance.GetPlayerStatus("mute", "1")

	// set unmute
	log.Info("start set mute to 0")
	Instance.SetMute("0")
	Instance.GetPlayerStatus("mute", "0")

	// set play mode: 0
	log.Info("start set play mode to 0")
	Instance.SetPlayMode("0")
	Instance.GetPlayerStatus("loop", "0")

	// set play mode: 1
	log.Info("start set play mode to 1")
	Instance.SetPlayMode("1")
	Instance.GetPlayerStatus("loop", "1")

	// set play mode: 2
	log.Info("start set play mode to 2")
	Instance.SetPlayMode("2")
	Instance.GetPlayerStatus("loop", "2")

	// stop play
	log.Info("stop playing musci")
	Instance.StopMusic()
	Instance.GetPlayerStatus("status", "stop")
}
