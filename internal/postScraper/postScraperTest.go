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

func TestLikes() {
	link := "/wall-69607966_26744"

	proxies := getProxyList()
	fmt.Println(len(proxies))

	success := make(chan int, len(proxies))
	failed := make(chan toScrape, len(proxies))

	scrapePost(toScrape{
		postID: -1,
		link:   link,
	}, success, failed, 0, proxies[0])
}
