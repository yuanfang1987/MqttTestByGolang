package main

import (
	"oceanwing/eufy/genie"
)

func main() {
	myGenie := genie.NewEufyGenie("http://10.10.10.254")
	myGenie.GetPlayerStatus("status", "play")
}
