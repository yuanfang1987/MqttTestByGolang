package user

import (
	"oceanwing/commontool"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func Test_user(t *testing.T) {
	commontool.InitLogInstance("debug")
	me := NewUser("matt.ma@oceanwing.com", "Lin910528&", "eufy-app", "8FHf22gaTKu7MZXqz5zytw")
	me.Login()
	// me.SetAwayMode("1fc2fca2-3e9f-4463-a243-2385b8390bea")
	me.StopAwayMode("1fc2fca2-3e9f-4463-a243-2385b8390bea")
	me.GetAwayModeInfo("1fc2fca2-3e9f-4463-a243-2385b8390bea")

	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
}
