package main

import (
	"fmt"
	"io/ioutil"
	"oceanwingqa/common/utils"
	"os"
	"strings"
)

func main() {
	// read all .gps file.
	fileList, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Printf("read dir files error: %s", err.Error())
		return
	}

	// create a new .gps file
	newGpsFile, err := os.Create("./" + utils.GenerateClientID() + "Merge.gps")
	if err != nil {
		fmt.Printf("create file error: %s", err.Error())
		return
	}
	defer newGpsFile.Close()

	var fContents []string
	var readErr error
	// and old file and write to new file
	for i, f := range fileList {
		fmt.Printf("index: %d, file name: %s\n", i, f.Name())
		if !strings.Contains(f.Name(), ".gps") {
			continue
		}
		fContents, readErr = utils.ReadFileContent(f.Name())
		if readErr != nil {
			fmt.Printf("read file: %s, error: %s", f.Name(), readErr.Error())
			continue
		}
		for _, v := range fContents {
			newGpsFile.WriteString(v + "\r\n")
		}
	}
}
