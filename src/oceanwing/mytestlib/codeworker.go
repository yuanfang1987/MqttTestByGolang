package mytestlib

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/PuerkitoBio/goquery"
)

func execShell() {
	// "ping", "-c4", "127.0.0.1"
	cmd := exec.Command("ping", "-c4", "127.0.0.1")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("error occur: %s\n", err)
		return
	}
	fmt.Printf("out put value is %s\n", out)
}

func mygoquery() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.amazon.com/b?ie=UTF8&node=13727921011", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request err: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Printf("parse to goquery error: %s", err.Error())
		return
	}
	li := doc.Find("div.categoryRefinementsSection").Find("li")
	fmt.Printf("li length: %d\n", li.Length())
	myfunc := func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		name := s.Find("span").First().Text()
		fmt.Printf("App name is: %s\n", name)
		fmt.Printf("App URL: %s\n", href)
	}
	li.Each(myfunc)
}
