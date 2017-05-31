package crawler

import (
	"net/http"
	"oceanwing/commontool"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
	log "github.com/cihub/seelog"
)

// AlexaSkill is struct
type AlexaSkill struct {
	client       *http.Client
	cookieList   []*http.Cookie
	headers      map[string]string
	categoryURLs []string
	categoryName []string
}

// NewAlexaSkill return a new instance.
func newAlexaSkill(headerMap map[string]string) *AlexaSkill {
	log.Debug("create a new AlexaSkill instance")
	return &AlexaSkill{
		client:  &http.Client{},
		headers: headerMap,
	}
}

func (a *AlexaSkill) addHeaders(req *http.Request) {
	for key, value := range a.headers {
		req.Header.Add(key, value)
		log.Debugf("add header, key: %s, value: %s", key, value)
	}
}

func (a *AlexaSkill) addCookies(req *http.Request) {
	for _, c := range a.cookieList {
		req.AddCookie(c)
		log.Debugf("add cookie: %s", c.String())
	}
}

func (a *AlexaSkill) sendRequest(method, url string, needCookie bool) (*goquery.Document, error) {
	//为了应对反爬虫，在每次发起请求之前，先暂停 X 秒钟, X 是一个随机数
	pauseTime := commontool.RandInt64(2, 7)
	log.Debugf("Sleep %d seconds before request next URL: %s", pauseTime, url)
	time.Sleep(time.Duration(pauseTime) * time.Second)
	var resp *http.Response
	var err error
	req, _ := http.NewRequest(method, url, nil)
	a.addHeaders(req)
	if needCookie {
		a.addCookies(req)
	}
	// try 3 times.
	for i := 0; i < 3; i++ {
		resp, err = a.client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			// get cookies
			a.cookieList = resp.Cookies()
			bd, _ := ioutil.ReadAll(resp.Body)
			log.Debugf("response code: %d", resp.StatusCode)
			log.Debugf("get content: %s", string(bd))
			doc, err := goquery.NewDocumentFromResponse(resp)
			if err == nil {
				return doc, nil
			}
			return nil, err
		}
		time.Sleep(time.Duration(commontool.RandInt64(2, 6)) * time.Second)
		log.Warnf("request fail, try times: %d, err: %s", i, err)
	}
	log.Errorf("request: %s, error: %s, ", req.URL.Path, err.Error())
	return nil, err
}

// get category urls and store their names.
func (a *AlexaSkill) getCategoryURLs(url string) {
	doc, err := a.sendRequest("GET", url, false)
	if err != nil {
		log.Errorf("goquery parse doc from response fail: %s", err.Error())
		os.Exit(1)
	}
	li := doc.Find("div.categoryRefinementsSection").Find("ul").Find("li")
	log.Infof("Found category URL number: %d", li.Length())
	myFunc := func(index int, sel *goquery.Selection) {
		// get category urls.
		href, bl := sel.Find("a").Attr("href")
		if bl {
			a.categoryURLs = append(a.categoryURLs, "https://www.amazon.com"+href)
			log.Infof("category index: %d, url: %s", index, href)
		}
		// get category name.
		name := sel.Find("span").First().Text()
		if name != "" {
			a.categoryName = append(a.categoryName, name)
			log.Infof("category name: %s", name)
		}
	}
	li.Each(myFunc)
}

// 当进入某个category之后，取当前页面的16个item的URL
func (a *AlexaSkill) getCurrentPageItems(doc *goquery.Document) ([]string, string) {
	var allItemsURL []string
	li := doc.Find("div.s-item-container")
	log.Infof("Found %d items on current page.", li.Length())
	getItemURLFunc := func(index int, sel *goquery.Selection) {
		url, bl := sel.Find("a").First().Attr("href")
		if bl {
			allItemsURL = append(allItemsURL, url)
			log.Infof("get item url: %s", url)
		}
	}
	li.Each(getItemURLFunc)
	// get next page url
	nextURL := a.getNextPageURL(doc)
	log.Infof("get next page link is: %s", nextURL)
	return allItemsURL, nextURL
}

