package mytest

import (
	"fmt"
	"testing"
	"time"
)

// func Test_pbhaha(t *testing.T) {
// 	pl := pbMarshal()
// 	pbUnMarshal(pl)
// }

// func Test_relfectAA(t *testing.T) {
// 	var p interface{}
// 	p = &Person{
// 		name: "yuanfang",
// 		age:  30,
// 	}

// 	t1 := reflect.TypeOf(p)
// 	v := reflect.ValueOf(p)
// 	fmt.Printf("ValueOf(), %v\n", v)
// 	fmt.Printf("TypeOf(), %v\n", t1)

// 	v2 := v.Elem()
// 	tv := t1.Elem()
// 	fmt.Printf("Elem(), %v\n", v2)
// 	fmt.Printf("Elem(), %v\n", tv)

// 	name1 := v2.FieldByName("name")
// 	fmt.Printf("name: %s\n", name1)

// 	v2typ := v2.Type()
// 	fmt.Printf("v2 type: %v", v2typ)
// }

func Test_timeParse(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("location invalid: %s\n", err)
		return
	}
	fmt.Printf("current location is: %s\n", loc.String())
	lay := "2006-01-02 15:04:05"
	now := time.Now()
	expectedTime := fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), 13, 20, 0)
	ti, err := time.ParseInLocation(lay, expectedTime, loc)
	if err != nil {
		fmt.Printf("parse time error: %s\n", err)
		return
	}
	fmt.Printf("parse time is: %d:%d\n", ti.Hour(), ti.Minute())
}
