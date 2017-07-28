package result

import (
	"oceanwingqa/common/utils"

	log "github.com/cihub/seelog"
	"github.com/tealeg/xlsx"
)

var (
	xFile  *xlsx.File
	xSheet *xlsx.Sheet
	arrStr chan []string
	err    error
)

// InitExcelFile 初始化一个excel变量
func InitExcelFile() {
	xFile = xlsx.NewFile()
	xSheet, err = xFile.AddSheet("result")
	if err != nil {
		panic("init excel file error")
	}
	arrStr = make(chan []string)
	handleWriting()
}

// WriteToExcel hh.
func WriteToExcel(ss []string) {
	go func() {
		arrStr <- ss
	}()
}

// SaveExcelFile save a file.
func SaveExcelFile() {
	err = xFile.Save(utils.GetTimeAsFileName() + "-result.xlsx")
	if err != nil {
		log.Errorf("save excel file error: %s", err)
	}
}

func writeResult(contents []string) {
	row := xSheet.AddRow()
	for _, v := range contents {
		row.AddCell().SetString(v)
	}
}

func handleWriting() {
	go func() {
		for {
			select {
			case ss := <-arrStr:
				writeResult(ss)
			}
		}
	}()
}
