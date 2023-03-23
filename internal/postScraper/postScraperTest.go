package postScraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
)

func TestScrape() {

	file, err := os.Open("1.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	k := 0

	hrefs := make(map[string]struct{})

	doc.Find(`a[class="ReplyItem__name"]`).Each(func(i int, s *goquery.Selection) {

		href, exists := s.Attr("href")
		if exists {
			k++
			if _, ok := hrefs[href]; !ok {
				hrefs[href] = struct{}{}
			}
		}
	})

	fmt.Println(k, len(hrefs))
}

func TestScrapePost() {
	// post 3054
	//link := "/wall-22468706_111109"
	link := "/wall-181047282_129604"
	//link := "/wall-151650962_83259"

	//proxies := getProxyList()
	//fmt.Println(len(proxies))

	success := make(chan int, 1)
	failed := make(chan toScrape, 1)

	scrapePost(toScrape{
		postID: -1,
		link:   link,
	}, success, failed, 0, "")
}
