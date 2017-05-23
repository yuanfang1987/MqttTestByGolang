package crawler

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"

	log "github.com/cihub/seelog"
)

func debuger() {
	log.Info("haha, stupid bird.")
}

// AccessAmazon  hh.
func AccessAmazon() {
	myClient := &http.Client{}
	myReq, _ := http.NewRequest("GET", "https://www.amazon.com/b?ie=UTF8&node=13727921011", nil)
	myReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0")
	myReq.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	resp, err := myClient.Do(myReq)
	if err != nil {
		log.Infof("reques fail: %s", err.Error())
		return
	}
	log.Info("start a new request.")
	defer resp.Body.Close()
	// debug
	allCookies := resp.Cookies()
	log.Infof("cookie length: %d", len(allCookies))
	for _, cookie := range allCookies {
		log.Infof("cookie info: %s", cookie.String())
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Infof("build new doc err: %s", err.Error())
		return
	}

	li := doc.Find("div.categoryRefinementsSection").Find("ul").Find("li")
	myFunc := func(index int, sel *goquery.Selection) {
		href, bl := sel.Find("a").Attr("href")
		if bl {
			log.Infof("category index: %d, url: %s", index, href)
		}
	}
	li.Each(myFunc)

	// end debug
	// bd, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("read body err: ", err)
	// 	return
	// }
	// fmt.Println("body is: ", string(bd))
}
