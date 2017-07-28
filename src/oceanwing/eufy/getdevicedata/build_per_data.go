package getdevicedata

import (
	"fmt"
)

func doit() {
	//查询数据库
	deviceList := getDevKeyAndDevID("T1011")
	for _, dev := range deviceList {
		// register user
		tryCounter := 0
		var u *user
		for {
			u = registerUser()
			if u != nil {
				break
			}
			tryCounter++
			if tryCounter == 5 {
				break
			}
		}

		if u == nil {
			fmt.Printf("dev key not bind: %s", dev["devkey"])
			continue
		}

		// bind device
		tryCounter2 := 0
		var b bool
		for {
			b = bindDevice(dev["devkey"], u)
			if b {
				break
			}
			tryCounter2++
			if tryCounter2 == 5 {
				break
			}
		}
		// 插数据到 timer

		// 插数据到 timer_away
	}
}
