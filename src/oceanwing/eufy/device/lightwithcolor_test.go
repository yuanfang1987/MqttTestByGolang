package device

import (
	"fmt"
	"testing"
)

func Test_rgb(t *testing.T) {
	getRGBData()

	for _, v := range rgb {
		name := v["RGBName"]
		rgbConf := v["rgbConfig"]
		if vv, ok := rgbConf.(*rgbInfo); ok {
			fmt.Printf("%s, red: %d, green: %d, blue: %d\n", name, vv.red, vv.green, vv.blue)
		}
	}
}
