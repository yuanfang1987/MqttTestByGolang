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

func Test_reflect1(t *testing.T) {
	fmt.Println("run Test_reflect1")
	var x = 3.4
	v := reflect.ValueOf(x)
	fmt.Println("Type: ", v.Type())
	fmt.Println("kind is float64 :", v.Kind() == reflect.Float64)
	fmt.Println("value: ", v.Float())
}

func Test_reflect2(t *testing.T) {
	fmt.Println("run Test_reflect2")
	type MyInt int
	var x MyInt = 7
	v := reflect.ValueOf(x)
	fmt.Println("Type: ", v.Type())
	fmt.Println("Kind: ", v.Kind())
}

func Test_reflect3(t *testing.T) {
	fmt.Println("run Test_reflect3")
	tonydon := &User{"TangXiaodong", 100, "0000123"}
	object := reflect.ValueOf(tonydon)
	fmt.Println("object Type: ", object.Type())
	fmt.Println("object Kind: ", object.Kind())
	myref := object.Elem()
	fmt.Println("myref Type: ", myref.Type())
	fmt.Println("myref Kind: ", myref.Kind())
	typeOfType := myref.Type()
	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)
		fmt.Printf("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
	}
	v := object.MethodByName("SayHello")
	fmt.Println("v Type: ", v.Type())
	fmt.Println("v Kind: ", v.Kind())
	v.Call([]reflect.Value{})
}
