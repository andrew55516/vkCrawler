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

	ctx, cancel2 := chromedp.NewContext(allocatorCtx, chromedp.WithDebugf(log.Printf))

	login(ctx)

	for len(posts) > 0 {

		likes := make([]string, 0)
		comments := make([]string, 0)
		var reposts int
		var owner string

		var ownerNodes, previewReactionNodes, repostNodes, commentNodes []*cdp.Node

		chromedp.Run(ctx, chromedp.Navigate("https://m.vk.com/"+posts[0].Link))

		//ctx1, cancel := context.WithTimeout(ctx, 10*time.Second)
		var html string
		chromedp.Run(ctx,
			//chromedp.Navigate("https://m.vk.com/"+posts[0].Link),
			//chromedp.WaitVisible(`body`, chromedp.BySearch),
			//chromedp.Nodes(`a[class="ReactionsPreview"]`, &previewReactionNodes, chromedp.ByQuery),
			chromedp.OuterHTML("html", &html),
			//chromedp.Nodes(`a[class="PostBottomButton PostBottomButton--non-shrinkable"]`, &repostNodes, chromedp.ByQuery),
			//chromedp.Nodes(`a[class="header__back  al_back mh mh_noleftmenu"]`, &ownerNodes, chromedp.ByQuery),
			//chromedp.Nodes(`a[class="ReplyItem__name"]`, &commentNodes, chromedp.ByQueryAll),
			//chromedp.Nodes(`a[class="RepliesThreadNext__link"]`, &nextCommentsNodes, chromedp.ByQueryAll),
		)

		//os.WriteFile("unknown.html", []byte(html), 0644)

		//var likesExists, commentsExists, repostsExists, nextCommentsExists bool

		reader := strings.NewReader(html)

		doc, err := goquery.NewDocumentFromReader(reader)

		if doc.Find(`a[class="PostBottomButton PostBottomButton--non-shrinkable"]`).Length() == 0 {
			cancel1()
			cancel2()
			resets++
			allocatorCtx, cancel1 = chromedp.NewExecAllocator(context.Background(), opts...)

			ctx, cancel2 = chromedp.NewContext(allocatorCtx, chromedp.WithDebugf(log.Printf))

			login(ctx)
			continue

		} else {
			chromedp.Run(ctx,
				chromedp.Nodes(`a[class="PostBottomButton PostBottomButton--non-shrinkable"]`, &repostNodes, chromedp.ByQuery),
				chromedp.Nodes(`a[class="header__back  al_back mh mh_noleftmenu"]`, &ownerNodes, chromedp.ByQuery),
			)
		}

		if doc.Find(`a[class="RepliesThreadNext__link"]`).Length() > 0 {
			chromedp.Run(ctx,
				chromedp.Click(`a[class="RepliesThreadNext__link"]`, chromedp.ByQueryAll),
			)
		}

		if doc.Find(`a[class="ReplyItem__name"]`).Length() > 0 {
			chromedp.Run(ctx,
				chromedp.Nodes(`a[class="ReplyItem__name"]`, &commentNodes, chromedp.ByQueryAll),
			)
		}

		if doc.Find(`a[class="ReactionsPreview"]`).Length() > 0 {
			chromedp.Run(ctx,
				//chromedp.Navigate("https://m.vk.com/"+posts[0].Link),
				chromedp.Nodes(`a[class="ReactionsPreview"]`, &previewReactionNodes, chromedp.ByQuery),
			)

			title := previewReactionNodes[0].AttributeValue("title")
			t := strings.Split(title, " ")
			likesAmount, err := strconv.Atoi(t[1])
			if err != nil {
				log.Fatal(err)
			}

			href := previewReactionNodes[0].AttributeValue("href")

			for i := 0; i*50 < likesAmount; i++ {
				link := fmt.Sprintf("https://m.vk.com/%s&offset=%v", href, i*50)
				var likeNodes []*cdp.Node

				ctx1, cancel := context.WithTimeout(ctx, 20*time.Second)
				chromedp.Run(ctx1,
					chromedp.Navigate(link),
					chromedp.WaitReady(`a[class^="inline_item"]`, chromedp.ByQueryAll),
					chromedp.Nodes(`a[class^="inline_item"]`, &likeNodes, chromedp.ByQueryAll),
				)

				if len(likeNodes) == 0 {
					cancel()
					cancel1()
					cancel2()
					resets++
					allocatorCtx, cancel1 = chromedp.NewExecAllocator(context.Background(), opts...)

					ctx, cancel2 = chromedp.NewContext(allocatorCtx, chromedp.WithDebugf(log.Printf))

					login(ctx)
					i--
					continue
				}
				cancel()

				for _, n := range likeNodes {
					likes = append(likes, n.AttributeValue("href"))
				}
			}

		}

		//cancel()

		//ctx1, cancel = context.WithTimeout(ctx, 10*time.Second)
		//chromedp.Run(ctx,
		//	chromedp.Navigate("https://m.vk.com/"+posts[0].Link),
		//	//chromedp.WaitVisible(`body`, chromedp.BySearch),
		//	chromedp.Nodes(`a[class="ReactionsPreview"]`, &previewReactionNodes, chromedp.ByQuery),
		//)
		//cancel()

		//if len(repostNodes) == 0 {
		//	cancel1()
		//	cancel2()
		//	resets++
		//	allocatorCtx, cancel1 = chromedp.NewExecAllocator(context.Background(), opts...)
		//
		//	ctx, cancel2 = chromedp.NewContext(allocatorCtx, chromedp.WithDebugf(log.Printf))
		//
		//	login(ctx)
		//	continue
		//}

		owner = ownerNodes[0].AttributeValue("href")

		title := repostNodes[0].AttributeValue("aria-label")
		t := strings.Split(title, " ")
		reposts, err = strconv.Atoi(t[0])
		if err != nil {
			log.Fatal(err)
		}

		//for _, n := range nextCommentsNodes {
		//	chromedp.Run(ctx,
		//		chromedp.MouseClickNode(n),
		//		chromedp.Sleep(300*time.Millisecond),
		//	)
		//}
		//
		//if len(commentNodes) > 0 {
		//	chromedp.Run(ctx,
		//		chromedp.WaitReady(`a[class="ReplyItem__name"]`, chromedp.ByQueryAll),
		//		chromedp.Nodes(`a[class="ReplyItem__name"]`, &commentNodes, chromedp.ByQueryAll),
		//	)
		//}

		comm := make(map[string]struct{})

		for _, n := range commentNodes {
			owner := n.AttributeValue("href")
			if _, ok := comm[owner]; !ok {
				comm[owner] = struct{}{}
			}
		}

		for c, _ := range comm {
			comments = append(comments, c)
		}

		//if len(previewReactionNodes) > 0 {
		//	title := previewReactionNodes[0].AttributeValue("title")
		//	t := strings.Split(title, " ")
		//	likesAmount, err := strconv.Atoi(t[1])
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//
		//	href := previewReactionNodes[0].AttributeValue("href")
		//
		//	for i := 0; i*50 < likesAmount; i++ {
		//		link := fmt.Sprintf("https://m.vk.com/%s&offset=%v", href, i*50)
		//		chromedp.Run(ctx,
		//			chromedp.Navigate(link),
		//			chromedp.WaitReady(`a[class^="inline_item"]`, chromedp.ByQueryAll),
		//			chromedp.Nodes(`a[class^="inline_item"]`, &likeNodes, chromedp.ByQueryAll),
		//		)
		//		for _, n := range likeNodes {
		//			likes = append(likes, n.AttributeValue("href"))
		//		}
		//	}
		//}

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

		ctx, cancel2 := chromedp.NewContext(allocatorCtx, chromedp.WithDebugf(log.Printf))

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

func updatePost(id int, owner string, likes []string, comments []string, reposts int, wg *sync.WaitGroup) {
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
