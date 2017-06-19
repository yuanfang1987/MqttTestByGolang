package mytest

import (
	"fmt"
)

type myInterface interface {
	myfunc1()
}

type yuanfang struct {
	name string
	age  int
}

func (y *yuanfang) myfunc1() {
	fmt.Printf("Yuanfang's name: %s\n", y.name)
}

type fengzi struct {
	yuanfang
	job string
}

// func (f *fengzi) myfunc1() {
// 	fmt.Printf("fengzi's name: %s\n", f.name)
// }

func test1() {
	var m myInterface
	n := &fengzi{
		job: "IT dog",
	}
	n.name = "fengzi"
	n.age = 28
	m = n
	m.myfunc1()
}
