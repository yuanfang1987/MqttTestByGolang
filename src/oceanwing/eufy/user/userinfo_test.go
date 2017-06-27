package user

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_user(t *testing.T) {
	// commontool.InitLogInstance("debug")
	// me := NewUser("matt.ma@oceanwing.com", "Lin910528&", "eufy-app", "8FHf22gaTKu7MZXqz5zytw")
	// me.Login()
	// // me.SetAwayMode(3, 120, "1fc2fca2-3e9f-4463-a243-2385b8390bea")
	// me.StopAwayMode("1fc2fca2-3e9f-4463-a243-2385b8390bea")
	// me.GetAwayModeInfo("1fc2fca2-3e9f-4463-a243-2385b8390bea")

	later := time.Now().Add(2 * time.Minute)
	laterHour := later.Hour()
	laterMinute := later.Minute()
	counter := 0
	for {
		time.Sleep(1 * time.Second)
		nn := time.Now()
		nnHour := nn.Hour()
		nnMinute := nn.Minute()
		fmt.Printf("expected time %d:%d, now time %d:%d\n", laterHour, laterMinute, nnHour, nnMinute)
		if laterHour == nnHour && laterMinute == nnMinute {
			fmt.Println("it's the time now.")
			break
		}
		counter++
		fmt.Printf("%d second pass...\n", counter)
	}

	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
}
