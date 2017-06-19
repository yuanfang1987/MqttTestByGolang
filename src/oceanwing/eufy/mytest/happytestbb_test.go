package mytest

import (
	"fmt"
	"testing"
	"time"
)

func Test_funcbb1(t *testing.T) {
	test1()
}

func Test_timeAdd(t *testing.T) {
	now := time.Now()
	fmt.Printf("Current time, hour: %d, minutes: %d, seconds: %d\n", now.Hour(), now.Minute(), now.Second())
	next := now.Add(10 * time.Minute)
	fmt.Printf("after add 10 minutes, hour: %d, minutes: %d, seconds: %d\n", next.Hour(), next.Minute(), next.Second())
}
