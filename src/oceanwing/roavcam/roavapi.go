package roavcam

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"oceanwing/eufy/result"
	"strings"
	"time"

	log "github.com/cihub/seelog"
)

var (
	client *HTTPClient
)

// HTTPClient hh..
type HTTPClient struct {
	httpclient *http.Client
	req        chan *http.Request
	resheart   chan []byte
	resFiles   chan []byte
	reqType    int
}

func init() {
	client = &HTTPClient{
		httpclient: &http.Client{Timeout: 30 * time.Second},
		req:        make(chan *http.Request, 2),
		resheart:   make(chan []byte, 2),
		resFiles:   make(chan []byte, 2),
	}
}

// outgoing 是所有HTTP 请求的出口
func (c *HTTPClient) outgoing() {
	go func() {
		for {
			select {
			case request := <-c.req:
				resp, err := c.httpclient.Do(request)
				if err == nil && resp.StatusCode == 200 {
					defer resp.Body.Close()
					bd, _ := ioutil.ReadAll(resp.Body)
					if c.reqType == 1 {
						c.resheart <- bd
					} else if c.reqType == 2 {
						c.resFiles <- bd
					}
					log.Info("send response data to channel res")
				} else {
					log.Infof("request to server fail: %s\n", err)
				}
			}
		}
	}()
}

// handleResponce 用于统一处理HTTP请求的返回数据
func (c *HTTPClient) handleResponce() {
	go func() {
		for {
			select {
			case res1 := <-c.resheart:
				log.Info("心跳返回: ")
				log.Info(string(res1))
			case res2 := <-c.resFiles:
				log.Info("文件列表")
				log.Info(string(res2))
				c.parXMLResult(res2)
			}
		}
	}()
}

func (c *HTTPClient) sendChannel(typ int, re *http.Request) {
	c.reqType = typ
	c.req <- re
}

func parseTime(s string) time.Time {
	layout := "2006/01/02 15:04:05"
	t, err := time.Parse(layout, s)
	if err != nil {
		log.Errorf("parse time fail: %s", err)
	}
	return t
}

func getTimeString(originStr string) string {
	// 2017_0719_120057_001A.MP4
	log.Infof("file name: %s", originStr)
	arrStr := strings.Split(originStr, "_")
	year := arrStr[0]
	month := arrStr[1][0:2]
	day := arrStr[1][2:]
	hour := arrStr[2][0:2]
	m := arrStr[2][2:4]
	s := arrStr[2][4:]
	toBeParse := year + "/" + month + "/" + day + " " + hour + ":" + m + ":" + s
	return toBeParse
}

func (c *HTTPClient) parXMLResult(data []byte) {
	v := &RoavXML{}
	err := xml.Unmarshal(data, v)
	if err != nil {
		log.Errorf("unmarlshal xml fail: %s", err)
		return
	}
	if len(v.AllFile) == 0 {
		return
	}
	for i := 0; i < len(v.AllFile); i++ {
		// 如果达到了最后一个，没有下一个来比较了，直接返回
		if i == (len(v.AllFile) - 1) {
			return
		}
		prev := v.AllFile[i]
		next := v.AllFile[i+1]

		prevEndTime := prev.Time
		nextStartTime := getTimeString(next.Name)

		prevTi := parseTime(prevEndTime)
		nextTi := parseTime(nextStartTime)

		if nextTi.Unix()-prevTi.Unix() > 5 {
			content := []string{prev.Name, prev.Time, next.Name, next.Time, "文件不连续"}
			result.WriteToExcel(content)
			log.Infof("assert fail: %v", content)
		}

		// log
		log.Infof("Prev file name: %s, end time: %s", prev.Name, prevEndTime)
		log.Infof("Next file name: %s, start time: %s", next.Name, nextStartTime)
		log.Infof("prev and next file diff time: %d seconds", nextTi.Unix()-prevTi.Unix())
	}

}

// SendRoavAPI hh.
func SendRoavAPI() {
	client.outgoing()
	client.handleResponce()

	heartBeatReq, _ := http.NewRequest("GET", "http://192.168.1.254/?custom=1&cmd=3016", nil)
	getFileListReq, _ := http.NewRequest("GET", "http://192.168.1.254/?custom=1&cmd=3015", nil)

	interval1 := time.NewTicker(time.Second * 5).C
	interval2 := time.NewTicker(time.Second * 600).C

	for {
		select {
		case <-interval1:
			client.sendChannel(1, heartBeatReq)
		case <-interval2:
			client.sendChannel(2, getFileListReq)
		}
	}

}

// RunService hh.
func RunService() {
	client.outgoing()
	client.handleResponce()
}

// SendHeartBeat hh.
func SendHeartBeat() {
	heartBeatReq, _ := http.NewRequest("GET", "http://192.168.1.254/?custom=1&cmd=3016", nil)
	client.sendChannel(1, heartBeatReq)
}

// GetFileListAndAssert h.
func GetFileListAndAssert() {
	getFileListReq, _ := http.NewRequest("GET", "http://192.168.1.254/?custom=1&cmd=3015", nil)
	client.sendChannel(2, getFileListReq)
}
