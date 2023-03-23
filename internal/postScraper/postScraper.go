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
	"vkCrawler/internal/helpers"
)

const root = "https://m.vk.com"

var Headers = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	"Accept-Language": "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0",
	"Sec-Fetch-Dest":  "document",
}

var wg sync.WaitGroup
var wg1 sync.WaitGroup

type toScrape struct {
	postID int
	link   string
}

type Comment struct {
	Owner       string
	ThreadOwner string
	Created     time.Time
}

var proxyLists = []string{
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
	"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt",
	"https://raw.githubusercontent.com/shiftytr/proxy-list/master/proxy.txt",
}

var updatedPosts = struct {
	updated int
	mu      sync.Mutex
}{
	updated: 0,
	mu:      sync.Mutex{},
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
			wg1.Add(1)
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
				fmt.Println(">>", len(proxies)-failedProxies)

			case i := <-success:
				//fmt.Println(i)
				if len(toScrapeList) > 0 {
					wg1.Add(1)
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
		fmt.Println("waiting for started goroutines...")
		wg1.Wait()
		fmt.Println("started goroutines done")

		end := false
		for !end {
			select {
			case p := <-failed:
				toScrapeList = append([]toScrape{p}, toScrapeList...)
			case <-success:
				scraped++
				fmt.Println(scraped)
			default:
				end = true
			}
		}

		close(success)
		close(failed)
		fmt.Println("restarting...")

		time.Sleep(10 * time.Second)

	}
}

func scrapePost(post toScrape, success chan int, failed chan toScrape, order int, proxy string) {
	defer wg1.Done()
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

	postID := strings.Split(post.link, "wall")[1]

	repliesSelector := fmt.Sprintf("#wall%s_replies", postID)

	comments := make([]Comment, 0)

	for {
		replies := doc.Find(repliesSelector).Children()

		if len(replies.Nodes) == 0 {
			break
		}

		var lastThreadOwner string

		nextComments := ""

		commentsFailed := false

		replies.Each(func(i int, s *goquery.Selection) {
			if s == nil || commentsFailed {
				return
			}
			class, _ := s.Attr("class")

			switch class {
			case "ReplyItem Post__rowPadding":
				t := s.Find("a[class=\"ReplyItem__name\"]")
				commOwner, _ := t.Attr("href")

				timeStr := s.Find(`a[class="item_date"]`).Text()
				created := helpers.StrToTime(timeStr)

				comm := Comment{
					Owner:       commOwner,
					ThreadOwner: owner,
					Created:     created,
				}

				comments = append(comments, comm)

				lastThreadOwner = commOwner

			case "ReplyItem ReplyItem_deleted Post__rowPadding":
				lastThreadOwner = "unknown"

			case "RepliesThread":

				threadID, _ := s.Attr("id")
				threadSelector := "#" + threadID

				threadDoc := doc
				first := true
				for {
					thread := threadDoc.Find(threadSelector)
					thread.Find(`div[class="ReplyItem Post__rowPadding"]`).Each(func(i int, s *goquery.Selection) {
						if s == nil {
							return
						}

						if first && lastThreadOwner == "unknown" {
							threadOwner, _ := s.Find(`a[class="mem_link"]`).Attr("href")
							if threadOwner != "" {
								lastThreadOwner = threadOwner
							}
							first = false
						}

						t := s.Find("a[class=\"ReplyItem__name\"]")
						subCommOwner, _ := t.Attr("href")

						timeStr := s.Find(`a[class="item_date"]`).Text()
						created := helpers.StrToTime(timeStr)

						comm := Comment{
							Owner:       subCommOwner,
							ThreadOwner: lastThreadOwner,
							Created:     created,
						}

						comments = append(comments, comm)
					})

					nextDiv := thread.Find("div[class=\"RepliesThreadNext Post__rowPadding\"]")

					if nextDiv == nil {
						break
					}

					nextHref, _ := nextDiv.Children().Attr("href")

					if nextHref == "" {
						break
					}

					req, err := http.NewRequest("GET", root+nextHref, nil)
					if err != nil {
						log.Fatal(err)
					}

					setHeaders(req)

					resp, err := client.Do(req)
					if err != nil {
						log.Println(err)
						commentsFailed = true
						return
					}

					threadDoc, err = goquery.NewDocumentFromReader(resp.Body)
					if err != nil {
						log.Fatal(err)
					}

					resp.Body.Close()
					time.Sleep(500 * time.Millisecond)

				}
			}

		})

		if commentsFailed {
			failed <- post
			return
		}

		next := replies.Nodes[len(replies.Nodes)-1]
		t := next.FirstChild.Attr[1].Val
		prev := strings.Contains(t, "prev")
		splitedHref := strings.Split(t, "#")
		if len(splitedHref) >= 2 && splitedHref[1] == "comments" && !prev {
			nextComments = t
		}

		if nextComments == "" {
			break
		}

		req, err := http.NewRequest("GET", root+nextComments, nil)
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
	}

	wg.Add(1)
	go updatePost(post.postID, owner, likes, comments, reposts)

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

	addProxies := getProxiesFromFreeProxyList()

	proxies = append(proxies, addProxies...)

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

func updatePost(id int, owner string, likes []string, comments []Comment, reposts int) {
	defer wg.Done()
	updatedPosts.mu.Lock()
	defer updatedPosts.mu.Unlock()
	updatedPosts.updated++

	db.UpdatePost(owner, len(likes), len(comments), reposts, id)

	for _, l := range likes {
		db.WriteDownLike(id, l)
	}

	for _, c := range comments {
		db.WriteDownComment(id, c.Owner, c.ThreadOwner, c.Created)
	}
}

func getProxiesFromFreeProxyList() []string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://free-proxy-list.net/", nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	setHeaders(req)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	proxiesText := doc.Find(`textarea[class="form-control"]`).Text()
	proxyList := strings.Split(strings.TrimSuffix(proxiesText, "\n"), "\n")[3:]

	return proxyList
}
