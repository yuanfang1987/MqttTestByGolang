package crawler

import (
	"encoding/csv"
	"os"
)

var (
	f        *os.File
	writer   *csv.Writer
	inComing chan []string
)

func createNewFile(fpath string) {
	f, _ := os.Create(fpath)
	writer = csv.NewWriter(f)
	title := []string{"CategoryName", "AppName", "URL", "ReviewStar", "ReviewStarNum", "CMD1", "CMD2", "CMD3"}
	writer.Write(title)
	writer.Flush()
	handleIncomingContent()
}

func closeFile() {
	f.Close()
}

func handleIncomingContent() {
	inComing = make(chan []string)
	go func() {
		for {
			select {
			case ss := <-inComing:
				writer.Write(ss)
				writer.Flush()
			}
		}
	}()
}

func writeToResult(content []string) {
	inComing <- content
}
