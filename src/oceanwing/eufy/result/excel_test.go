package result

import (
	"fmt"
	"testing"
)

func Test_aa(t *testing.T) {
	InitExcelFile()
	arr := []string{"haha", "hehe", "qinqin"}
	arr2 := []string{"fengzi", "shejing"}
	arr3 := []string{"shenzhen", "guangzhou", "nanning", "qinzhou"}
	WriteToExcel(arr)
	WriteToExcel(arr2)
	WriteToExcel(arr3)
	SaveExcelFile()
}

func Test_ssli(t *testing.T) {
	ss := []string{"aa", "bb", "cc"}
	bb := []string{"dd", "ee", "ff"}
	ss = append(ss, bb...)
	for _, v := range ss {
		fmt.Println(v)
	}
}
