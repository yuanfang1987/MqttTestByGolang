package mytest

import (
	"fmt"
	"testing"
)

func Test_bufferchan(t *testing.T) {
	c := make(chan string, 2)
	fmt.Printf("new channel len: %d\n", len(c))
	c <- "yuanfang"
	c <- "matt"
	fmt.Printf("chan len: %d\n", len(c))
	<-c
	fmt.Printf("chan len after get one: %d\n", len(c))
	<-c
	fmt.Printf("chan len after get two: %d\n", len(c))
}

// func Test_pbhaha(t *testing.T) {
// 	pl := pbMarshal()
// 	pbUnMarshal(pl)
// }

// func Test_binary(t *testing.T) {
// 	fmt.Printf("binary: %b\n", 3)
// }

// func Test_binary2(t *testing.T) {
// 	var weekinfo string
// 	v := 101
// 	if (v & 1) > 0 {
// 		weekinfo = "星期一, "
// 	}
// 	if (v & 2) > 0 {
// 		weekinfo = weekinfo + "星期二, "
// 	}
// 	if (v & 4) > 0 {
// 		weekinfo = weekinfo + "星期三, "
// 	}
// 	if (v & 8) > 0 {
// 		weekinfo = weekinfo + "星期四, "
// 	}
// 	if (v & 16) > 0 {
// 		weekinfo = weekinfo + "星期五, "
// 	}
// 	if (v & 32) > 0 {
// 		weekinfo = weekinfo + "星期六, "
// 	}
// 	if (v & 64) > 0 {
// 		weekinfo = weekinfo + "星期日"
// 	}
// 	fmt.Println(weekinfo)
// }
