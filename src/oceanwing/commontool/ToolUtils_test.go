package commontool

import (
	"fmt"
	"testing"
	"time"
)

// func Test_RandInt64(t *testing.T) {
// 	a := RandInt64(2, 6)
// 	fmt.Printf("haha, vlaue is: %d\n", a)
// 	fmt.Printf("current time is: %s\n", GetCurrentTime())
// 	<-time.After(5 * time.Second)
// 	fmt.Printf("after 5 seconds: %s\n", GetCurrentTime())
// }

func Test_func(t *testing.T) {
	index := 1
	for i := 0; i < 3; i++ {
		switch index {
		case 1:
			fmt.Println("value is 1")
			index = 4
		case 4:
			fmt.Println("value is 4")
			index = 9
		case 9:
			fmt.Println("value is 9")
			index = 1
		}
	}
}

func Test_timer(t *testing.T) {
	fmt.Printf("current time is: %s\n", GetCurrentTime())
	t1 := time.NewTimer(8 * time.Second)
	<-t1.C
	fmt.Printf("current time is: %s\n", GetCurrentTime())
	t1.Reset(5 * time.Second)
	<-t1.C
	fmt.Printf("current time is: %s\n", GetCurrentTime())
}
