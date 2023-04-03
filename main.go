package main

import (
	"log"
	db "vkCrawler/db/sqlc"
)

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
	//postScraper.ScrapeAllPosts()
	//cromedpWorker.ScrapeUnknownPosts()

	//postScraper.TestScrapePost()

	err := db.FillAllNodes()
	if err != nil {
		log.Fatal(err)
	}

	err = db.FillAllEdges()
	if err != nil {
		log.Fatal(err)
	}

}
