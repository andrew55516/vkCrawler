package main

import "vkCrawler/internal/postScraper"

func main() {
	//scraped := db.GetLastPost().ID
	//for scraped < 9000 {
	//	var wg sync.WaitGroup
	//	cromedpWorker.Crawl(&wg)
	//	wg.Wait()
	//	scraped = db.GetLastPost().ID
	//	if scraped < 9000 {
	//		time.Sleep(10 * time.Minute)
	//	}
	//}
	//
	postScraper.ScrapeAllPosts()
	//cromedpWorker.ScrapeUnknownPosts()

	//postScraper.TestScrapePost()

}
