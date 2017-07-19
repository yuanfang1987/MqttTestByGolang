package roavcam

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test_asdf11(t *testing.T) {
	//myXml22()
}

func Test_partt(t *testing.T) {
	lay := "2006/01/02 15:04:05"
	// loc, err := time.LoadLocation("Local")
	// if err != nil {
	// 	fmt.Printf("get loc error: %s", err)
	// 	return
	// }

	t1, err := time.Parse(lay, "2017/07/19 12:01:42")

	if err != nil {
		fmt.Printf("parse error: %s\n", err)
		return
	}
	fmt.Printf("hour: %d\n", t1.Hour())
	fmt.Printf("minute: %d\n", t1.Minute())
	fmt.Printf("second: %d\n", t1.Second())

	t2, err := time.Parse(lay, "2017/07/19 12:01:52")
	diff := t2.Unix() - t1.Unix()
	fmt.Printf("diff value: %d\n", diff)

	fmt.Println("===============================================")
}

func Test_partt2(t *testing.T) {
	originStr := "2017_0719_120057_001A.MP4"
	arrStr := strings.Split(originStr, "_")
	year := arrStr[0]
	month := arrStr[1][0:2]
	day := arrStr[1][2:]
	hour := arrStr[2][0:2]
	m := arrStr[2][2:4]
	s := arrStr[2][4:]
	toBeParse := year + "/" + month + "/" + day + " " + hour + ":" + m + ":" + s
	fmt.Printf("expected value: %s\n", toBeParse)

	fmt.Println("================================================")
}
