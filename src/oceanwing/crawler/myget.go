package crawler

// "github.com/PuerkitoBio/goquery"
import (
	"io/ioutil"
	"net/http"

	"fmt"

	log "github.com/cihub/seelog"
)

func debuger() {
	log.Info("haha, stupid bird.")
}

// AccessAmazon  hh.
func AccessAmazon() {
	myClient := &http.Client{Timeout: 30}
	myReq, _ := http.NewRequest("GET", "https://www.amazon.com/b?ie=UTF8&node=13727921011", nil)
	myReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/538.1 (KHTML, like Gecko) PhantomJS/2.1.1 Safari/538.1")
	resp, err := myClient.Do(myReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body err: ", err)
		return
	}
	fmt.Println("body is: ", string(bd))
}
