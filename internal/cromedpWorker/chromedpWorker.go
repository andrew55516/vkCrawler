package cromedpWorker

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	db "vkCrawler/db/sqlc"
	"vkCrawler/internal/helpers"
	"vkCrawler/internal/postScraper"
)

const searchLink = "https://m.vk.com/search?c[q]=%23%D0%BA%D0%BB%D0%B8%D0%BC%D0%B0%D1%82&c[section]=statuses"

const pwd = "Crawler3799@!"

type post struct {
	link    string
	created time.Time
}

var resets int

var nums = [2]string{"89995370529", "89288099849"}

var opts = append(chromedp.DefaultExecAllocatorOptions[:],
	chromedp.Flag("headless", false),
	chromedp.Flag("disable-gpu", true),
	chromedp.Flag("no-sandbox", true),
	chromedp.Flag("ignore-certificate-errors", true),
	chromedp.WindowSize(800, 600),
)

var mu = sync.Mutex{}

func init() {

}

func login(ctx context.Context) {
	k := resets % 2
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://vk.com/`),
		chromedp.WaitReady(`input[name="login"]`, chromedp.ByQuery),
		chromedp.SendKeys(`#index_email`, nums[k]),
		chromedp.Click(`//button[@class="FlatButton FlatButton--primary FlatButton--size-l FlatButton--wide VkIdForm__button VkIdForm__signInButton"]`),
	)
	if err != nil {
		log.Fatal(err)
	}

	if k == 1 {
		err = chromedp.Run(ctx,
			chromedp.WaitVisible(`#otp`),
			chromedp.Click(`button[class="vkc__PureButton__button vkc__Link__link vkc__Link__primary vkc__Bottom__switchToPassword"`, chromedp.ByQuery),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = chromedp.Run(ctx,
		chromedp.WaitReady(`input[name="password"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="password"]`, pwd, chromedp.ByQuery),
		chromedp.Click(`span[class="vkuiButton__in"]`, chromedp.ByQuery),
		chromedp.WaitVisible(`#post_field`),
	)

	if err != nil {
		log.Fatal(err)
	}

}

func ScrapeUnknownPosts() {
	var wg sync.WaitGroup
	defer wg.Wait()

	posts := db.GetUnknownPosts()

	allocatorCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel2 := chromedp.NewContext(allocatorCtx)

	login(ctx)

	for len(posts) > 0 {

		likes := make([]string, 0)
		comments := make([]postScraper.Comment, 0)
		var reposts int
		var owner string

		chromedp.Run(ctx, chromedp.Navigate("https://m.vk.com/"+posts[0].Link))

		var html string
		chromedp.Run(ctx,
			chromedp.OuterHTML("html", &html),
		)

		reader := strings.NewReader(html)

		doc, _ := goquery.NewDocumentFromReader(reader)

		if doc.Find(`a[class="PostBottomButton PostBottomButton--non-shrinkable"]`).Length() == 0 {
			cancel1()
			cancel2()
			resets++
			allocatorCtx, cancel1 = chromedp.NewExecAllocator(context.Background(), opts...)

			ctx, cancel2 = chromedp.NewContext(allocatorCtx, chromedp.WithDebugf(log.Printf))

			login(ctx)
			continue

		}

		s := doc.Find(`a[class="header__back  al_back mh mh_noleftmenu"]`)
		href, _ := s.Attr("href")
		owner = href

		s = doc.Find(`a[class="PostBottomButton PostBottomButton--non-shrinkable"]`)
		label, exists := s.Attr("aria-label")
		if exists {
			t := strings.Split(label, " ")
			if k, err := strconv.Atoi(t[0]); err == nil {
				reposts = k
			}
		}

		if doc.Find(`a[class="RepliesThreadNext__link"]`).Length() > 0 {
			chromedp.Run(ctx,
				chromedp.Click(`a[class="RepliesThreadNext__link"]`, chromedp.ByQueryAll),
			)

			chromedp.Run(ctx,
				chromedp.OuterHTML("html", &html),
			)

			reader := strings.NewReader(html)

			doc, _ = goquery.NewDocumentFromReader(reader)
		}

		if doc.Find(`a[class="ReplyItem__name"]`).Length() > 0 {
			postID := strings.Split(posts[0].Link, "wall")[1]
			repliesSelector := fmt.Sprintf("#wall%s_replies", postID)

			replies := doc.Find(repliesSelector).Children()

			var lastThreadOwner string

			replies.Each(func(i int, s *goquery.Selection) {
				if s == nil {
					return
				}
				class, _ := s.Attr("class")

				switch class {
				case "ReplyItem Post__rowPadding":

					t := s.Find("a[class=\"ReplyItem__name\"]")
					commOwner, _ := t.Attr("href")

					timeStr := s.Find(`a[class="item_date"]`).Text()
					created := helpers.StrToTime(timeStr)

					comm := postScraper.Comment{
						Owner:       commOwner,
						ThreadOwner: owner,
						Created:     created,
					}

					comments = append(comments, comm)

					lastThreadOwner = commOwner

				case "RepliesThread":
					s.Find(`div[class="ReplyItem Post__rowPadding"]`).Each(func(i int, s *goquery.Selection) {
						if s == nil {
							return
						}

						t := s.Find("a[class=\"ReplyItem__name\"]")
						subCommOwner, _ := t.Attr("href")

						timeStr := s.Find(`a[class="item_date"]`).Text()
						created := helpers.StrToTime(timeStr)

						comm := postScraper.Comment{
							Owner:       subCommOwner,
							ThreadOwner: lastThreadOwner,
							Created:     created,
						}

						comments = append(comments, comm)
					})
				}
			})
		}

		if doc.Find(`a[class="ReactionsPreview"]`).Length() > 0 {

			s = doc.Find(`a[class="ReactionsPreview"]`)

			title, _ := s.Attr("title")
			t := strings.Split(title, " ")
			likesAmount, err := strconv.Atoi(t[1])
			if err != nil {
				log.Fatal(err)
			}

			href, _ = s.Attr("href")

			for i := 0; i*50 < likesAmount; i++ {
				link := fmt.Sprintf("https://m.vk.com/%s&offset=%v", href, i*50)

				chromedp.Run(ctx,
					chromedp.Navigate(link),
					chromedp.OuterHTML("html", &html),
				)

				reader := strings.NewReader(html)

				doc, _ = goquery.NewDocumentFromReader(reader)

				if doc.Find(`a[class^="inline_item"]`).Length() > 0 {
					doc.Find(`a[class^="inline_item"]`).Each(func(i int, s *goquery.Selection) {
						href, exists := s.Attr("href")
						if exists {
							likes = append(likes, href)
						}
					})
				} else {
					//cancel()
					cancel1()
					cancel2()
					resets++
					allocatorCtx, cancel1 = chromedp.NewExecAllocator(context.Background(), opts...)

					ctx, cancel2 = chromedp.NewContext(allocatorCtx)

					login(ctx)
					i--
					continue
				}
			}

		}
		fmt.Println(owner, likes, comments, reposts)

		wg.Add(1)
		go updatePost(int(posts[0].ID), owner, likes, comments, reposts, &wg)
		posts = posts[1:]
	}

	cancel1()
	cancel2()
}

func Crawl(wg *sync.WaitGroup) {
	end := false
	for !end {
		allocatorCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)

		ctx, cancel2 := chromedp.NewContext(allocatorCtx)

		login(ctx)

		lastPost := db.GetLastPost()
		toDate := lastPost.CreatedAt.Unix()
		if toDate < 0 {
			toDate = 0
		}
		end = getPostsToDate(ctx, toDate, wg)
		cancel1()
		cancel2()
	}
}

