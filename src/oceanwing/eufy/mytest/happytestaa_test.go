package mytest

import (
	"fmt"
	"reflect"
	"testing"
)

// func Test_pbhaha(t *testing.T) {
// 	pl := pbMarshal()
// 	pbUnMarshal(pl)
// }

func Test_relfectAA(t *testing.T) {
	var p interface{}
	p = &Person{
		name: "yuanfang",
		age:  30,
	}

	t1 := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	fmt.Printf("ValueOf(), %v\n", v)
	fmt.Printf("TypeOf(), %v\n", t1)

	v2 := v.Elem()
	tv := t1.Elem()
	fmt.Printf("Elem(), %v\n", v2)
	fmt.Printf("Elem(), %v\n", tv)

	name1 := v2.FieldByName("name")
	fmt.Printf("name: %s\n", name1)

	v2typ := v2.Type()
	fmt.Printf("v2 type: %v", v2typ)
}
