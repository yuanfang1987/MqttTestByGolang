package commontool

import (
	"fmt"
	"testing"
)

// func Test_RandInt64(t *testing.T) {
// 	a := RandInt64(2, 6)
// 	fmt.Printf("haha, vlaue is: %d\n", a)
// 	fmt.Printf("current time is: %s\n", GetCurrentTime())
// 	<-time.After(5 * time.Second)
// 	fmt.Printf("after 5 seconds: %s\n", GetCurrentTime())
// }

// func Test_func(t *testing.T) {
// 	index := 1
// 	for i := 0; i < 3; i++ {
// 		switch index {
// 		case 1:
// 			fmt.Println("value is 1")
// 			index = 4
// 		case 4:
// 			fmt.Println("value is 4")
// 			index = 9
// 		case 9:
// 			fmt.Println("value is 9")
// 			index = 1
// 		}
// 	}
// }

// func Test_timer(t *testing.T) {
// 	fmt.Printf("current time is: %s\n", GetCurrentTime())
// 	t1 := time.NewTimer(8 * time.Second)
// 	<-t1.C
// 	fmt.Printf("current time is: %s\n", GetCurrentTime())
// 	t1.Reset(5 * time.Second)
// 	<-t1.C
// 	fmt.Printf("current time is: %s\n", GetCurrentTime())
// }

// func Test_WirteFile(t *testing.T) {
// 	f, err := os.OpenFile("D:/temp/golangTemp/test.txt", os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()
// 	f.WriteString("13yuanfang,888\r\n")
// }

// func Test_readdir(t *testing.T) {
// 	// read all .gps file
// 	fileList, err := ioutil.ReadDir("D:/temp/golangTemp/GPS")
// 	if err != nil {
// 		fmt.Printf("read dir error: %s", err.Error())
// 		return
// 	}

// 	// create a new .gps file
// 	newGpsFile, err := os.Create("D:/temp/golangTemp/" + GenerateClientID() + ".gps")
// 	if err != nil {
// 		fmt.Printf("create file error: %s", err.Error())
// 		return
// 	}
// 	defer newGpsFile.Close()

// 	var fContents []string
// 	var readErr error
// 	// and old file and write to new file
// 	for i, f := range fileList {
// 		fmt.Printf("index: %d, file name: %s\n", i, f.Name())
// 		fContents, readErr = ReadFileContent("D:/temp/golangTemp/GPS/" + f.Name())
// 		if readErr != nil {
// 			continue
// 		}
// 		for _, v := range fContents {
// 			newGpsFile.WriteString(v + "\r\n")
// 		}
// 	}
// }

func Test_convertByteToString(t *testing.T) {
	var aa byte
	aa = 1
	fmt.Printf("aa is: %d\n", int(aa))
}
