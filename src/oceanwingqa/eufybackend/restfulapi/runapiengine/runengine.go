package runapiengine

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"oceanwingqa/common/utils"
	"oceanwingqa/eufybackend/restfulapi/api/base"
	"oceanwingqa/eufybackend/restfulapi/api/build"
	"strings"
	"time"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
	"github.com/tealeg/xlsx"
)

var (
	xFile         *xlsx.File
	testCaseSheet *xlsx.Sheet
	testDataSheet *xlsx.Sheet
	client        *HTTPClient
	err           error
	serverHost    string
)

// HTTPClient hh..
type HTTPClient struct {
	httpclient *http.Client
	req        chan *http.Request
	res        chan []byte
	jsonResult chan *splJSON.Json
	uid        string
	token      string
}

func init() {
	xFile, err = xlsx.OpenFile("eufyApiTest.xlsx")
	if err != nil {
		panic("open file error.")
	}
	testCaseSheet = xFile.Sheet["TestCase"]
	testDataSheet = xFile.Sheet["TestData"]
	// new a http client
	client = &HTTPClient{
		httpclient: &http.Client{Timeout: 30 * time.Second},
		req:        make(chan *http.Request),
		res:        make(chan []byte),
		jsonResult: make(chan *splJSON.Json),
	}
	client.outgoing()
	client.handleResponce()
}

// SetHostURL hh.
func SetHostURL(url string) {
	serverHost = url
}

// Runapitest 运行测试入口
func Runapitest() {
	var api base.RESTfulAPI
	for i, row := range testCaseSheet.Rows {
		// 跳过第一行, 列头名
		if i == 0 {
			log.Debug("This is the first row.")
			continue
		}
		cells := row.Cells
		// 跳过空行
		if len(cells) == 0 {
			log.Debug("This is a empty row")
			continue
		}

		category, _ := cells[4].String()
		urlPath, _ := cells[5].String()
		httpMethod, _ := cells[6].String()
		testDataID, _ := cells[7].String()

		// 获取测试数据
		testDataMap := getTestData(testDataID)
		// 根据 category 新建一个 API 实例
		api = build.CreateNewAPI(category, urlPath, httpMethod, testDataMap)
		// 根据测试数据，构造出一个请求 body.
		actURL, body := api.BuildRequestBody()
		// 重新构造URL
		if actURL != "" {
			urlPath = actURL
		}
		// 发出请求并获取结果
		jsonResponse := client.doItNow(httpMethod, urlPath, body, testDataMap)
		// 解析并判断结果
		resultMAP := api.DecodeAndAssertResult(jsonResponse)
		// 如果存在 uid 和 token， 则更新 client 的 uid 和 token，因为可能使用了新的用户登录
		if UID, ok := resultMAP["uid"]; ok {
			client.uid = UID
		}
		if TOKEN, ok := resultMAP["token"]; ok {
			client.token = TOKEN
		}
		resultString := resultMAP["resultString"]
		// 把结果写入原文件
		cells[8].SetString(passOrNot(resultString))
		if resultString != "" {
			cells[9].SetString(resultString)
		}
	}
	// 生成一个新的结果文件
	newFileName := utils.GetTimeAsFileName() + "-TestResult.xlsx"
	err = xFile.Save(newFileName)
	if err != nil {
		log.Errorf("Save test result file error: %s", err)
	}
}

func getTestData(testdataID string) map[string]string {
	// 初始化一个字典
	dataMap := make(map[string]string)
	for _, row := range testDataSheet.Rows {
		id, _ := row.Cells[0].String()
		// 找到 test data id 所在的 row
		if id == testdataID {
			for i, cell := range row.Cells {
				// 跳过第 1 个单元格
				if i == 0 {
					continue
				}
				text, e := cell.String()
				// 成功取出单元格的值，且不为空，且分割后有两个值
				if e == nil && text != "" {
					str := strings.Split(text, ":")
					if len(str) > 1 {
						dataMap[str[0]] = str[1]
						log.Infof("Get test data, key: %s, value: %s", str[0], str[1])
					}
				}
			}
			break
		}
	}
	return dataMap
}

func passOrNot(str string) string {
	if str == "" {
		return "Pass"
	}
	return "Fail"
}

// 构建请求体，包括必要的 header及 post body
func (c *HTTPClient) buildRequest(method, url string, bd []byte, headerMap map[string]string) *http.Request {
	body := bytes.NewBuffer(bd)
	req, err := http.NewRequest(method, serverHost+url, body)
	if err != nil {
		log.Errorf("build new request error: %s", err)
		return nil
	}

	// 时区
	if tz, ok := headerMap["timezone"]; ok {
		req.Header.Add("timezone", tz)
	} else {
		req.Header.Add("timezone", "Asia/Shanghai")
	}

	// 国家
	if ct, ok := headerMap["country"]; ok {
		req.Header.Add("country", ct)
	} else {
		req.Header.Add("country", "CN")
	}

	// 语言
	if lang, ok := headerMap["language"]; ok {
		req.Header.Add("language", lang)
	} else {
		req.Header.Add("language", "zh-Hans-CN")
	}

	req.Header.Add("openudid", "yuanfang1987") // 这个值没什么意义，随便填， 没有也行.

	// category
	if cate, ok := headerMap["category"]; ok {
		req.Header.Add("category", cate)
	} else {
		req.Header.Add("category", "eufy-app")
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	if c.uid != "" {
		req.Header.Add("uid", c.uid)
	}
	if c.token != "" {
		req.Header.Add("token", c.token)
	}
	log.Debugf("build a request with method: %s, and url: %s", method, url)
	return req
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
					c.res <- bd
					log.Info("send response data to channel res")
				} else {
					log.Errorf("request to server fail: %s", err)
					c.res <- nil
				}
			}
		}
	}()
}

// handleResponce 用于统一处理HTTP请求的返回数据，主要是转为simple JSON的对象
func (c *HTTPClient) handleResponce() {
	go func() {
		for {
			select {
			case resp := <-c.res:
				JSONJinstance, err := splJSON.NewJson(resp)
				if err != nil {
					c.jsonResult <- nil
					log.Errorf("decode the res data to simpleJSON error: %s", err)
				} else {
					c.jsonResult <- JSONJinstance
					log.Info("send simpleJSON to jsonResult")
				}
			}
		}
	}()
}

func (c *HTTPClient) doItNow(method, url string, body []byte, headers map[string]string) *splJSON.Json {
	c.req <- c.buildRequest(method, url, body, headers)
	j := <-c.jsonResult
	log.Infof("Get response body: %v", j)
	return j
}
