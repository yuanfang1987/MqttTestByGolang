package main

import (
	"flag"
	"oceanwing/commontool"
	"oceanwing/crawler"

	log "github.com/cihub/seelog"
)

// haha.
func main() {
	defer log.Flush()
	specifiedCateIndex := flag.Int("CateIndex", 0, "category index")
	specifiedPageIndex := flag.Int("PageIndex", 0, "page index")
	logLevel := flag.String("logLevel", "debug", "log level")
	flag.Parse()
	commontool.InitLogInstance(*logLevel)
	heads := make(map[string]string)
	heads["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0"
	heads["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	heads["Accept-Encoding"] = "gzip, deflate, sdch, br"
	heads["Accept-Language"] = "zh-CN,zh;q=0.8"
	heads["Connection"] = "keep-alive"
	crawler.RunVersion2("https://www.amazon.com/b?ie=UTF8&node=13727921011", heads, *specifiedCateIndex, *specifiedPageIndex)
}