// 获取app的title, review star, review number, cmd 等数据
func (a *AlexaSkill) getData(categoryName, url string, doc *goquery.Document) []string {
	var contentList []string
	// get cmd title
	title := doc.Find("h1.a2s-title-content").Text()
	title = strings.TrimSpace(title)
	log.Infof("App category is: %s", categoryName)
	log.Infof("App title is: %s", title)
	log.Infof("App url is: %s", url)
	// get review star
	starStr, bl := doc.Find("div.a2s-pdd-reviews").Find("i").Attr("class")
	if bl {
		starStr = getReviewStar(starStr)
		log.Infof("App review star is: %s", starStr)
	}
	// get review number
	starNum := doc.Find("div.a2s-pdd-reviews").Find("span").Eq(1).Text()
	// strings.Replace(starNum, " ", "", -1)
	starNum = strings.TrimSpace(starNum)
	log.Infof("App review number is: %s", starNum)

	contentList = append(contentList, categoryName, title, url, starStr, starNum)
	// get cmd span.a2s-utterance-text
	li := doc.Find("ol[class='a-carousel'][role='list']").First().Find("li")
	getCMDFunc := func(index int, sel *goquery.Selection) {
		cmdName := sel.Find("span.a2s-utterance-text").Text()
		if cmdName != "" {
			//去掉两边空格
			cmdName = strings.TrimSpace(cmdName)
			//去掉双引号
			cmdName = strings.Replace(cmdName, `"`, "", -1)
			log.Infof("App's %d cmd is: %s", index, cmdName)
			contentList = append(contentList, cmdName)
		}
	}
	li.Each(getCMDFunc)
	return contentList
}

//对review star 做二次处理
func getReviewStar(str string) string {
	ss := strings.Split(str, " ")
	if len(ss) >= 3 {
		s1 := strings.Replace(ss[2], "a-star-small-", "", -1)
		if strings.Contains(s1, "-") {
			s1 = strings.Replace(s1, "-", ".", -1)
		}
		return s1
	}
	return "0"
}

// 获取下一页(NextPage)的URL
func (a *AlexaSkill) getNextPageURL(doc *goquery.Document) string {
	link, bl := doc.Find("a#pagnNextLink").Attr("href")
	if bl {
		return "https://www.amazon.com" + link
	}
	return ""
}

// 大循环+小循环，遍历所有的category的所有items
func (a *AlexaSkill) runEngineVersion2(cateIndex, pageIndex int) {
	var currentItemURLs []string
	var cateName string
	var content []string
	var nextPageURL string
	var flag bool
	var doc *goquery.Document
	var err error
	var currentPageIndex int
	// i, categoryURL := range a.categoryURLs
	for i := cateIndex; i < len(a.categoryURLs); i++ {
		currentPageIndex = 1
		doc, err = a.sendRequest("GET", a.categoryURLs[i], true)
		if err != nil {
			log.Errorf("Enter category fail: %s", a.categoryURLs[i])
			continue
		}
		log.Infof("进入第 %d 个 Category, URL: %s", i, a.categoryURLs[i])
		cateName = a.categoryName[i]
		flag = true

		for flag {
			// 取当前页的 16 个item的URL以及 Next Page 的URL
			log.Infof("Current page is Number.%d page.", currentPageIndex)
			currentItemURLs, nextPageURL = a.getCurrentPageItems(doc)
			log.Infof("current item number: %d", len(currentItemURLs))
			for _, itemURL := range currentItemURLs {
				// 如果指定到具体的某一页才开始收集数据， 则忽略掉前面的
				if pageIndex != 0 && currentPageIndex < pageIndex {
					break
				}
				//
				pageIndex = 0
				doc, err = a.sendRequest("GET", itemURL, true)
				if err != nil {
					log.Errorf("Enter Item page: %s fail, err: %s", itemURL, err.Error())
					continue
				}
				content = a.getData(cateName, itemURL, doc)
				writeToResult(content)
			}
			// 进入下一页
			if nextPageURL != "" {
				doc, err = a.sendRequest("GET", nextPageURL, true)
				if err != nil {
					log.Errorf("Enter next page fail: %s, current page is: %d", err.Error(), currentPageIndex)
				}
				log.Infof("Category: %s, page index: %d", a.categoryURLs[i], currentPageIndex)
				currentPageIndex++
			} else {
				flag = false
				log.Infof("No more next page.")
			}
		}
	}
}

// RunVersion2 hh.
func RunVersion2(url string, heads map[string]string, cIndex, pIndex int) {
	// create a csv file
	createNewFile("alexaData.csv")
	defer closeFile()

	inst := newAlexaSkill(heads)
	inst.getCategoryURLs(url)
	inst.runEngineVersion2(cIndex, pIndex)
}