func getPostsToDate(ctx context.Context, toDate int64, wg *sync.WaitGroup) bool {
	curSearchLink := searchLink
	if toDate == 0 {
		curSearchLink = searchLink + "&offset="
	} else {
		curSearchLink = searchLink + fmt.Sprintf("&c[end_time]=%d&offset=", toDate)
	}

	end := false

	var nodes []*cdp.Node
	for i := 0; i < 30; i++ {
		var newNodes []*cdp.Node
		ctx1, _ := context.WithTimeout(ctx, 15*time.Second)
		searchLinkWithOffset := curSearchLink + strconv.Itoa(i*30)
		time.Sleep(1 * time.Second)
		chromedp.Run(ctx1,
			chromedp.Navigate(searchLinkWithOffset),
			chromedp.WaitReady(`a[class="wi_date al_wall"]`, chromedp.ByQuery),
			chromedp.Nodes(`a[class="wi_date al_wall"]`, &newNodes, chromedp.ByQueryAll),
		)
		if len(newNodes) == 0 {
			end = true
			resets++
			break
		}
		nodes = append(nodes, newNodes...)
	}

	if len(nodes) == 0 {
		return true
	}

	posts := make([]post, len(nodes))
	for i, node := range nodes {
		link := node.Attributes[3]
		link = strings.Split(link, "?")[0]

		var created time.Time
		if len(node.Attributes) == 10 {
			created = helpers.StrToTime(node.Attributes[9])
		} else {
			created = helpers.StrToTime(node.Children[0].NodeValue)
		}

		posts[i] = post{
			link:    link,
			created: created,
		}

	}

	wg.Wait()
	wg.Add(1)
	go putPosts(posts, wg)

	if end {
		wg.Wait()
		return false
	}

	nextToDate := posts[len(posts)-1].created.Unix()
	if nextToDate == toDate {
		nextToDate += int64(time.Hour.Seconds())
	}

	return getPostsToDate(ctx, nextToDate, wg)

}

func putPosts(posts []post, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, post := range posts {
		err := db.WriteDownPost(post.link, "unknown", -1, -1, -1, post.created)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func updatePost(id int, owner string, likes []string, comments []postScraper.Comment, reposts int, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	//wg.Wait()

	db.UpdatePost(owner, len(likes), len(comments), reposts, id)
	for _, l := range likes {
		db.WriteDownLike(id, l)
	}

	for _, c := range comments {
		db.WriteDownComment(id, c.Owner, c.ThreadOwner, c.Created)
	}
}
