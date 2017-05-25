package results

import (
	"encoding/csv"
	"oceanwing/commontool"
	"os"
	"strconv"
)

var (
	f           *os.File
	writer      *csv.Writer
	inComingStr chan []string
	counter     int
)

// NewResultFile hh.
func NewResultFile(fPath string) {
	f, _ := os.Create(fPath)
	writer = csv.NewWriter(f)
	title := []string{"编号", "类别", "Case名称", "测试结果", "测试时间"}
	writer.Write(title)
	writer.Flush()
	handleIncomingContent()
}

// CloseResultFile hh.
func CloseResultFile() {
	f.Close()
}

// WriteToResultFile hh.
func WriteToResultFile(args ...string) {
	counter++
	index := strconv.Itoa(counter)
	var content []string
	// 编号， 自动增长
	content = append(content, index)
	// 类别，case, 结果
	content = append(content, args...)
	// 时间，自动生成
	content = append(content, commontool.GetCurrentTime())
	inComingStr <- content
}

func handleIncomingContent() {
	inComingStr = make(chan []string)
	go func() {
		for {
			select {
			case ss := <-inComingStr:
				writer.Write(ss)
				writer.Flush()
			}
		}
	}()
}
