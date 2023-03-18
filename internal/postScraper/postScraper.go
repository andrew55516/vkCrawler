package postScraper

import (
	"bufio"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	db "vkCrawler/db/sqlc"
)

const root = "https://m.vk.com"

var Headers = map[string]string{
	"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	//"Accept-Encoding": "gzip, deflate, br",
	"Accept-Language": "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0",
	"Sec-Fetch-Dest":  "document",
}

var mu = &sync.Mutex{}
var wg sync.WaitGroup

type toScrape struct {
	postID int
	link   string
}

var proxyLists = []string{
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
	"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt",
	"https://raw.githubusercontent.com/shiftytr/proxy-list/master/proxy.txt",
}

func ScrapeAllPosts() {
	defer wg.Wait()

	posts := db.GetAllPosts()
	postsAmount := len(posts)

	toScrapeList := make([]toScrape, len(posts))

	for i, p := range posts {
		toScrapeList[i] = toScrape{
			postID: int(p.ID),
			link:   p.Link,
		}
	}

	scraped := 0

	for {

		proxies := getProxyList()
		fmt.Println(len(proxies))

		success := make(chan int, len(proxies))
		failed := make(chan toScrape, len(proxies))

		for i := 0; i < len(proxies); i++ {
			go scrapePost(toScrapeList[0], success, failed, i, proxies[i])
			toScrapeList = toScrapeList[1:]

		}

		failedProxies := 0

		for scraped < postsAmount {
			if len(proxies)-failedProxies < 5 {
				fmt.Println("proxies failed")
				break
			}
			select {
			case p := <-failed:
				toScrapeList = append([]toScrape{p}, toScrapeList...)
				failedProxies++
				fmt.Println(">>", failedProxies)

			case i := <-success:
				//fmt.Println(i)
				if len(toScrapeList) > 0 {
					go scrapePost(toScrapeList[0], success, failed, i, proxies[i])
					toScrapeList = toScrapeList[1:]
				}
				scraped++
				fmt.Println(scraped)

			}
		}

		if scraped == postsAmount {
			close(success)
			close(failed)
			break
		}

		for i := 0; i < len(proxies)-failedProxies; i++ {
			select {
			case p := <-failed:
				toScrapeList = append([]toScrape{p}, toScrapeList...)

			case <-success:
				//fmt.Println(i)
				scraped++
				fmt.Println(scraped)

			}
		}

		close(success)
		close(failed)

		time.Sleep(10 * time.Second)

	}

}

func scrapePost(post toScrape, success chan int, failed chan toScrape, order int, proxy string) {
	owner := "unknown"
	var reposts int

	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", root+post.link, nil)
	if err != nil {
		log.Fatal(err)
	}

	setHeaders(req)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		failed <- post
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		failed <- post
		return
	}

	likes := make([]string, 0)
	comments := make(map[string]struct{})

	s := doc.Find(`a[class="header__back  al_back mh mh_noleftmenu"]`)
	if s == nil {
		success <- order
		return
	}
	href, exists := s.Attr("href")
	if exists {
		owner = href
	} else {
		success <- order
		return
	}

	s = doc.Find(`a[class="PostBottomButton PostBottomButton--non-shrinkable"]`)
	label, exists := s.Attr("aria-label")
	if exists {
		t := strings.Split(label, " ")
		if k, err := strconv.Atoi(t[0]); err == nil {
			reposts = k
		}
	}

	s = doc.Find(`a[class="ReactionsPreview"]`)
	if s != nil {
		href, exists = s.Attr("href")
		if exists {
			t, _ := s.Attr("title")
			likesAmount, _ := strconv.Atoi(strings.Split(t, " ")[1])

			for i := 0; i*50 <= likesAmount; i++ {
				link := fmt.Sprintf("%s%s&offset=%v", root, href, i*50)
				likesReq, _ := http.NewRequest("GET", link, nil)
				setHeaders(likesReq)

				resp, err := client.Do(likesReq)
				if err != nil {
					log.Println(err)
					failed <- post
					return
				}

				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					log.Println(err)
					failed <- post
					return
				}

				resp.Body.Close()

				doc.Find(`a[class^="inline_item"]`).Each(func(i int, s *goquery.Selection) {
					href, exists := s.Attr("href")
					if exists {
						likes = append(likes, href)
					}
				})

				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	visited := make(map[string]struct{})

	for {
		doc.Find("a[class=\"ReplyItem__name\"]").Each(func(i int, s *goquery.Selection) {
			if s == nil {
				return
			}
			owner, exists := s.Attr("href")
			if exists {
				if _, ok := comments[owner]; !ok {
					comments[owner] = struct{}{}
				}
			}
		})
		var next string
		doc.Find("div[class=\"RepliesThreadNext Post__rowPadding\"]").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if s == nil {
				return true
			}
			href, exists := s.Children().Attr("href")
			if exists {
				if _, ok := visited[href]; !ok {
					visited[href] = struct{}{}
					next = href
					return false
				}
			}
			return true
		})
		if next == "" {
			break
		}

		req, err := http.NewRequest("GET", root+next, nil)
		if err != nil {
			log.Fatal(err)
		}

		setHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			failed <- post
			return
		}

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		resp.Body.Close()
		time.Sleep(500 * time.Millisecond)
	}

	commentList := make([]string, 0, len(comments))
	for c, _ := range comments {
		commentList = append(commentList, c)
	}

	fmt.Println(len(likes), len(commentList), reposts)

	wg.Add(1)
	go updatePost(post.postID, owner, likes, commentList, reposts)

	success <- order

}

func setHeaders(req *http.Request) *http.Request {
	for key, value := range Headers {
		req.Header.Set(key, value)
	}
	return req
}

func getProxyList() []string {
	proxies := make([]string, 0)
	testedProxies := make([]string, 0)

	file, err := os.Create("proxy-list.txt")
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range proxyLists {

		resp, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Write(body)
		if err != nil {
			log.Fatal(err)
		}

	}

	file.Close()

	file, err = os.Open("proxy-list.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		proxies = append(proxies, line)
	}
	fmt.Println(len(proxies))

	ch := make(chan string, len(proxies))

	for _, p := range proxies {
		go testProxy(p, ch)
	}

	for i := 0; i < len(proxies); i++ {
		select {
		case p := <-ch:
			if p != "" {
				testedProxies = append(testedProxies, p)
			}
		}
	}

	return testedProxies
}

func testProxy(proxy string, ch chan string) {

	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		log.Println(err)
		ch <- ""
		return
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://vk.com/", nil)
	if err != nil {
		log.Println(err)
		ch <- ""
		return
	}

	setHeaders(req)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		ch <- ""
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if doc == nil {
		ch <- ""
		return
	}

	s := doc.Find("#index_email")
	_, exists := s.Attr("name")
	if exists {
		ch <- proxy
	} else {
		ch <- ""
	}
}

func updatePost(id int, owner string, likes []string, comments []string, reposts int) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()

	db.UpdatePost(owner, len(likes), len(comments), reposts, id)
	for _, l := range likes {
		db.WriteDownLike(id, l)
	}

	for _, c := range comments {
		db.WriteDownComment(id, c)
	}
}
